package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maximicciullo/personal-finance-api/internal/models"
	"github.com/maximicciullo/personal-finance-api/internal/services"
)

type TransactionController struct {
	service services.TransactionService
}

func NewTransactionController(service services.TransactionService) *TransactionController {
	return &TransactionController{
		service: service,
	}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var req models.CreateTransactionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	transaction, err := c.service.CreateTransaction(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	ctx.JSON(http.StatusCreated, transaction)
}

func (c *TransactionController) GetTransactions(ctx *gin.Context) {
	filters := c.parseFilters(ctx)

	transactions, err := c.service.GetTransactions(filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve transactions",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

func (c *TransactionController) GetTransaction(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid transaction ID",
			"status":  http.StatusBadRequest,
		})
		return
	}

	transaction, err := c.service.GetTransaction(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Transaction not found",
			"status":  http.StatusNotFound,
		})
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

func (c *TransactionController) DeleteTransaction(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid transaction ID",
			"status":  http.StatusBadRequest,
		})
		return
	}

	err = c.service.DeleteTransaction(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Transaction not found",
			"status":  http.StatusNotFound,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction deleted successfully",
	})
}

func (c *TransactionController) parseFilters(ctx *gin.Context) models.TransactionFilters {
	filters := models.TransactionFilters{
		Type:     ctx.Query("type"),
		Category: ctx.Query("category"),
		Currency: ctx.Query("currency"),
	}

	// Parse date filters if provided
	if fromDateStr := ctx.Query("from_date"); fromDateStr != "" {
		if fromDate, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			filters.FromDate = &fromDate
		}
	}

	if toDateStr := ctx.Query("to_date"); toDateStr != "" {
		if toDate, err := time.Parse("2006-01-02", toDateStr); err == nil {
			filters.ToDate = &toDate
		}
	}

	return filters
}