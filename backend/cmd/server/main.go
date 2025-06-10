package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vk-worker/internal/config"
	"vk-worker/internal/logger"
	"vk-worker/internal/server"
)

func main() {
	cfg := config.MustLoad("config.yaml")
	logger.Setup(cfg.LogLevel)

	srv := server.New(&cfg)
	if err := srv.Start(); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	if err := srv.Stop(5 * time.Second); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}

	slog.Info("Server exited")
}
