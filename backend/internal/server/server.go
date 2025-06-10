package server

import (
	"context"
	"net/http"
	"time"

	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"vk-worker/internal/api"
	"vk-worker/internal/config"
	"vk-worker/internal/service/workermanager"
)

type Server struct {
	httpServer *http.Server
	manager    workermanager.WorkerManager
	cancel     context.CancelFunc
}

func New(cfg *config.Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	manager := workermanager.New(cfg.QueueSize)

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

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.ServerPort,
			Handler: r,
		},
		manager: manager,
		cancel:  cancel,
	}
}

func (s *Server) Start() error {
	slog.Info("Starting server", "addr", s.httpServer.Addr)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()
	return nil
}

func (s *Server) Stop(timeout time.Duration) error {
	s.cancel()
	s.manager.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
