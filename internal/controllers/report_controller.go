package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maximicciullo/personal-finance-api/internal/middleware"
	"github.com/maximicciullo/personal-finance-api/internal/services"
	"go.uber.org/zap"
)

type ReportController struct {
	service services.ReportService
	logger  *middleware.BusinessLoggerInstance
}

func NewReportController(service services.ReportService) *ReportController {
	return &ReportController{
		service: service,
		logger:  middleware.BusinessLogger(),
	}
}

func (c *ReportController) GetMonthlyReport(ctx *gin.Context) {
	yearParam := ctx.Param("year")
	monthParam := ctx.Param("month")

	c.logger.Controller("GetMonthlyReport started",
		zap.String("year_param", yearParam),
		zap.String("month_param", monthParam),
		zap.String("client_ip", ctx.ClientIP()),
	)

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.logger.Error("controller", "GetMonthlyReport - invalid year format", err,
			zap.String("year_param", yearParam),
		)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid year format",
			"status":  http.StatusBadRequest,
		})
		return
	}

	month, err := strconv.Atoi(monthParam)
	if err != nil {
		c.logger.Error("controller", "GetMonthlyReport - invalid month format", err,
			zap.String("month_param", monthParam),
		)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid month format",
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.logger.Controller("GetMonthlyReport - parameters validated",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	start := time.Now()
	report, err := c.service.GetMonthlyReport(year, month)
	duration := time.Since(start)

	c.logger.Performance("GetMonthlyReport service call", duration,
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Bool("success", err == nil),
	)

	if err != nil {
		c.logger.Error("controller", "GetMonthlyReport - service error", err,
			zap.Int("year", year),
			zap.Int("month", month),
		)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.logger.Controller("GetMonthlyReport completed successfully",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("transaction_count", report.Summary.TransactionCount),
		zap.Duration("total_duration", duration),
	)

	ctx.JSON(http.StatusOK, report)
}

func (c *ReportController) GetCurrentMonthReport(ctx *gin.Context) {
	now := time.Now()
	c.logger.Controller("GetCurrentMonthReport started",
		zap.Int("current_year", now.Year()),
		zap.Int("current_month", int(now.Month())),
		zap.String("client_ip", ctx.ClientIP()),
	)

	start := time.Now()
	report, err := c.service.GetCurrentMonthReport()
	duration := time.Since(start)

	c.logger.Performance("GetCurrentMonthReport service call", duration,
		zap.Bool("success", err == nil),
	)

	if err != nil {
		c.logger.Error("controller", "GetCurrentMonthReport - service error", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to generate current month report",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.logger.Controller("GetCurrentMonthReport completed successfully",
		zap.String("month", report.Month),
		zap.Int("year", report.Year),
		zap.Int("transaction_count", report.Summary.TransactionCount),
		zap.Duration("total_duration", duration),
	)

	ctx.JSON(http.StatusOK, report)
}