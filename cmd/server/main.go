package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kyambuthia/go-dot-server/internal/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := server.LoadConfigFromEnv()
	if err := cfg.Validate(); err != nil {
		logger.Error("invalid server configuration", "error", err)
		os.Exit(1)
	}

	srvr := server.BuildHTTPServer(cfg)
	logger.Info("server starting", "url", "https://"+cfg.Addr, "public_dir", cfg.PublicDir)

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-sigCtx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.Info("shutdown signal received, stopping server")
		if err := srvr.Shutdown(shutdownCtx); err != nil {
			logger.Error("graceful shutdown failed", "error", err)
		}
	}()

	if err := srvr.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server crashed", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
