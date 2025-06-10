package api

import (
	"context"
	"net/http"
	"strconv"

	"log/slog"
	"vk-worker/internal/service/workermanager"

	"github.com/gin-gonic/gin"
)

type API struct {
	manager workermanager.WorkerManager
	ctx     context.Context
}

func NewAPI(ctx context.Context, manager workermanager.WorkerManager) *API {
	return &API{manager: manager, ctx: ctx}
}

type SendRequest struct {
	Messages []string `json:"messages"`
}

func (api *API) AddWorkers(c *gin.Context) {
	countStr := c.Param("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid count"})
		return
	}

	for i := 0; i < count; i++ {
		api.manager.AddWorkerWithContext(api.ctx)
	}
	slog.Info("Added workers via API", "count", count)
	c.JSON(http.StatusOK, gin.H{"message": "workers added", "count": count})
}

func (api *API) RemoveWorkers(c *gin.Context) {
	countStr := c.Param("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid count"})
		return
	}

	removed := 0
	for i := 0; i < count; i++ {
		if api.manager.RemoveWorker() {
			removed++
		} else {
			break
		}
	}
	slog.Info("Removed workers via API", "requested", count, "removed", removed)
	c.JSON(http.StatusOK, gin.H{
		"message":             "workers removed",
		"requested_to_remove": count,
		"actually_removed":    removed,
	})
}

func (api *API) SendMessages(c *gin.Context) {
	var req SendRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	sentCount := 0
	for _, msg := range req.Messages {
		if api.manager.Send(msg) {
			sentCount++
		}
	}
	slog.Info("Sent messages via API", "count", sentCount)
	c.JSON(http.StatusOK, gin.H{"sent": sentCount, "total": len(req.Messages)})
}

func (api *API) GetStats(c *gin.Context) {
	stats := api.manager.GetStats()
	c.JSON(http.StatusOK, gin.H{
		"workers":            stats.Workers,
		"queue_length":       stats.QueueLength,
		"messages_processed": stats.MessagesProcessed,
		"messages_total":     stats.MessagesTotal,
	})
}

func (api *API) StopAll(c *gin.Context) {
	api.manager.StopAll()
	slog.Info("Stopped all workers via API")
	c.JSON(http.StatusOK, gin.H{"message": "all workers stopped"})
}
