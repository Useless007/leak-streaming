package movies

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenSigner struct {
	client *redis.Client
}

func NewRedisTokenSigner(client *redis.Client) *RedisTokenSigner {
	return &RedisTokenSigner{client: client}
}

func (s *RedisTokenSigner) SignToken(movieID, viewerID string, ttl time.Duration) (string, error) {
	if s == nil || s.client == nil {
		return generateRandomToken(), nil
	}

	token := generateRandomToken()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	value := movieID
	if viewerID != "" {
		value = fmt.Sprintf("%s|%s", movieID, viewerID)
	}

	if err := s.client.Set(ctx, redisPlaybackKey(token), value, ttl).Err(); err != nil {
		return "", err
	}

	return token, nil
}

func (s *RedisTokenSigner) ValidateToken(token, movieID string) (bool, error) {
	if s == nil || s.client == nil {
		return true, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	value, err := s.client.Get(ctx, redisPlaybackKey(token)).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	parts := strings.Split(value, "|")
	return parts[0] == movieID, nil
}

type InMemoryTokenSigner struct {
	store map[string]string
}

func NewInMemoryTokenSigner() *InMemoryTokenSigner {
	return &InMemoryTokenSigner{store: make(map[string]string)}
}

func (s *InMemoryTokenSigner) SignToken(movieID, viewerID string, ttl time.Duration) (string, error) {
	token := generateRandomToken()
	s.store[token] = movieID
	return token, nil
}

func (s *InMemoryTokenSigner) ValidateToken(token, movieID string) (bool, error) {
	stored, ok := s.store[token]
	if !ok {
		return false, nil
	}
	return stored == movieID, nil
}

func redisPlaybackKey(token string) string {
	return "playback:token:" + token
}

func generateRandomToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
