package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maximicciullo/personal-finance-api/internal/services"
)

type ReportController struct {
	service services.ReportService
}

func NewReportController(service services.ReportService) *ReportController {
	return &ReportController{
		service: service,
	}
}

func (c *ReportController) GetMonthlyReport(ctx *gin.Context) {
	yearParam := ctx.Param("year")
	monthParam := ctx.Param("month")

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid year format",
			"status":  http.StatusBadRequest,
		})
		return
	}

	month, err := strconv.Atoi(monthParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid month format",
			"status":  http.StatusBadRequest,
		})
		return
	}

	report, err := c.service.GetMonthlyReport(year, month)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

func (c *ReportController) GetCurrentMonthReport(ctx *gin.Context) {
	report, err := c.service.GetCurrentMonthReport()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to generate current month report",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, report)
}