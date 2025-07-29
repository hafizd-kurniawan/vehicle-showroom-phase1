package usecase

import (
	"fmt"
	"time"

	"vehicle-showroom/internal/entity"
	"vehicle-showroom/internal/repository"
)

type ReportUsecase interface {
	GetVehicleProfitabilityReport(startDate, endDate string) ([]entity.VehicleProfitability, error)
	GetSalesReport(startDate, endDate string) ([]entity.SalesTransaction, error)
	GetPurchaseReport(startDate, endDate string) ([]entity.PurchaseTransaction, error)
}

type reportUsecase struct {
	reportRepo repository.ReportRepository
}

func NewReportUsecase(reportRepo repository.ReportRepository) ReportUsecase {
	return &reportUsecase{
		reportRepo: reportRepo,
	}
}

func (u *reportUsecase) GetVehicleProfitabilityReport(startDateStr, endDateStr string) ([]entity.VehicleProfitability, error) {
	start, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	end, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}
	// To include the whole end day
	end = end.Add(24*time.Hour - 1*time.Nanosecond)

	return u.reportRepo.GetVehicleProfitabilityReport(start, end)
}

func (u *reportUsecase) GetSalesReport(startDateStr, endDateStr string) ([]entity.SalesTransaction, error) {
	start, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	end, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}
	end = end.Add(24*time.Hour - 1*time.Nanosecond)

	return u.reportRepo.GetSalesReport(start, end)
}

func (u *reportUsecase) GetPurchaseReport(startDateStr, endDateStr string) ([]entity.PurchaseTransaction, error) {
	start, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	end, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}
	end = end.Add(24*time.Hour - 1*time.Nanosecond)

	return u.reportRepo.GetPurchaseReport(start, end)
}
