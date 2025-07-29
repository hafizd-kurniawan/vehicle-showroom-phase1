package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"vehicle-showroom/internal/entity"
	"vehicle-showroom/internal/usecase"
)

type VehicleHandler struct {
	vehicleUsecase usecase.VehicleUsecase
}

func NewVehicleHandler(vehicleUsecase usecase.VehicleUsecase) *VehicleHandler {
	return &VehicleHandler{
		vehicleUsecase: vehicleUsecase,
	}
}

func (h *VehicleHandler) Create(c *gin.Context) {
	var req entity.CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userValue, _ := c.Get("user")
	user := userValue.(*entity.User)
	purchasedBy := user.ID

	vehicle, err := h.vehicleUsecase.Create(&req, purchasedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create vehicle",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    vehicle,
	})
}

func (h *VehicleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid vehicle ID",
			"message": "Vehicle ID must be a number",
		})
		return
	}

	vehicle, err := h.vehicleUsecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Vehicle not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vehicle,
	})
}

func (h *VehicleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")

	response, err := h.vehicleUsecase.List(page, limit, search, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list vehicles",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (h *VehicleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid vehicle ID",
			"message": "Vehicle ID must be a number",
		})
		return
	}

	var req entity.UpdateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	vehicle, err := h.vehicleUsecase.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update vehicle",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vehicle,
	})
}

func (h *VehicleHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid vehicle ID",
			"message": "Vehicle ID must be a number",
		})
		return
	}

	var req entity.UpdateVehicleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	vehicle, err := h.vehicleUsecase.UpdateStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update vehicle status",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vehicle,
	})
}

func (h *VehicleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid vehicle ID",
			"message": "Vehicle ID must be a number",
		})
		return
	}

	if err := h.vehicleUsecase.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to delete vehicle",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Vehicle deleted successfully",
	})
}
