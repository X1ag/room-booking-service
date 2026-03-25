package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"test-backend-1-X1ag/internal/app"
	"test-backend-1-X1ag/internal/config"
	"test-backend-1-X1ag/internal/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	baseLogger, err := logger.NewZerologLogger(cfg.Logger)
	if err != nil {
		log.Panicln("cant create logger, err: ", err)
		return
	}

	application, err := app.New(ctx, cfg, baseLogger)
	if err != nil {
		log.Fatalf("create app: %v", err)
	}
	defer application.Close()

	srv := &http.Server{
		Addr:         cfg.HTTP.Addr(),
		Handler:      application.Router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	go func() {
		application.Logger.Info().Str("addr", cfg.HTTP.Addr()).Msg("HTTP server started")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("run HTTP server: %v", err)
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-shutdownCtx.Done()
	application.Logger.Info().Msg("shutdown signal received")

	ctxTimeout, cancelShutdown := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		application.Logger.Error().Err(err).Msg("graceful shutdown failed")
		return
	}

	application.Logger.Info().Msg("http server stopped gracefully")
}
