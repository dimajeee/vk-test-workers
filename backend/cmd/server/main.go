package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"
	"vk-worker/internal/api"
	"vk-worker/internal/config"
	"vk-worker/internal/logger"
	"vk-worker/internal/service/workermanager"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	logger.Setup(cfg.LogLevel)
	slog.Info("Starting server", "port", cfg.ServerPort)

	manager := workermanager.New(cfg.QueueSize)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < cfg.InitialWorkers; i++ {
		manager.AddWorkerWithContext(ctx)
	}

	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiSrv := api.NewAPI(ctx, manager)
	apiSrv.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	cancel()
	//manager.StopAll()
	manager.Wait()

	ctxShutdown, cancelShatdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShatdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		manager.CloseInput()
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}
