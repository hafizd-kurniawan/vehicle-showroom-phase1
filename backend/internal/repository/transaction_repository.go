package repository

import (
  "database/sql"
  "fmt"
  "strings"
  "time"

  "github.com/jmoiron/sqlx"
  "vehicle-showroom/internal/entity"
)

type TransactionRepository interface {
  // Purchase Transactions
  CreatePurchase(tx *entity.PurchaseTransaction) error
  GetPurchaseByID(id int) (*entity.PurchaseTransaction, error)
  ListPurchases(page, limit int, search string) ([]entity.PurchaseTransaction, int, error)
  
  // Sales Transactions
  CreateSales(tx *entity.SalesTransaction) error
  GetSalesByID(id int) (*entity.SalesTransaction, error)
  ListSales(page, limit int, search string) ([]entity.SalesTransaction, int, error)
  
  // Transaction Numbers
  GeneratePurchaseTransactionNumber() (string, error)
  GenerateSalesTransactionNumber() (string, error)
  GeneratePurchaseInvoiceNumber() (string, error)
  GenerateSalesInvoiceNumber() (string, error)
  
  // Dashboard Stats
  GetDashboardStats() (*entity.DashboardStats, error)
}

type transactionRepository struct {
  db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
  return &transactionRepository{db: db}
}

func (r *transactionRepository) CreatePurchase(tx *entity.PurchaseTransaction) error {
  query := `
    INSERT INTO purchase_transactions (
      transaction_number, invoice_number, vehicle_id, customer_id, vehicle_price,
      tax_amount, total_amount, payment_method, payment_reference, transaction_date,
      cashier_id, status, notes
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    RETURNING id, created_at
  `
  
  err := r.db.QueryRow(
    query,
    tx.TransactionNumber,
    tx.InvoiceNumber,
    tx.VehicleID,
    tx.CustomerID,
    tx.VehiclePrice,
    tx.TaxAmount,
    tx.TotalAmount,
    tx.PaymentMethod,
    tx.PaymentReference,
    tx.TransactionDate,
    tx.CashierID,
    tx.Status,
    tx.Notes,
  ).Scan(&tx.ID, &tx.CreatedAt)
  
  if err != nil {
    return fmt.Errorf("failed to create purchase transaction: %w", err)
  }
  
  return nil
}

func (r *transactionRepository) GetPurchaseByID(id int) (*entity.PurchaseTransaction, error) {
  tx := &entity.PurchaseTransaction{}
  query := `
    SELECT pt.id, pt.transaction_number, pt.invoice_number, pt.vehicle_id, pt.customer_id,
           pt.vehicle_price, pt.tax_amount, pt.total_amount, pt.payment_method,
           pt.payment_reference, pt.transaction_date, pt.cashier_id, pt.status,
           pt.notes, pt.created_at
    FROM purchase_transactions pt
    WHERE pt.id = $1
  `
  
  err := r.db.Get(tx, query, id)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("failed to get purchase transaction by id: %w", err)
  }
  
  // Load related data
  r.loadPurchaseRelatedData(tx)
  
  return tx, nil
}

func (r *transactionRepository) ListPurchases(page, limit int, search string) ([]entity.PurchaseTransaction, int, error) {
  offset := (page - 1) * limit
  
  whereClause := "WHERE 1=1"
  args := []interface{}{}
  argIndex := 1
  
  if search != "" {
    whereClause += fmt.Sprintf(" AND (pt.transaction_number ILIKE $%d OR pt.invoice_number ILIKE $%d OR c.name ILIKE $%d OR v.brand ILIKE $%d OR v.model ILIKE $%d)", 
                              argIndex, argIndex+1, argIndex+2, argIndex+3, argIndex+4)
    searchPattern := "%" + search + "%"
    args = append(args, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
    argIndex += 5
  }
  
  // Get total count
  countQuery := fmt.Sprintf(`
    SELECT COUNT(*) 
    FROM purchase_transactions pt
    LEFT JOIN customers c ON pt.customer_id = c.id
    LEFT JOIN vehicles v ON pt.vehicle_id = v.id
    %s
  `, whereClause)
  
  var total int
  err := r.db.Get(&total, countQuery, args...)
  if err != nil {
    return nil, 0, fmt.Errorf("failed to get purchase transaction count: %w", err)
  }
  
  // Get transactions
  query := fmt.Sprintf(`
    SELECT pt.id, pt.transaction_number, pt.invoice_number, pt.vehicle_id, pt.customer_id,
           pt.vehicle_price, pt.tax_amount, pt.total_amount, pt.payment_method,
           pt.payment_reference, pt.transaction_date, pt.cashier_id, pt.status,
           pt.notes, pt.created_at
    FROM purchase_transactions pt
    LEFT JOIN customers c ON pt.customer_id = c.id
    LEFT JOIN vehicles v ON pt.vehicle_id = v.id
    %s
    ORDER BY pt.created_at DESC
    LIMIT $%d OFFSET $%d
  `, whereClause, argIndex, argIndex+1)
  
  args = append(args, limit, offset)
  
  var transactions []entity.PurchaseTransaction
  err = r.db.Select(&transactions, query, args...)
  if err != nil {
    return nil, 0, fmt.Errorf("failed to list purchase transactions: %w", err)
  }
  
  // Load related data for each transaction
  for i := range transactions {
    r.loadPurchaseRelatedData(&transactions[i])
  }
  
  return transactions, total, nil
}

func (r *transactionRepository) CreateSales(tx *entity.SalesTransaction) error {
  query := `
    INSERT INTO sales_transactions (
      transaction_number, invoice_number, vehicle_id, customer_id, vehicle_price,
      tax_amount, discount_amount, total_amount, payment_method, payment_reference,
      transaction_date, cashier_id, status, notes
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    RETURNING id, created_at
  `
  
  err := r.db.QueryRow(
    query,
    tx.TransactionNumber,
    tx.InvoiceNumber,
    tx.VehicleID,
    tx.CustomerID,
    tx.VehiclePrice,
    tx.TaxAmount,
    tx.DiscountAmount,
    tx.TotalAmount,
    tx.PaymentMethod,
    tx.PaymentReference,
    tx.TransactionDate,
    tx.CashierID,
    tx.Status,
    tx.Notes,
  ).Scan(&tx.ID, &tx.CreatedAt)
  
  if err != nil {
    return fmt.Errorf("failed to create sales transaction: %w", err)
  }
  
  return nil
}

func (r *transactionRepository) GetSalesByID(id int) (*entity.SalesTransaction, error) {
  tx := &entity.SalesTransaction{}
  query := `
    SELECT st.id, st.transaction_number, st.invoice_number, st.vehicle_id, st.customer_id,
           st.vehicle_price, st.tax_amount, st.discount_amount, st.total_amount,
           st.payment_method, st.payment_reference, st.transaction_date, st.cashier_id,
           st.status, st.notes, st.created_at
    FROM sales_transactions st
    WHERE st.id = $1
  `
  
  err := r.db.Get(tx, query, id)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("failed to get sales transaction by id: %w", err)
  }
  
  // Load related data
  r.loadSalesRelatedData(tx)
  
  return tx, nil
}

func (r *transactionRepository) ListSales(page, limit int, search string) ([]entity.SalesTransaction, int, error) {
  offset := (page - 1) * limit
  
  whereClause := "WHERE 1=1"
  args := []interface{}{}
  argIndex := 1
  
  if search != "" {
    whereClause += fmt.Sprintf(" AND (st.transaction_number ILIKE $%d OR st.invoice_number ILIKE $%d OR c.name ILIKE $%d OR v.brand ILIKE $%d OR v.model ILIKE $%d)", 
                              argIndex, argIndex+1, argIndex+2, argIndex+3, argIndex+4)
    searchPattern := "%" + search + "%"
    args = append(args, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
    argIndex += 5
  }
  
  // Get total count
  countQuery := fmt.Sprintf(`
    SELECT COUNT(*) 
    FROM sales_transactions st
    LEFT JOIN customers c ON st.customer_id = c.id
    LEFT JOIN vehicles v ON st.vehicle_id = v.id
    %s
  `, whereClause)
  
  var total int
  err := r.db.Get(&total, countQuery, args...)
  if err != nil {
    return nil, 0, fmt.Errorf("failed to get sales transaction count: %w", err)
  }
  
  // Get transactions
  query := fmt.Sprintf(`
    SELECT st.id, st.transaction_number, st.invoice_number, st.vehicle_id, st.customer_id,
           st.vehicle_price, st.tax_amount, st.discount_amount, st.total_amount,
           st.payment_method, st.payment_reference, st.transaction_date, st.cashier_id,
           st.status, st.notes, st.created_at
    FROM sales_transactions st
    LEFT JOIN customers c ON st.customer_id = c.id
    LEFT JOIN vehicles v ON st.vehicle_id = v.id
    %s
    ORDER BY st.created_at DESC
    LIMIT $%d OFFSET $%d
  `, whereClause, argIndex, argIndex+1)
  
  args = append(args, limit, offset)
  
  var transactions []entity.SalesTransaction
  err = r.db.Select(&transactions, query, args...)
  if err != nil {
    return nil, 0, fmt.Errorf("failed to list sales transactions: %w", err)
  }
  
  // Load related data for each transaction
  for i := range transactions {
    r.loadSalesRelatedData(&transactions[i])
  }
  
  return transactions, total, nil
}

func (r *transactionRepository) GeneratePurchaseTransactionNumber() (string, error) {
  today := time.Now().Format("20060102")
  var lastNumber string
  query := `
    SELECT transaction_number 
    FROM purchase_transactions 
    WHERE transaction_number LIKE 'PUR-' || $1 || '-%' 
    ORDER BY transaction_number DESC 
    LIMIT 1
  `
  
  err := r.db.Get(&lastNumber, query, today)
  if err != nil && err != sql.ErrNoRows {
    return "", fmt.Errorf("failed to get last purchase transaction number: %w", err)
  }
  
  nextNumber := 1
  if lastNumber != "" {
    parts := strings.Split(lastNumber, "-")
    if len(parts) == 3 {
      var num int
      fmt.Sscanf(parts[2], "%d", &num)
      nextNumber = num + 1
    }
  }
  
  return fmt.Sprintf("PUR-%s-%03d", today, nextNumber), nil
}

func (r *transactionRepository) GenerateSalesTransactionNumber() (string, error) {
  today := time.Now().Format("20060102")
  var lastNumber string
  query := `
    SELECT transaction_number 
    FROM sales_transactions 
    WHERE transaction_number LIKE 'SAL-' || $1 || '-%' 
    ORDER BY transaction_number DESC 
    LIMIT 1
  `
  
  err := r.db.Get(&lastNumber, query, today)
  if err != nil && err != sql.ErrNoRows {
    return "", fmt.Errorf("failed to get last sales transaction number: %w", err)
  }
  
  nextNumber := 1
  if lastNumber != "" {
    parts := strings.Split(lastNumber, "-")
    if len(parts) == 3 {
      var num int
      fmt.Sscanf(parts[2], "%d", &num)
      nextNumber = num + 1
    }
  }
  
  return fmt.Sprintf("SAL-%s-%03d", today, nextNumber), nil
}

func (r *transactionRepository) GeneratePurchaseInvoiceNumber() (string, error) {
  today := time.Now().Format("20060102")
  var lastNumber string
  query := `
    SELECT invoice_number 
    FROM purchase_transactions 
    WHERE invoice_number LIKE 'INV-PUR-' || $1 || '-%' 
    ORDER BY invoice_number DESC 
    LIMIT 1
  `
  
  err := r.db.Get(&lastNumber, query, today)
  if err != nil && err != sql.ErrNoRows {
    return "", fmt.Errorf("failed to get last purchase invoice number: %w", err)
  }
  
  nextNumber := 1
  if lastNumber != "" {
    parts := strings.Split(lastNumber, "-")
    if len(parts) == 4 {
      var num int
      fmt.Sscanf(parts[3], "%d", &num)
      nextNumber = num + 1
    }
  }
  
  return fmt.Sprintf("INV-PUR-%s-%03d", today, nextNumber), nil
}

func (r *transactionRepository) GenerateSalesInvoiceNumber() (string, error) {
  today := time.Now().Format("20060102")
  var lastNumber string
  query := `
    SELECT invoice_number 
    FROM sales_transactions 
    WHERE invoice_number LIKE 'INV-SAL-' || $1 || '-%' 
    ORDER BY invoice_number DESC 
    LIMIT 1
  `
  
  err := r.db.Get(&lastNumber, query, today)
  if err != nil && err != sql.ErrNoRows {
    return "", fmt.Errorf("failed to get last sales invoice number: %w", err)
  }
  
  nextNumber := 1
  if lastNumber != "" {
    parts := strings.Split(lastNumber, "-")
    if len(parts) == 4 {
      var num int
      fmt.Sscanf(parts[3], "%d", &num)
      nextNumber = num + 1
    }
  }
  
  return fmt.Sprintf("INV-SAL-%s-%03d", today, nextNumber), nil
}

func (r *transactionRepository) GetDashboardStats() (*entity.DashboardStats, error) {
  stats := &entity.DashboardStats{}
  
  // Get vehicle counts
  err := r.db.Get(&stats.TotalVehicles, "SELECT COUNT(*) FROM vehicles")
  if err != nil {
    return nil, fmt.Errorf("failed to get total vehicles: %w", err)
  }
  
  err = r.db.Get(&stats.VehiclesForSale, "SELECT COUNT(*) FROM vehicles WHERE status = 'ready_to_sell'")
  if err != nil {
    return nil, fmt.Errorf("failed to get vehicles for sale: %w", err)
  }
  
  err = r.db.Get(&stats.VehiclesInRepair, "SELECT COUNT(*) FROM vehicles WHERE status = 'in_repair'")
  if err != nil {
    return nil, fmt.Errorf("failed to get vehicles in repair: %w", err)
  }
  
  err = r.db.Get(&stats.VehiclesSold, "SELECT COUNT(*) FROM vehicles WHERE status = 'sold'")
  if err != nil {
    return nil, fmt.Errorf("failed to get vehicles sold: %w", err)
  }
  
  // Get customer count
  err = r.db.Get(&stats.TotalCustomers, "SELECT COUNT(*) FROM customers WHERE is_active = true")
  if err != nil {
    return nil, fmt.Errorf("failed to get total customers: %w", err)
  }
  
  // Get today's transactions
  today := time.Now().Format("2006-01-02")
  err = r.db.Get(&stats.TodayPurchases, "SELECT COUNT(*) FROM purchase_transactions WHERE DATE(transaction_date) = $1", today)
  if err != nil {
    return nil, fmt.Errorf("failed to get today purchases: %w", err)
  }
  
  err = r.db.Get(&stats.TodaySales, "SELECT COUNT(*) FROM sales_transactions WHERE DATE(transaction_date) = $1", today)
  if err != nil {
    return nil, fmt.Errorf("failed to get today sales: %w", err)
  }
  
  // Get today's revenue
  var todayRevenue sql.NullFloat64
  err = r.db.Get(&todayRevenue, "SELECT COALESCE(SUM(total_amount), 0) FROM sales_transactions WHERE DATE(transaction_date) = $1", today)
  if err != nil {
    return nil, fmt.Errorf("failed to get today revenue: %w", err)
  }
  stats.TodayRevenue = todayRevenue.Float64
  
  // Get monthly revenue
  monthStart := time.Now().Format("2006-01-01")
  var monthlyRevenue sql.NullFloat64
  err = r.db.Get(&monthlyRevenue, "SELECT COALESCE(SUM(total_amount), 0) FROM sales_transactions WHERE transaction_date >= $1", monthStart)
  if err != nil {
    return nil, fmt.Errorf("failed to get monthly revenue: %w", err)
  }
  stats.MonthlyRevenue = monthlyRevenue.Float64
  
  // Calculate total profit (simplified: sales - purchases)
  var totalSales, totalPurchases sql.NullFloat64
  err = r.db.Get(&totalSales, "SELECT COALESCE(SUM(total_amount), 0) FROM sales_transactions")
  if err != nil {
    return nil, fmt.Errorf("failed to get total sales: %w", err)
  }
  
  err = r.db.Get(&totalPurchases, "SELECT COALESCE(SUM(total_amount), 0) FROM purchase_transactions")
  if err != nil {
    return nil, fmt.Errorf("failed to get total purchases: %w", err)
  }
  
  stats.TotalProfit = totalSales.Float64 - totalPurchases.Float64
  
  return stats, nil
}

func (r *transactionRepository) loadPurchaseRelatedData(tx *entity.PurchaseTransaction) {
  // Load vehicle
  vehicle := &entity.Vehicle{}
  vehicleQuery := `
    SELECT id, vehicle_code, chassis_number, license_plate, brand, model, variant, year, color
    FROM vehicles WHERE id = $1
  `
  err := r.db.Get(vehicle, vehicleQuery, tx.VehicleID)
  if err == nil {
    tx.Vehicle = vehicle
  }
  
  // Load customer
  customer := &entity.Customer{}
  customerQuery := `
    SELECT id, customer_code, name, phone, email, type
    FROM customers WHERE id = $1
  `
  err = r.db.Get(customer, customerQuery, tx.CustomerID)
  if err == nil {
    tx.Customer = customer
  }
  
  // Load cashier
  cashier := &entity.User{}
  cashierQuery := `
    SELECT id, username, full_name, role
    FROM users WHERE id = $1
  `
  err = r.db.Get(cashier, cashierQuery, tx.CashierID)
  if err == nil {
    tx.Cashier = cashier
  }
}

func (r *transactionRepository) loadSalesRelatedData(tx *entity.SalesTransaction) {
  // Load vehicle
  vehicle := &entity.Vehicle{}
  vehicleQuery := `
    SELECT id, vehicle_code, chassis_number, license_plate, brand, model, variant, year, color
    FROM vehicles WHERE id = $1
  `
  err := r.db.Get(vehicle, vehicleQuery, tx.VehicleID)
  if err == nil {
    tx.Vehicle = vehicle
  }
  
  // Load customer
  customer := &entity.Customer{}
  customerQuery := `
    SELECT id, customer_code, name, phone, email, type
    FROM customers WHERE id = $1
  `
  err = r.db.Get(customer, customerQuery, tx.CustomerID)
  if err == nil {
    tx.Customer = customer
  }
  
  // Load cashier
  cashier := &entity.User{}
  cashierQuery := `
    SELECT id, username, full_name, role
    FROM users WHERE id = $1
  `
  err = r.db.Get(cashier, cashierQuery, tx.CashierID)
  if err == nil {
    tx.Cashier = cashier
  }
}
