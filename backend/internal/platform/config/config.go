package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/leak-streaming/leak-streaming/backend/internal/platform/cache"
	"github.com/leak-streaming/leak-streaming/backend/internal/platform/telemetry"
)

type HTTPConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func (h HTTPConfig) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

type Config struct {
	Env       string
	HTTP      HTTPConfig
	Database  DatabaseConfig
	Redis     cache.RedisConfig
	Telemetry telemetry.Config
	Stream    StreamConfig
}

type StreamConfig struct {
	TokenTTL time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	MigrationsDir   string
}

func Load() (Config, error) {
	http := HTTPConfig{
		Host:            getEnv("HTTP_HOST", "0.0.0.0"),
		Port:            getEnvAsInt("HTTP_PORT", 8080),
		ReadTimeout:     getEnvAsDuration("HTTP_READ_TIMEOUT_MS", 5*time.Second),
		WriteTimeout:    getEnvAsDuration("HTTP_WRITE_TIMEOUT_MS", 10*time.Second),
		IdleTimeout:     getEnvAsDuration("HTTP_IDLE_TIMEOUT_MS", 120*time.Second),
		ShutdownTimeout: getEnvAsDuration("HTTP_SHUTDOWN_TIMEOUT_MS", 15*time.Second),
	}

	return Config{
		Env:  getEnv("APP_ENV", "development"),
		HTTP: http,
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "127.0.0.1"),
			Port:            getEnvAsInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "leakstream"),
			Password:        getEnv("DB_PASSWORD", "leakstream"),
			Name:            getEnv("DB_NAME", "leakstream"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDurationSeconds("DB_CONN_MAX_LIFETIME_SEC", 300),
			MigrationsDir:   getEnv("DB_MIGRATIONS_DIR", "./internal/persistence/migrations"),
		},
		Redis: cache.RedisConfig{
			Host:            getEnv("REDIS_HOST", "127.0.0.1"),
			Port:            getEnvAsInt("REDIS_PORT", 6379),
			Password:        getEnv("REDIS_PASSWORD", ""),
			DB:              getEnvAsInt("REDIS_DB", 0),
			UseTLS:          getEnvAsBool("REDIS_TLS_ENABLED", false),
			ReadTimeout:     getEnvAsDuration("REDIS_READ_TIMEOUT_MS", 500*time.Millisecond),
			WriteTimeout:    getEnvAsDuration("REDIS_WRITE_TIMEOUT_MS", 500*time.Millisecond),
			DialTimeout:     getEnvAsDuration("REDIS_DIAL_TIMEOUT_MS", 500*time.Millisecond),
			PoolSize:        getEnvAsInt("REDIS_POOL_SIZE", 20),
			MinIdleConns:    getEnvAsInt("REDIS_MIN_IDLE_CONNS", 2),
			MaxRetries:      getEnvAsInt("REDIS_MAX_RETRIES", 3),
			HealthCheckFreq: getEnvAsDuration("REDIS_HEALTHCHECK_MS", 30*time.Second),
		},
		Telemetry: telemetry.Config{
			ServiceName:  getEnv("OTEL_SERVICE_NAME", "leak-streaming-api"),
			Environment:  getEnv("OTEL_ENVIRONMENT", getEnv("APP_ENV", "development")),
			OTLPEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
			Insecure:     getEnvAsBool("OTEL_EXPORTER_OTLP_INSECURE", true),
			SampleRatio:  getEnvAsFloat("OTEL_SAMPLE_RATIO", 0.25),
		},
		Stream: StreamConfig{
			TokenTTL: getEnvAsDurationSeconds("STREAM_TOKEN_TTL_SEC", 300),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr, ok := os.LookupEnv(key)
	if !ok || valueStr == "" {
		return fallback
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}

	return value
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	valueStr, ok := os.LookupEnv(key)
	if !ok || valueStr == "" {
		return fallback
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}

	return time.Duration(value) * time.Millisecond
}

func getEnvAsBool(key string, fallback bool) bool {
	valueStr, ok := os.LookupEnv(key)
	if !ok || valueStr == "" {
		return fallback
	}

	switch valueStr {
	case "true", "1", "yes", "on", "enabled":
		return true
	case "false", "0", "no", "off", "disabled":
		return false
	default:
		return fallback
	}
}

func getEnvAsFloat(key string, fallback float64) float64 {
	valueStr, ok := os.LookupEnv(key)
	if !ok || valueStr == "" {
		return fallback
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fallback
	}

	return value
}

func getEnvAsDurationSeconds(key string, fallback int) time.Duration {
	valueStr, ok := os.LookupEnv(key)
	if !ok || valueStr == "" {
		return time.Duration(fallback) * time.Second
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return time.Duration(fallback) * time.Second
	}

	return time.Duration(value) * time.Second
}
