package router

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"

	"github.com/leak-streaming/leak-streaming/backend/internal/api/health"
	apimiddleware "github.com/leak-streaming/leak-streaming/backend/internal/api/middleware"
	apimovies "github.com/leak-streaming/leak-streaming/backend/internal/api/movies"
	"github.com/leak-streaming/leak-streaming/backend/internal/platform/config"
	"github.com/leak-streaming/leak-streaming/backend/internal/platform/telemetry"
	servicemovies "github.com/leak-streaming/leak-streaming/backend/internal/service/movies"
)

func NewServer(cfg config.Config, log *slog.Logger, redisClient *redis.Client, movieService *servicemovies.Service) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.HTTP.WriteTimeout))
	r.Use(apimiddleware.SecureHeaders())
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"X-Correlation-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(telemetry.CorrelationMiddleware)
	rateLimitCfg := apimiddleware.RateLimitConfig{
		RequestsPerMinute: 120,
		Burst:             240,
		Window:            time.Minute,
		KeyExtractor: func(r *http.Request) string {
			forwarded := r.Header.Get("X-Forwarded-For")
			if forwarded == "" {
				return ""
			}
			parts := strings.Split(forwarded, ",")
			return strings.TrimSpace(parts[0])
		},
	}
	if redisClient != nil {
		r.Use(apimiddleware.RateLimitRedis(redisClient, rateLimitCfg))
	} else {
		r.Use(apimiddleware.RateLimit(rateLimitCfg))
	}
	r.Use(apimiddleware.RequestLogger(log))

	health.RegisterRoutes(r)

	if movieService != nil {
		listHandler := apimovies.NewListHandler(movieService)
		detailsHandler := apimovies.NewDetailsHandler(movieService)
		streamHandler := apimovies.NewStreamTokenHandler(movieService)
		manifestHandler := apimovies.NewManifestHandler(movieService)
		segmentHandler := apimovies.NewSegmentHandler(movieService)
		createHandler := apimovies.NewCreateHandler(movieService)
		r.Route("/movies", func(r chi.Router) {
			r.Get("/", listHandler.ServeHTTP)
			r.Post("/", createHandler.ServeHTTP)
			r.Get("/{slug}", detailsHandler.ServeHTTP)
			r.Post("/{slug}/playback-token", streamHandler.ServeHTTP)
			r.Get("/{slug}/manifest.m3u8", manifestHandler.ServeHTTP)
			r.Get("/{slug}/segment", segmentHandler.ServeHTTP)
		})
	}

	return &http.Server{
		Addr:         cfg.HTTP.Address(),
		Handler:      r,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}
}

func Shutdown(server *http.Server, ctx context.Context) error {
	return server.Shutdown(ctx)
}
