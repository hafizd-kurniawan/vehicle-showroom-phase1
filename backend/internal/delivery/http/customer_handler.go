package http

import (
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
  "vehicle-showroom/internal/entity"
  "vehicle-showroom/internal/usecase"
)

type CustomerHandler struct {
  customerUsecase usecase.CustomerUsecase
}

func NewCustomerHandler(customerUsecase usecase.CustomerUsecase) *CustomerHandler {
  return &CustomerHandler{
    customerUsecase: customerUsecase,
  }
}

func (h *CustomerHandler) Create(c *gin.Context) {
  var req entity.CreateCustomerRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Invalid request",
      "message": err.Error(),
    })
    return
  }

  // For now, use a default user ID (in real app, get from JWT token)
  createdBy := 1

  customer, err := h.customerUsecase.Create(&req, createdBy)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Failed to create customer",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusCreated, gin.H{
    "success": true,
    "data":    customer,
  })
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
  idStr := c.Param("id")
  id, err := strconv.Atoi(idStr)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Invalid customer ID",
      "message": "Customer ID must be a number",
    })
    return
  }

  customer, err := h.customerUsecase.GetByID(id)
  if err != nil {
    c.JSON(http.StatusNotFound, gin.H{
      "error":   "Customer not found",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data":    customer,
  })
}

func (h *CustomerHandler) List(c *gin.Context) {
  page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
  limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
  search := c.Query("search")

  response, err := h.customerUsecase.List(page, limit, search)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error":   "Failed to list customers",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data":    response,
  })
}

func (h *CustomerHandler) Update(c *gin.Context) {
  idStr := c.Param("id")
  id, err := strconv.Atoi(idStr)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Invalid customer ID",
      "message": "Customer ID must be a number",
    })
    return
  }

  var req entity.UpdateCustomerRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Invalid request",
      "message": err.Error(),
    })
    return
  }

  customer, err := h.customerUsecase.Update(id, &req)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Failed to update customer",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "data":    customer,
  })
}

func (h *CustomerHandler) Delete(c *gin.Context) {
  idStr := c.Param("id")
  id, err := strconv.Atoi(idStr)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Invalid customer ID",
      "message": "Customer ID must be a number",
    })
    return
  }

  if err := h.customerUsecase.Delete(id); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error":   "Failed to delete customer",
      "message": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "success": true,
    "message": "Customer deleted successfully",
  })
}
