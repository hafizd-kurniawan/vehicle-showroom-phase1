package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"vehicle-showroom/internal/entity"
	"vehicle-showroom/internal/usecase"
)

type RepairHandler struct {
	repairUsecase usecase.RepairUsecase
}

func NewRepairHandler(repairUsecase usecase.RepairUsecase) *RepairHandler {
	return &RepairHandler{
		repairUsecase: repairUsecase,
	}
}

func (h *RepairHandler) Create(c *gin.Context) {
	var req entity.CreateRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userValue, _ := c.Get("user")
	user := userValue.(*entity.User)

	repair, err := h.repairUsecase.Create(&req, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create repair",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    repair,
	})
}

func (h *RepairHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid repair ID",
			"message": "Repair ID must be a number",
		})
		return
	}

	repair, err := h.repairUsecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Repair not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    repair,
	})
}

func (h *RepairHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")

	response, err := h.repairUsecase.List(page, limit, search, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list repairs",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *RepairHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid repair ID",
			"message": "Repair ID must be a number",
		})
		return
	}

	var req entity.UpdateRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	repair, err := h.repairUsecase.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update repair",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    repair,
	})
}

func (h *RepairHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid repair ID",
			"message": "Repair ID must be a number",
		})
		return
	}

	var req entity.UpdateRepairStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	repair, err := h.repairUsecase.UpdateStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update repair status",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    repair,
	})
}

func (h *RepairHandler) AddPart(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid repair ID",
			"message": "Repair ID must be a number",
		})
		return
	}

	var req entity.AddPartToRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userValue, _ := c.Get("user")
	user := userValue.(*entity.User)

	repairPart, err := h.repairUsecase.AddPart(id, &req, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to add part to repair",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    repairPart,
	})
}

func (h *RepairHandler) RemovePart(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid repair ID",
			"message": "Repair ID must be a number",
		})
		return
	}

	partIdStr := c.Param("partId")
	partId, err := strconv.Atoi(partIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid part ID",
			"message": "Part ID must be a number",
		})
		return
	}

	if err := h.repairUsecase.RemovePart(id, partId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to remove part from repair",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Part removed from repair successfully",
	})
}
