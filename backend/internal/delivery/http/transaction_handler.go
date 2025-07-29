package http

import (
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
  "vehicle-showroom/internal/entity"
  "vehicle-showroom/internal/usecase"
)

type TransactionHandler struct {
  transactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecase) *TransactionHandler {
  return &TransactionHandler{
    transactionUsecase: transactionUsecase,
  }
}

func (h *TransactionHandler) CreatePurchase(c *gin.Context) {
  var req entity.CreatePurchaseTransactionRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Invalid request",
      "message": err.Error(),
    })
    return
  }

  // For now, use a default cashier ID (in real app, get from JWT token)
  cashierID := 2

  transaction, err := h.transactionUsecase.CreatePurchase(&req, cashierID)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Failed to create purchase transaction",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusCreated, gin.H{
    "success": true,
    "data":    transaction,
  })
}

func (h *TransactionHandler) GetPurchaseByID(c *gin.Context) {
  idStr := c.Param("id")
  id, err := strconv.Atoi(idStr)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Invalid transaction ID",
      "message": "Transaction ID must be a number",
    })
    return
  }

  transaction, err := h.transactionUsecase.GetPurchaseByID(id)
  if err != nil {
    c.JSON(http.StatusNotFound, gin.H{
      "error":   "Purchase transaction not found",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data":    transaction,
  })
}

func (h *TransactionHandler) ListPurchases(c *gin.Context) {
  page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
  limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
  search := c.Query("search")

  response, err := h.transactionUsecase.ListPurchases(page, limit, search)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error":   "Failed to list purchase transactions",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data":    response,
  })
}

func (h *TransactionHandler) CreateSales(c *gin.Context) {
  var req entity.CreateSalesTransactionRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Invalid request",
      "message": err.Error(),
    })
    return
  }

  // For now, use a default cashier ID (in real app, get from JWT token)
  cashierID := 2

  transaction, err := h.transactionUsecase.CreateSales(&req, cashierID)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Failed to create sales transaction",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusCreated, gin.H{
    "success": true,
    "data":    transaction,
  })
}

func (h *TransactionHandler) GetSalesByID(c *gin.Context) {
  idStr := c.Param("id")
  id, err := strconv.Atoi(idStr)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Invalid transaction ID",
      "message": "Transaction ID must be a number",
    })
    return
  }

  transaction, err := h.transactionUsecase.GetSalesByID(id)
  if err != nil {
    c.JSON(http.StatusNotFound, gin.H{
      "error":   "Sales transaction not found",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data":    transaction,
  })
}

func (h *TransactionHandler) ListSales(c *gin.Context) {
  page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
  limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
  search := c.Query("search")

  response, err := h.transactionUsecase.ListSales(page, limit, search)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error":   "Failed to list sales transactions",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data":    response,
  })
}

func (h *TransactionHandler) GetDashboardStats(c *gin.Context) {
  stats, err := h.transactionUsecase.GetDashboardStats()
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error":   "Failed to get dashboard stats",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data":    stats,
  })
}
