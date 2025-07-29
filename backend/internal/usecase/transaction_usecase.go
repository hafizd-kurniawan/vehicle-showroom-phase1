package usecase

import (
  "fmt"
  "time"

  "vehicle-showroom/internal/entity"
  "vehicle-showroom/internal/repository"
)

type TransactionUsecase interface {
  // Purchase Transactions
  CreatePurchase(req *entity.CreatePurchaseTransactionRequest, cashierID int) (*entity.PurchaseTransaction, error)
  GetPurchaseByID(id int) (*entity.PurchaseTransaction, error)
  ListPurchases(page, limit int, search string) (*entity.TransactionListResponse, error)
  
  // Sales Transactions
  CreateSales(req *entity.CreateSalesTransactionRequest, cashierID int) (*entity.SalesTransaction, error)
  GetSalesByID(id int) (*entity.SalesTransaction, error)
  ListSales(page, limit int, search string) (*entity.TransactionListResponse, error)
  
  // Dashboard
  GetDashboardStats() (*entity.DashboardStats, error)
}

type transactionUsecase struct {
  transactionRepo repository.TransactionRepository
  vehicleRepo     repository.VehicleRepository
  customerRepo    repository.CustomerRepository
}

func NewTransactionUsecase(
  transactionRepo repository.TransactionRepository,
  vehicleRepo repository.VehicleRepository,
  customerRepo repository.CustomerRepository,
) TransactionUsecase {
  return &transactionUsecase{
    transactionRepo: transactionRepo,
    vehicleRepo:     vehicleRepo,
    customerRepo:    customerRepo,
  }
}

func (u *transactionUsecase) CreatePurchase(req *entity.CreatePurchaseTransactionRequest, cashierID int) (*entity.PurchaseTransaction, error) {
  // Validate vehicle exists
  vehicle, err := u.vehicleRepo.GetByID(req.VehicleID)
  if err != nil {
    return nil, fmt.Errorf("failed to get vehicle: %w", err)
  }
  if vehicle == nil {
    return nil, fmt.Errorf("vehicle not found")
  }
  
  // Validate customer exists
  customer, err := u.customerRepo.GetByID(req.CustomerID)
  if err != nil {
    return nil, fmt.Errorf("failed to get customer: %w", err)
  }
  if customer == nil {
    return nil, fmt.Errorf("customer not found")
  }
  
  // Generate transaction and invoice numbers
  transactionNumber, err := u.transactionRepo.GeneratePurchaseTransactionNumber()
  if err != nil {
    return nil, fmt.Errorf("failed to generate transaction number: %w", err)
  }
  
  invoiceNumber, err := u.transactionRepo.GeneratePurchaseInvoiceNumber()
  if err != nil {
    return nil, fmt.Errorf("failed to generate invoice number: %w", err)
  }
  
  // Calculate total amount
  totalAmount := req.VehiclePrice + req.TaxAmount
  
  // Create transaction
  transaction := &entity.PurchaseTransaction{
    TransactionNumber: transactionNumber,
    InvoiceNumber:     invoiceNumber,
    VehicleID:         req.VehicleID,
    CustomerID:        req.CustomerID,
    VehiclePrice:      req.VehiclePrice,
    TaxAmount:         req.TaxAmount,
    TotalAmount:       totalAmount,
    PaymentMethod:     req.PaymentMethod,
    PaymentReference:  req.PaymentReference,
    TransactionDate:   time.Now(),
    CashierID:         cashierID,
    Status:            "completed",
    Notes:             req.Notes,
  }
  
  if err := u.transactionRepo.CreatePurchase(transaction); err != nil {
    return nil, fmt.Errorf("failed to create purchase transaction: %w", err)
  }
  
  // Update vehicle with purchase information
  vehicle.PurchasePrice = &req.VehiclePrice
  vehicle.PurchasedFromCustomerID = &req.CustomerID
  vehicle.PurchasedByCashier = &cashierID
  now := time.Now()
  vehicle.PurchasedAt = &now
  vehicle.Status = "purchased"
  
  if err := u.vehicleRepo.Update(vehicle); err != nil {
    return nil, fmt.Errorf("failed to update vehicle: %w", err)
  }
  
  // Get the created transaction with related data
  return u.transactionRepo.GetPurchaseByID(transaction.ID)
}

func (u *transactionUsecase) GetPurchaseByID(id int) (*entity.PurchaseTransaction, error) {
  transaction, err := u.transactionRepo.GetPurchaseByID(id)
  if err != nil {
    return nil, fmt.Errorf("failed to get purchase transaction: %w", err)
  }
  
  if transaction == nil {
    return nil, fmt.Errorf("purchase transaction not found")
  }
  
  return transaction, nil
}

func (u *transactionUsecase) ListPurchases(page, limit int, search string) (*entity.TransactionListResponse, error) {
  if page <= 0 {
    page = 1
  }
  if limit <= 0 || limit > 100 {
    limit = 10
  }
  
  transactions, total, err := u.transactionRepo.ListPurchases(page, limit, search)
  if err != nil {
    return nil, fmt.Errorf("failed to list purchase transactions: %w", err)
  }
  
  return &entity.TransactionListResponse{
    Transactions: transactions,
    Total:        total,
    Page:         page,
    Limit:        limit,
  }, nil
}

func (u *transactionUsecase) CreateSales(req *entity.CreateSalesTransactionRequest, cashierID int) (*entity.SalesTransaction, error) {
  // Validate vehicle exists and is available for sale
  vehicle, err := u.vehicleRepo.GetByID(req.VehicleID)
  if err != nil {
    return nil, fmt.Errorf("failed to get vehicle: %w", err)
  }
  if vehicle == nil {
    return nil, fmt.Errorf("vehicle not found")
  }
  if vehicle.Status != "ready_to_sell" && vehicle.Status != "reserved" {
    return nil, fmt.Errorf("vehicle is not available for sale")
  }
  
  // Validate customer exists
  customer, err := u.customerRepo.GetByID(req.CustomerID)
  if err != nil {
    return nil, fmt.Errorf("failed to get customer: %w", err)
  }
  if customer == nil {
    return nil, fmt.Errorf("customer not found")
  }
  
  // Generate transaction and invoice numbers
  transactionNumber, err := u.transactionRepo.GenerateSalesTransactionNumber()
  if err != nil {
    return nil, fmt.Errorf("failed to generate transaction number: %w", err)
  }
  
  invoiceNumber, err := u.transactionRepo.GenerateSalesInvoiceNumber()
  if err != nil {
    return nil, fmt.Errorf("failed to generate invoice number: %w", err)
  }
  
  // Calculate total amount
  totalAmount := req.VehiclePrice + req.TaxAmount - req.DiscountAmount
  
  // Create transaction
  transaction := &entity.SalesTransaction{
    TransactionNumber: transactionNumber,
    InvoiceNumber:     invoiceNumber,
    VehicleID:         req.VehicleID,
    CustomerID:        req.CustomerID,
    VehiclePrice:      req.VehiclePrice,
    TaxAmount:         req.TaxAmount,
    DiscountAmount:    req.DiscountAmount,
    TotalAmount:       totalAmount,
    PaymentMethod:     req.PaymentMethod,
    PaymentReference:  req.PaymentReference,
    TransactionDate:   time.Now(),
    CashierID:         cashierID,
    Status:            "completed",
    Notes:             req.Notes,
  }
  
  if err := u.transactionRepo.CreateSales(transaction); err != nil {
    return nil, fmt.Errorf("failed to create sales transaction: %w", err)
  }
  
  // Update vehicle with sales information
  vehicle.FinalSellingPrice = &req.VehiclePrice
  vehicle.SoldToCustomerID = &req.CustomerID
  vehicle.SoldByCashier = &cashierID
  now := time.Now()
  vehicle.SoldAt = &now
  vehicle.Status = "sold"
  
  if err := u.vehicleRepo.Update(vehicle); err != nil {
    return nil, fmt.Errorf("failed to update vehicle: %w", err)
  }
  
  // Get the created transaction with related data
  return u.transactionRepo.GetSalesByID(transaction.ID)
}

func (u *transactionUsecase) GetSalesByID(id int) (*entity.SalesTransaction, error) {
  transaction, err := u.transactionRepo.GetSalesByID(id)
  if err != nil {
    return nil, fmt.Errorf("failed to get sales transaction: %w", err)
  }
  
  if transaction == nil {
    return nil, fmt.Errorf("sales transaction not found")
  }
  
  return transaction, nil
}

func (u *transactionUsecase) ListSales(page, limit int, search string) (*entity.TransactionListResponse, error) {
  if page <= 0 {
    page = 1
  }
  if limit <= 0 || limit > 100 {
    limit = 10
  }
  
  transactions, total, err := u.transactionRepo.ListSales(page, limit, search)
  if err != nil {
    return nil, fmt.Errorf("failed to list sales transactions: %w", err)
  }
  
  return &entity.TransactionListResponse{
    Transactions: transactions,
    Total:        total,
    Page:         page,
    Limit:        limit,
  }, nil
}

func (u *transactionUsecase) GetDashboardStats() (*entity.DashboardStats, error) {
  stats, err := u.transactionRepo.GetDashboardStats()
  if err != nil {
    return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
  }
  
  return stats, nil
}
