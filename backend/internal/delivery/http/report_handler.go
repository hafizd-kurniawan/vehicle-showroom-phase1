package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"vehicle-showroom/internal/entity"
	"vehicle-showroom/internal/usecase"
)

type ReportHandler struct {
	reportUsecase usecase.ReportUsecase
}

func NewReportHandler(reportUsecase usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{
		reportUsecase: reportUsecase,
	}
}

func (h *ReportHandler) GetVehicleProfitability(c *gin.Context) {
	var req entity.DateRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range", "message": err.Error()})
		return
	}

	report, err := h.reportUsecase.GetVehicleProfitabilityReport(req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": report})
}

func (h *ReportHandler) GetSalesReport(c *gin.Context) {
	var req entity.DateRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range", "message": err.Error()})
		return
	}

	report, err := h.reportUsecase.GetSalesReport(req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": report})
}

func (h *ReportHandler) GetPurchaseReport(c *gin.Context) {
	var req entity.DateRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date range", "message": err.Error()})
		return
	}

	report, err := h.reportUsecase.GetPurchaseReport(req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": report})
}
