package usecase

import (
	"fmt"

	"vehicle-showroom/internal/entity"
	"vehicle-showroom/internal/repository"
)

type SparePartUsecase interface {
	Create(req *entity.CreateSparePartRequest) (*entity.SparePart, error)
	GetByID(id int) (*entity.SparePart, error)
	List(page, limit int, search string) (*entity.SparePartListResponse, error)
	Update(id int, req *entity.UpdateSparePartRequest) (*entity.SparePart, error)
	Delete(id int) error
}

type sparePartUsecase struct {
	sparePartRepo repository.SparePartRepository
}

func NewSparePartUsecase(sparePartRepo repository.SparePartRepository) SparePartUsecase {
	return &sparePartUsecase{
		sparePartRepo: sparePartRepo,
	}
}

func (u *sparePartUsecase) Create(req *entity.CreateSparePartRequest) (*entity.SparePart, error) {
	// Generate part code
	partCode, err := u.sparePartRepo.GeneratePartCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate part code: %w", err)
	}

	sparePart := &entity.SparePart{
		PartCode:      partCode,
		Name:          req.Name,
		Description:   req.Description,
		Brand:         req.Brand,
		CostPrice:     req.CostPrice,
		SellingPrice:  req.SellingPrice,
		StockQuantity: req.StockQuantity,
		MinStockLevel: req.MinStockLevel,
		UnitMeasure:   req.UnitMeasure,
		IsActive:      true,
	}

	if err := u.sparePartRepo.Create(sparePart); err != nil {
		return nil, fmt.Errorf("failed to create spare part: %w", err)
	}

	return sparePart, nil
}

func (u *sparePartUsecase) GetByID(id int) (*entity.SparePart, error) {
	sparePart, err := u.sparePartRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get spare part: %w", err)
	}

	if sparePart == nil {
		return nil, fmt.Errorf("spare part not found")
	}

	return sparePart, nil
}

func (u *sparePartUsecase) List(page, limit int, search string) (*entity.SparePartListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	spareParts, total, err := u.sparePartRepo.List(page, limit, search)
	if err != nil {
		return nil, fmt.Errorf("failed to list spare parts: %w", err)
	}

	return &entity.SparePartListResponse{
		SpareParts: spareParts,
		Total:      total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (u *sparePartUsecase) Update(id int, req *entity.UpdateSparePartRequest) (*entity.SparePart, error) {
	sparePart, err := u.sparePartRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get spare part: %w", err)
	}

	if sparePart == nil {
		return nil, fmt.Errorf("spare part not found")
	}

	sparePart.Name = req.Name
	sparePart.Description = req.Description
	sparePart.Brand = req.Brand
	sparePart.CostPrice = req.CostPrice
	sparePart.SellingPrice = req.SellingPrice
	sparePart.MinStockLevel = req.MinStockLevel
	sparePart.UnitMeasure = req.UnitMeasure
	if req.IsActive != nil {
		sparePart.IsActive = *req.IsActive
	}

	if err := u.sparePartRepo.Update(sparePart); err != nil {
		return nil, fmt.Errorf("failed to update spare part: %w", err)
	}

	return sparePart, nil
}

func (u *sparePartUsecase) Delete(id int) error {
	sparePart, err := u.sparePartRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get spare part: %w", err)
	}

	if sparePart == nil {
		return fmt.Errorf("spare part not found")
	}

	if err := u.sparePartRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete spare part: %w", err)
	}

	return nil
}
