package usecase

import (
  "fmt"
  "time"

  "vehicle-showroom/internal/entity"
  "vehicle-showroom/internal/repository"
)

type VehicleUsecase interface {
  Create(req *entity.CreateVehicleRequest, purchasedBy int) (*entity.Vehicle, error)
  GetByID(id int) (*entity.Vehicle, error)
  List(page, limit int, search, status string) (*entity.VehicleListResponse, error)
  Update(id int, req *entity.UpdateVehicleRequest) (*entity.Vehicle, error)
  UpdateStatus(id int, status string) (*entity.Vehicle, error)
  Delete(id int) error
}

type vehicleUsecase struct {
  vehicleRepo  repository.VehicleRepository
  customerRepo repository.CustomerRepository
}

func NewVehicleUsecase(vehicleRepo repository.VehicleRepository, customerRepo repository.CustomerRepository) VehicleUsecase {
  return &vehicleUsecase{
    vehicleRepo:  vehicleRepo,
    customerRepo: customerRepo,
  }
}

func (u *vehicleUsecase) Create(req *entity.CreateVehicleRequest, purchasedBy int) (*entity.Vehicle, error) {
  // Validate customer if provided
  if req.PurchasedFromCustomerID != nil {
    customer, err := u.customerRepo.GetByID(*req.PurchasedFromCustomerID)
    if err != nil {
      return nil, fmt.Errorf("failed to validate customer: %w", err)
    }
    if customer == nil {
      return nil, fmt.Errorf("customer not found")
    }
  }
  
  // Generate vehicle code
  vehicleCode, err := u.vehicleRepo.GenerateVehicleCode()
  if err != nil {
    return nil, fmt.Errorf("failed to generate vehicle code: %w", err)
  }
  
  now := time.Now()
  vehicle := &entity.Vehicle{
    VehicleCode:             vehicleCode,
    ChassisNumber:           req.ChassisNumber,
    LicensePlate:            req.LicensePlate,
    Brand:                   req.Brand,
    Model:                   req.Model,
    Variant:                 req.Variant,
    Year:                    req.Year,
    Color:                   req.Color,
    Mileage:                 req.Mileage,
    FuelType:                req.FuelType,
    Transmission:            req.Transmission,
    PurchasePrice:           req.PurchasePrice,
    TotalRepairCost:         0,
    Status:                  "purchased",
    PurchasedFromCustomerID: req.PurchasedFromCustomerID,
    PurchasedByCashier:      &purchasedBy,
    PurchasedAt:             &now,
    PurchaseNotes:           req.PurchaseNotes,
    ConditionNotes:          req.ConditionNotes,
  }
  
  if err := u.vehicleRepo.Create(vehicle); err != nil {
    return nil, fmt.Errorf("failed to create vehicle: %w", err)
  }
  
  // Get the created vehicle with related data
  return u.vehicleRepo.GetByID(vehicle.ID)
}

func (u *vehicleUsecase) GetByID(id int) (*entity.Vehicle, error) {
  vehicle, err := u.vehicleRepo.GetByID(id)
  if err != nil {
    return nil, fmt.Errorf("failed to get vehicle: %w", err)
  }
  
  if vehicle == nil {
    return nil, fmt.Errorf("vehicle not found")
  }
  
  return vehicle, nil
}

func (u *vehicleUsecase) List(page, limit int, search, status string) (*entity.VehicleListResponse, error) {
  if page <= 0 {
    page = 1
  }
  if limit <= 0 || limit > 100 {
    limit = 10
  }
  
  vehicles, total, err := u.vehicleRepo.List(page, limit, search, status)
  if err != nil {
    return nil, fmt.Errorf("failed to list vehicles: %w", err)
  }
  
  return &entity.VehicleListResponse{
    Vehicles: vehicles,
    Total:    total,
    Page:     page,
    Limit:    limit,
  }, nil
}

func (u *vehicleUsecase) Update(id int, req *entity.UpdateVehicleRequest) (*entity.Vehicle, error) {
  vehicle, err := u.vehicleRepo.GetByID(id)
  if err != nil {
    return nil, fmt.Errorf("failed to get vehicle: %w", err)
  }
  
  if vehicle == nil {
    return nil, fmt.Errorf("vehicle not found")
  }
  
  vehicle.LicensePlate = req.LicensePlate
  vehicle.Brand = req.Brand
  vehicle.Model = req.Model
  vehicle.Variant = req.Variant
  vehicle.Year = req.Year
  vehicle.Color = req.Color
  vehicle.Mileage = req.Mileage
  vehicle.FuelType = req.FuelType
  vehicle.Transmission = req.Transmission
  vehicle.SuggestedSellingPrice = req.SuggestedSellingPrice
  vehicle.PurchaseNotes = req.PurchaseNotes
  vehicle.ConditionNotes = req.ConditionNotes
  
  if err := u.vehicleRepo.Update(vehicle); err != nil {
    return nil, fmt.Errorf("failed to update vehicle: %w", err)
  }
  
  return u.vehicleRepo.GetByID(id)
}

func (u *vehicleUsecase) UpdateStatus(id int, status string) (*entity.Vehicle, error) {
  vehicle, err := u.vehicleRepo.GetByID(id)
  if err != nil {
    return nil, fmt.Errorf("failed to get vehicle: %w", err)
  }
  
  if vehicle == nil {
    return nil, fmt.Errorf("vehicle not found")
  }
  
  if err := u.vehicleRepo.UpdateStatus(id, status); err != nil {
    return nil, fmt.Errorf("failed to update vehicle status: %w", err)
  }
  
  return u.vehicleRepo.GetByID(id)
}

func (u *vehicleUsecase) Delete(id int) error {
  vehicle, err := u.vehicleRepo.GetByID(id)
  if err != nil {
    return fmt.Errorf("failed to get vehicle: %w", err)
  }
  
  if vehicle == nil {
    return fmt.Errorf("vehicle not found")
  }
  
  if err := u.vehicleRepo.Delete(id); err != nil {
    return fmt.Errorf("failed to delete vehicle: %w", err)
  }
  
  return nil
}
