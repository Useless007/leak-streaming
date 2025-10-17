package middleware

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

type RateLimitConfig struct {
	RequestsPerMinute int
	Burst             int
	KeyExtractor      func(*http.Request) string
	Window            time.Duration
}

func RateLimit(cfg RateLimitConfig) func(http.Handler) http.Handler {
	if cfg.RequestsPerMinute <= 0 {
		cfg.RequestsPerMinute = 60
	}

	if cfg.Burst <= 0 {
		cfg.Burst = cfg.RequestsPerMinute
	}
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}

	limiterStore := newLimiterStore(cfg)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := cfg.KeyExtractor(r)
			if key == "" {
				key = clientIP(r)
			}

			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			limiter := limiterStore.get(key)
			if !limiter.Allow() {
				w.Header().Set("Retry-After", "60")
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RateLimitRedis(client *redis.Client, cfg RateLimitConfig) func(http.Handler) http.Handler {
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	if cfg.RequestsPerMinute <= 0 {
		cfg.RequestsPerMinute = 60
	}
	if cfg.Burst <= 0 {
		cfg.Burst = cfg.RequestsPerMinute
	}

	limiter := newRedisLimiter(client, cfg)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := cfg.KeyExtractor(r)
			if key == "" {
				key = clientIP(r)
			}
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			allowed, retryAfter := limiter.Allow(r.Context(), key)
			if !allowed {
				if retryAfter > 0 {
					seconds := int(retryAfter / time.Second)
					if retryAfter%time.Second != 0 {
						seconds++
					}
					if seconds <= 0 {
						seconds = 1
					}
					w.Header().Set("Retry-After", strconv.Itoa(seconds))
				}
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type limiterStore struct {
	cfg RateLimitConfig
	mu  sync.Mutex
	m   map[string]*rate.Limiter
	ttl time.Duration
}

func newLimiterStore(cfg RateLimitConfig) *limiterStore {
	return &limiterStore{
		cfg: cfg,
		m:   make(map[string]*rate.Limiter),
		ttl: cfg.Window,
	}
}

func (s *limiterStore) get(key string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	limiter, ok := s.m[key]
	if !ok {
		limit := rate.Every(s.cfg.Window / time.Duration(s.cfg.RequestsPerMinute))
		limiter = rate.NewLimiter(limit, s.cfg.Burst)
		s.m[key] = limiter
	}

	return limiter
}

type redisLimiter struct {
	client *redis.Client
	cfg    RateLimitConfig
	script *redis.Script
}

func newRedisLimiter(client *redis.Client, cfg RateLimitConfig) *redisLimiter {
	lua := redis.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local allowed = 0
local ttl = redis.call("PTTL", key)
if ttl < 0 then
  ttl = window
end
local current = redis.call("INCR", key)
if current == 1 then
  redis.call("PEXPIRE", key, window)
  ttl = window
end
if current <= limit then
  allowed = 1
end
return {allowed, ttl}
    `)

	return &redisLimiter{
		client: client,
		cfg:    cfg,
		script: lua,
	}
}

func (l *redisLimiter) Allow(ctx context.Context, key string) (bool, time.Duration) {
	limit := l.cfg.Burst
	if limit <= 0 {
		limit = l.cfg.RequestsPerMinute
	}
	result, err := l.script.Run(ctx, l.client, []string{l.redisKey(key)}, limit, l.cfg.Window.Milliseconds()).Result()
	if err != nil {
		return true, 0
	}

	values, ok := result.([]interface{})
	if !ok || len(values) != 2 {
		return true, 0
	}

	allowed, ok := values[0].(int64)
	if !ok {
		return true, 0
	}

	ttl, ok := values[1].(int64)
	if !ok || ttl <= 0 {
		return allowed == 1, 0
	}

	return allowed == 1, time.Duration(ttl) * time.Millisecond
}

func (l *redisLimiter) redisKey(key string) string {
	return "rate:limiter:" + key
}

func clientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		parts := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	return ip
}
