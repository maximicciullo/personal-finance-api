package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maximicciullo/personal-finance-api/internal/middleware"
	"github.com/maximicciullo/personal-finance-api/internal/models"
	"github.com/maximicciullo/personal-finance-api/internal/services"
	"go.uber.org/zap"
)

type TransactionController struct {
	service services.TransactionService
	logger  *middleware.BusinessLoggerInstance
}

func NewTransactionController(service services.TransactionService) *TransactionController {
	return &TransactionController{
		service: service,
		logger:  middleware.BusinessLogger(),
	}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	c.logger.Controller("CreateTransaction started",
		zap.String("client_ip", ctx.ClientIP()),
		zap.String("user_agent", ctx.Request.UserAgent()),
	)

	var req models.CreateTransactionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Error("controller", "CreateTransaction - JSON binding failed", err,
			zap.Any("request_body", req),
		)
		
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.logger.Controller("CreateTransaction - request validated",
		zap.String("type", req.Type),
		zap.Float64("amount", req.Amount),
		zap.String("currency", req.Currency),
		zap.String("category", req.Category),
	)

	start := time.Now()
	transaction, err := c.service.CreateTransaction(&req)
	duration := time.Since(start)

	c.logger.Performance("CreateTransaction service call", duration,
		zap.Bool("success", err == nil),
	)

	if err != nil {
		c.logger.Error("controller", "CreateTransaction - service error", err,
			zap.Any("request", req),
		)
		
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.logger.Controller("CreateTransaction completed successfully",
		zap.Int("transaction_id", transaction.ID),
		zap.Duration("total_duration", duration),
	)

	ctx.JSON(http.StatusCreated, transaction)
}

func (c *TransactionController) GetTransactions(ctx *gin.Context) {
	c.logger.Controller("GetTransactions started",
		zap.String("query_params", ctx.Request.URL.RawQuery),
	)

	filters := c.parseFilters(ctx)
	
	c.logger.Controller("GetTransactions - filters parsed",
		zap.Any("filters", filters),
	)

	start := time.Now()
	transactions, err := c.service.GetTransactions(filters)
	duration := time.Since(start)

	c.logger.Performance("GetTransactions service call", duration,
		zap.Int("transaction_count", len(transactions)),
		zap.Bool("success", err == nil),
	)

	if err != nil {
		c.logger.Error("controller", "GetTransactions - service error", err,
			zap.Any("filters", filters),
		)
		
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve transactions",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.logger.Controller("GetTransactions completed successfully",
		zap.Int("transaction_count", len(transactions)),
		zap.Duration("total_duration", duration),
	)

	ctx.JSON(http.StatusOK, transactions)
}

func (c *TransactionController) GetTransaction(ctx *gin.Context) {
	idParam := ctx.Param("id")
	
	c.logger.Controller("GetTransaction started",
		zap.String("transaction_id", idParam),
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.logger.Error("controller", "GetTransaction - invalid ID format", err,
			zap.String("id_param", idParam),
		)
		
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid transaction ID",
			"status":  http.StatusBadRequest,
		})
		return
	}

	start := time.Now()
	transaction, err := c.service.GetTransaction(id)
	duration := time.Since(start)

	c.logger.Performance("GetTransaction service call", duration,
		zap.Int("transaction_id", id),
		zap.Bool("success", err == nil),
	)

	if err != nil {
		c.logger.Error("controller", "GetTransaction - service error", err,
			zap.Int("transaction_id", id),
		)
		
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Transaction not found",
			"status":  http.StatusNotFound,
		})
		return
	}

	c.logger.Controller("GetTransaction completed successfully",
		zap.Int("transaction_id", id),
		zap.Duration("total_duration", duration),
	)

	ctx.JSON(http.StatusOK, transaction)
}

func (c *TransactionController) DeleteTransaction(ctx *gin.Context) {
	idParam := ctx.Param("id")
	
	c.logger.Controller("DeleteTransaction started",
		zap.String("transaction_id", idParam),
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.logger.Error("controller", "DeleteTransaction - invalid ID format", err,
			zap.String("id_param", idParam),
		)
		
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid transaction ID",
			"status":  http.StatusBadRequest,
		})
		return
	}

	start := time.Now()
	err = c.service.DeleteTransaction(id)
	duration := time.Since(start)

	c.logger.Performance("DeleteTransaction service call", duration,
		zap.Int("transaction_id", id),
		zap.Bool("success", err == nil),
	)

	if err != nil {
		c.logger.Error("controller", "DeleteTransaction - service error", err,
			zap.Int("transaction_id", id),
		)
		
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Transaction not found",
			"status":  http.StatusNotFound,
		})
		return
	}

	c.logger.Controller("DeleteTransaction completed successfully",
		zap.Int("transaction_id", id),
		zap.Duration("total_duration", duration),
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Transaction deleted successfully",
	})
}

func (c *TransactionController) UpdateTransaction(ctx *gin.Context) {
	idParam := ctx.Param("id")
	
	c.logger.Controller("UpdateTransaction started",
		zap.String("transaction_id", idParam),
	)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.logger.Error("controller", "UpdateTransaction - invalid ID format", err,
			zap.String("id_param", idParam),
		)
		
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid transaction ID",
			"status":  http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateTransactionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Error("controller", "UpdateTransaction - JSON binding failed", err,
			zap.Any("request_body", req),
		)
		
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.logger.Controller("UpdateTransaction - request validated",
		zap.Int("transaction_id", id),
		zap.Any("update_request", req),
	)

	start := time.Now()
	transaction, err := c.service.UpdateTransaction(id, &req)
	duration := time.Since(start)

	c.logger.Performance("UpdateTransaction service call", duration,
		zap.Int("transaction_id", id),
		zap.Bool("success", err == nil),
	)

	if err != nil {
		c.logger.Error("controller", "UpdateTransaction - service error", err,
			zap.Int("transaction_id", id),
			zap.Any("request", req),
		)
		
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Transaction not found",
			"status":  http.StatusNotFound,
		})
		return
	}

	c.logger.Controller("UpdateTransaction completed successfully",
		zap.Int("transaction_id", id),
		zap.Duration("total_duration", duration),
	)

	ctx.JSON(http.StatusOK, transaction)
}

func (c *TransactionController) parseFilters(ctx *gin.Context) models.TransactionFilters {
	filters := models.TransactionFilters{
		Type:     ctx.Query("type"),
		Category: ctx.Query("category"),
		Currency: ctx.Query("currency"),
	}

	c.logger.Debug("controller", "Parsing query filters",
		zap.String("type", filters.Type),
		zap.String("category", filters.Category),
		zap.String("currency", filters.Currency),
	)

	// Parse date filters if provided
	if fromDateStr := ctx.Query("from_date"); fromDateStr != "" {
		if fromDate, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			filters.FromDate = &fromDate
			c.logger.Debug("controller", "Parsed from_date filter",
				zap.Time("from_date", fromDate),
			)
		} else {
			c.logger.Error("controller", "Invalid from_date format", err,
				zap.String("from_date_str", fromDateStr),
			)
		}
	}

	if toDateStr := ctx.Query("to_date"); toDateStr != "" {
		if toDate, err := time.Parse("2006-01-02", toDateStr); err == nil {
			filters.ToDate = &toDate
			c.logger.Debug("controller", "Parsed to_date filter",
				zap.Time("to_date", toDate),
			)
		} else {
			c.logger.Error("controller", "Invalid to_date format", err,
				zap.String("to_date_str", toDateStr),
			)
		}
	}

	return filters
}