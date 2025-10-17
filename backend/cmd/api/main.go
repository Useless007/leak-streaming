package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/leak-streaming/leak-streaming/backend/internal/api/router"
	"github.com/leak-streaming/leak-streaming/backend/internal/persistence/repository"
	"github.com/leak-streaming/leak-streaming/backend/internal/platform/cache"
	"github.com/leak-streaming/leak-streaming/backend/internal/platform/config"
	"github.com/leak-streaming/leak-streaming/backend/internal/platform/logger"
	"github.com/leak-streaming/leak-streaming/backend/internal/platform/telemetry"
	movieservice "github.com/leak-streaming/leak-streaming/backend/internal/service/movies"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	log := logger.New(cfg.Env)

	telemetryShutdown := func(context.Context) error { return nil }
	if shutdown, err := telemetry.Setup(ctx, cfg.Telemetry, log); err != nil {
		log.Warn("failed to initialize telemetry", "error", err)
	} else {
		telemetryShutdown = shutdown
		defer func() {
			if err := telemetryShutdown(context.Background()); err != nil {
				log.Warn("failed to shutdown telemetry provider", "error", err)
			}
		}()
	}

	redisClient, err := cache.New(ctx, cfg.Redis)
	if err != nil {
		log.Warn("failed to connect to redis", "error", err)
		redisClient = nil
	} else {
		defer func() {
			if err := redisClient.Close(); err != nil {
				log.Warn("failed to close redis client", "error", err)
			}
		}()
	}

	repo := repository.NewMovieRepository(nil)
	var tokenSigner movieservice.TokenSigner
	if redisClient != nil {
		tokenSigner = movieservice.NewRedisTokenSigner(redisClient)
	} else {
		tokenSigner = movieservice.NewInMemoryTokenSigner()
	}
	movieService := movieservice.NewService(repo, tokenSigner, cfg.Stream.TokenTTL)

	server := router.NewServer(cfg, log, redisClient, movieService)

	go func() {
		log.Info("api server starting", "addr", cfg.HTTP.Address())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("api server encountered an error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := router.Shutdown(server, shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "error", err)
		return
	}

	log.Info("api server stopped gracefully")
}
