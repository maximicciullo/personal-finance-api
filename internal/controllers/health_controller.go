package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) HealthCheck(ctx *gin.Context) {
	response := gin.H{
		"status":    "healthy",
		"service":   "personal-finance-api",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    "running",
	}

	ctx.JSON(http.StatusOK, response)
}
