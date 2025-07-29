package usecase

import (
  "fmt"

  "vehicle-showroom/internal/entity"
  "vehicle-showroom/internal/repository"
)

type CustomerUsecase interface {
  Create(req *entity.CreateCustomerRequest, createdBy int) (*entity.Customer, error)
  GetByID(id int) (*entity.Customer, error)
  List(page, limit int, search string) (*entity.CustomerListResponse, error)
  Update(id int, req *entity.UpdateCustomerRequest) (*entity.Customer, error)
  Delete(id int) error
}

type customerUsecase struct {
  customerRepo repository.CustomerRepository
}

func NewCustomerUsecase(customerRepo repository.CustomerRepository) CustomerUsecase {
  return &customerUsecase{
    customerRepo: customerRepo,
  }
}

func (u *customerUsecase) Create(req *entity.CreateCustomerRequest, createdBy int) (*entity.Customer, error) {
  // Generate customer code
  customerCode, err := u.customerRepo.GenerateCustomerCode()
  if err != nil {
    return nil, fmt.Errorf("failed to generate customer code: %w", err)
  }
  
  customer := &entity.Customer{
    CustomerCode: customerCode,
    Name:         req.Name,
    Phone:        req.Phone,
    Email:        req.Email,
    Address:      req.Address,
    IDCardNumber: req.IDCardNumber,
    Type:         req.Type,
    CreatedBy:    &createdBy,
    IsActive:     true,
  }
  
  if err := u.customerRepo.Create(customer); err != nil {
    return nil, fmt.Errorf("failed to create customer: %w", err)
  }
  
  return customer, nil
}

func (u *customerUsecase) GetByID(id int) (*entity.Customer, error) {
  customer, err := u.customerRepo.GetByID(id)
  if err != nil {
    return nil, fmt.Errorf("failed to get customer: %w", err)
  }
  
  if customer == nil {
    return nil, fmt.Errorf("customer not found")
  }
  
  return customer, nil
}

func (u *customerUsecase) List(page, limit int, search string) (*entity.CustomerListResponse, error) {
  if page <= 0 {
    page = 1
  }
  if limit <= 0 || limit > 100 {
    limit = 10
  }
  
  customers, total, err := u.customerRepo.List(page, limit, search)
  if err != nil {
    return nil, fmt.Errorf("failed to list customers: %w", err)
  }
  
  return &entity.CustomerListResponse{
    Customers: customers,
    Total:     total,
    Page:      page,
    Limit:     limit,
  }, nil
}

func (u *customerUsecase) Update(id int, req *entity.UpdateCustomerRequest) (*entity.Customer, error) {
  customer, err := u.customerRepo.GetByID(id)
  if err != nil {
    return nil, fmt.Errorf("failed to get customer: %w", err)
  }
  
  if customer == nil {
    return nil, fmt.Errorf("customer not found")
  }
  
  customer.Name = req.Name
  customer.Phone = req.Phone
  customer.Email = req.Email
  customer.Address = req.Address
  customer.IDCardNumber = req.IDCardNumber
  customer.Type = req.Type
  
  if err := u.customerRepo.Update(customer); err != nil {
    return nil, fmt.Errorf("failed to update customer: %w", err)
  }
  
  return customer, nil
}

func (u *customerUsecase) Delete(id int) error {
  customer, err := u.customerRepo.GetByID(id)
  if err != nil {
    return fmt.Errorf("failed to get customer: %w", err)
  }
  
  if customer == nil {
    return fmt.Errorf("customer not found")
  }
  
  if err := u.customerRepo.Delete(id); err != nil {
    return fmt.Errorf("failed to delete customer: %w", err)
  }
  
  return nil
}
