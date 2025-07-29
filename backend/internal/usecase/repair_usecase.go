package usecase

import (
	"fmt"

	"vehicle-showroom/internal/entity"
	"vehicle-showroom/internal/repository"
)

type RepairUsecase interface {
	Create(req *entity.CreateRepairRequest, createdBy int) (*entity.Repair, error)
	GetByID(id int) (*entity.Repair, error)
	List(page, limit int, search, status string) (*entity.RepairListResponse, error)
	Update(id int, req *entity.UpdateRepairRequest) (*entity.Repair, error)
	UpdateStatus(id int, status string) (*entity.Repair, error)
	AddPart(repairId int, req *entity.AddPartToRepairRequest, processedBy int) (*entity.RepairPart, error)
	RemovePart(repairId, partId int) error
}

type repairUsecase struct {
	repairRepo    repository.RepairRepository
	vehicleRepo   repository.VehicleRepository
	sparePartRepo repository.SparePartRepository
}

func NewRepairUsecase(
	repairRepo repository.RepairRepository,
	vehicleRepo repository.VehicleRepository,
	sparePartRepo repository.SparePartRepository,
) RepairUsecase {
	return &repairUsecase{
		repairRepo:    repairRepo,
		vehicleRepo:   vehicleRepo,
		sparePartRepo: sparePartRepo,
	}
}

func (u *repairUsecase) Create(req *entity.CreateRepairRequest, createdBy int) (*entity.Repair, error) {
	// Validate vehicle exists
	vehicle, err := u.vehicleRepo.GetByID(req.VehicleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}
	if vehicle == nil {
		return nil, fmt.Errorf("vehicle not found")
	}

	// Generate repair number
	repairNumber, err := u.repairRepo.GenerateRepairNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate repair number: %w", err)
	}

	repair := &entity.Repair{
		RepairNumber:   repairNumber,
		VehicleID:      req.VehicleID,
		Title:          req.Title,
		Description:    req.Description,
		LaborCost:      0,
		TotalPartsCost: 0,
		TotalCost:      0,
		Status:         "pending",
		MechanicID:     req.MechanicID,
	}

	if err := u.repairRepo.Create(repair); err != nil {
		return nil, fmt.Errorf("failed to create repair: %w", err)
	}

	// Update vehicle status to in_repair
	if err := u.vehicleRepo.UpdateStatus(req.VehicleID, "in_repair"); err != nil {
		return nil, fmt.Errorf("failed to update vehicle status: %w", err)
	}

	// Get the created repair with related data
	return u.repairRepo.GetByID(repair.ID)
}

func (u *repairUsecase) GetByID(id int) (*entity.Repair, error) {
	repair, err := u.repairRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get repair: %w", err)
	}

	if repair == nil {
		return nil, fmt.Errorf("repair not found")
	}

	return repair, nil
}

func (u *repairUsecase) List(page, limit int, search, status string) (*entity.RepairListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	repairs, total, err := u.repairRepo.List(page, limit, search, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list repairs: %w", err)
	}

	return &entity.RepairListResponse{
		Repairs: repairs,
		Total:   total,
		Page:    page,
		Limit:   limit,
	}, nil
}

func (u *repairUsecase) Update(id int, req *entity.UpdateRepairRequest) (*entity.Repair, error) {
	repair, err := u.repairRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get repair: %w", err)
	}

	if repair == nil {
		return nil, fmt.Errorf("repair not found")
	}

	repair.Title = req.Title
	repair.Description = req.Description
	if req.LaborCost != nil {
		repair.LaborCost = *req.LaborCost
	}
	repair.MechanicID = req.MechanicID
	repair.WorkNotes = req.WorkNotes

	if err := u.repairRepo.Update(repair); err != nil {
		return nil, fmt.Errorf("failed to update repair: %w", err)
	}

	// Update repair costs
	if err := u.repairRepo.UpdateRepairCosts(id); err != nil {
		return nil, fmt.Errorf("failed to update repair costs: %w", err)
	}

	return u.repairRepo.GetByID(id)
}

func (u *repairUsecase) UpdateStatus(id int, status string) (*entity.Repair, error) {
	repair, err := u.repairRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get repair: %w", err)
	}

	if repair == nil {
		return nil, fmt.Errorf("repair not found")
	}

	if err := u.repairRepo.UpdateStatus(id, status); err != nil {
		return nil, fmt.Errorf("failed to update repair status: %w", err)
	}

	// If repair is completed, update vehicle status and total repair cost
	if status == "completed" {
		// Update vehicle's total repair cost
		vehicle, err := u.vehicleRepo.GetByID(repair.VehicleID)
		if err == nil && vehicle != nil {
			// Get updated repair with costs
			updatedRepair, err := u.repairRepo.GetByID(id)
			if err == nil {
				vehicle.TotalRepairCost += updatedRepair.TotalCost
				u.vehicleRepo.Update(vehicle)
			}
		}

		// Update vehicle status to ready_to_sell
		u.vehicleRepo.UpdateStatus(repair.VehicleID, "ready_to_sell")
	}

	return u.repairRepo.GetByID(id)
}

func (u *repairUsecase) AddPart(repairId int, req *entity.AddPartToRepairRequest, processedBy int) (*entity.RepairPart, error) {
	// Validate repair exists
	repair, err := u.repairRepo.GetByID(repairId)
	if err != nil {
		return nil, fmt.Errorf("failed to get repair: %w", err)
	}
	if repair == nil {
		return nil, fmt.Errorf("repair not found")
	}

	// Validate spare part exists and has sufficient stock
	sparePart, err := u.sparePartRepo.GetByID(req.SparePartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spare part: %w", err)
	}
	if sparePart == nil {
		return nil, fmt.Errorf("spare part not found")
	}

	if sparePart.StockQuantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: available %d, requested %d", sparePart.StockQuantity, req.Quantity)
	}

	// Calculate costs
	unitCost := sparePart.CostPrice
	totalCost := unitCost * float64(req.Quantity)

	repairPart := &entity.RepairPart{
		RepairID:     repairId,
		SparePartID:  req.SparePartID,
		QuantityUsed: req.Quantity,
		UnitCost:     unitCost,
		TotalCost:    totalCost,
	}

	if err := u.repairRepo.AddPart(repairPart); err != nil {
		return nil, fmt.Errorf("failed to add part to repair: %w", err)
	}

	// Update spare part stock
	newStock := sparePart.StockQuantity - req.Quantity
	if err := u.sparePartRepo.UpdateStock(req.SparePartID, newStock); err != nil {
		return nil, fmt.Errorf("failed to update spare part stock: %w", err)
	}

	// Update repair costs
	if err := u.repairRepo.UpdateRepairCosts(repairId); err != nil {
		return nil, fmt.Errorf("failed to update repair costs: %w", err)
	}

	// TODO: Create stock movement record

	return repairPart, nil
}

func (u *repairUsecase) RemovePart(repairId, partId int) error {
	// Validate repair exists
	repair, err := u.repairRepo.GetByID(repairId)
	if err != nil {
		return fmt.Errorf("failed to get repair: %w", err)
	}
	if repair == nil {
		return fmt.Errorf("repair not found")
	}

	// Get repair part details before removing
	parts, err := u.repairRepo.GetRepairParts(repairId)
	if err != nil {
		return fmt.Errorf("failed to get repair parts: %w", err)
	}

	var partToRemove *entity.RepairPart
	for _, part := range parts {
		if part.ID == partId {
			partToRemove = &part
			break
		}
	}

	if partToRemove == nil {
		return fmt.Errorf("repair part not found")
	}

	// Remove part from repair
	if err := u.repairRepo.RemovePart(repairId, partId); err != nil {
		return fmt.Errorf("failed to remove part from repair: %w", err)
	}

	// Restore spare part stock
	sparePart, err := u.sparePartRepo.GetByID(partToRemove.SparePartID)
	if err == nil && sparePart != nil {
		newStock := sparePart.StockQuantity + partToRemove.QuantityUsed
		u.sparePartRepo.UpdateStock(partToRemove.SparePartID, newStock)
	}

	// Update repair costs
	if err := u.repairRepo.UpdateRepairCosts(repairId); err != nil {
		return fmt.Errorf("failed to update repair costs: %w", err)
	}

	return nil
}
