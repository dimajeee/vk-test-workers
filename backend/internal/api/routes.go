package api

import "github.com/gin-gonic/gin"

func (api *API) RegisterRoutes(r *gin.Engine) {
	r.POST("/workers/add/:count", api.AddWorkers)
	r.POST("/workers/remove/:count", api.RemoveWorkers)
	r.POST("/send", api.SendMessages)
	r.GET("/stats", api.GetStats)
	r.POST("/stop", api.StopAll)
}
