package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maximicciullo/personal-finance-api/internal/middleware"
	"go.uber.org/zap"
)

type HealthController struct {
	logger *middleware.BusinessLoggerInstance
}

func NewHealthController() *HealthController {
	return &HealthController{
		logger: middleware.BusinessLogger(),
	}
}

func (c *HealthController) HealthCheck(ctx *gin.Context) {
	c.logger.Controller("HealthCheck started",
		zap.String("client_ip", ctx.ClientIP()),
		zap.String("user_agent", ctx.Request.UserAgent()),
	)

	start := time.Now()
	
	response := gin.H{
		"status":    "healthy",
		"service":   "personal-finance-api",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    "running",
	}

	duration := time.Since(start)
	c.logger.Performance("HealthCheck response preparation", duration)

	c.logger.Controller("HealthCheck completed successfully",
		zap.Duration("duration", duration),
		zap.String("status", "healthy"),
	)

	ctx.JSON(http.StatusOK, response)
}