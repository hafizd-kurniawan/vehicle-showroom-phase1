package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"vehicle-showroom/internal/entity"
	"vehicle-showroom/internal/usecase"
)

type SparePartHandler struct {
	sparePartUsecase usecase.SparePartUsecase
}

func NewSparePartHandler(sparePartUsecase usecase.SparePartUsecase) *SparePartHandler {
	return &SparePartHandler{
		sparePartUsecase: sparePartUsecase,
	}
}

func (h *SparePartHandler) Create(c *gin.Context) {
	var req entity.CreateSparePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	sparePart, err := h.sparePartUsecase.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create spare part",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    sparePart,
	})
}

func (h *SparePartHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid spare part ID",
			"message": "Spare part ID must be a number",
		})
		return
	}

	sparePart, err := h.sparePartUsecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Spare part not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    sparePart,
	})
}

func (h *SparePartHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	response, err := h.sparePartUsecase.List(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list spare parts",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *SparePartHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid spare part ID",
			"message": "Spare part ID must be a number",
		})
		return
	}

	var req entity.UpdateSparePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	sparePart, err := h.sparePartUsecase.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update spare part",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    sparePart,
	})
}

func (h *SparePartHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid spare part ID",
			"message": "Spare part ID must be a number",
		})
		return
	}

	if err := h.sparePartUsecase.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to delete spare part",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Spare part deleted successfully",
	})
}
