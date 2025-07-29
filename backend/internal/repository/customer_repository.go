package repository

import (
  "database/sql"
  "fmt"
  "strings"

  "github.com/jmoiron/sqlx"
  "vehicle-showroom/internal/entity"
)

type CustomerRepository interface {
  Create(customer *entity.Customer) error
  GetByID(id int) (*entity.Customer, error)
  GetByCode(code string) (*entity.Customer, error)
  List(page, limit int, search string) ([]entity.Customer, int, error)
  Update(customer *entity.Customer) error
  Delete(id int) error
  GenerateCustomerCode() (string, error)
}

type customerRepository struct {
  db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) CustomerRepository {
  return &customerRepository{db: db}
}

func (r *customerRepository) Create(customer *entity.Customer) error {
  query := `
    INSERT INTO customers (customer_code, name, phone, email, address, id_card_number, type, created_by, is_active)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING id, created_at, updated_at
  `
  
  err := r.db.QueryRow(
    query,
    customer.CustomerCode,
    customer.Name,
    customer.Phone,
    customer.Email,
    customer.Address,
    customer.IDCardNumber,
    customer.Type,
    customer.CreatedBy,
    customer.IsActive,
  ).Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)
  
  if err != nil {
    return fmt.Errorf("failed to create customer: %w", err)
  }
  
  return nil
}

func (r *customerRepository) GetByID(id int) (*entity.Customer, error) {
  customer := &entity.Customer{}
  query := `
    SELECT id, customer_code, name, phone, email, address, id_card_number, type, 
           created_at, updated_at, created_by, is_active
    FROM customers
    WHERE id = $1 AND is_active = true
  `
  
  err := r.db.Get(customer, query, id)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("failed to get customer by id: %w", err)
  }
  
  return customer, nil
}

func (r *customerRepository) GetByCode(code string) (*entity.Customer, error) {
  customer := &entity.Customer{}
  query := `
    SELECT id, customer_code, name, phone, email, address, id_card_number, type, 
           created_at, updated_at, created_by, is_active
    FROM customers
    WHERE customer_code = $1 AND is_active = true
  `
  
  err := r.db.Get(customer, query, code)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("failed to get customer by code: %w", err)
  }
  
  return customer, nil
}

func (r *customerRepository) List(page, limit int, search string) ([]entity.Customer, int, error) {
  offset := (page - 1) * limit
  
  whereClause := "WHERE is_active = true"
  args := []interface{}{}
  argIndex := 1
  
  if search != "" {
    whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR customer_code ILIKE $%d OR phone ILIKE $%d OR email ILIKE $%d)", 
                              argIndex, argIndex+1, argIndex+2, argIndex+3)
    searchPattern := "%" + search + "%"
    args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
    argIndex += 4
  }
  
  // Get total count
  countQuery := fmt.Sprintf("SELECT COUNT(*) FROM customers %s", whereClause)
  var total int
  err := r.db.Get(&total, countQuery, args...)
  if err != nil {
    return nil, 0, fmt.Errorf("failed to get customer count: %w", err)
  }
  
  // Get customers
  query := fmt.Sprintf(`
    SELECT id, customer_code, name, phone, email, address, id_card_number, type, 
           created_at, updated_at, created_by, is_active
    FROM customers
    %s
    ORDER BY created_at DESC
    LIMIT $%d OFFSET $%d
  `, whereClause, argIndex, argIndex+1)
  
  args = append(args, limit, offset)
  
  var customers []entity.Customer
  err = r.db.Select(&customers, query, args...)
  if err != nil {
    return nil, 0, fmt.Errorf("failed to list customers: %w", err)
  }
  
  return customers, total, nil
}

func (r *customerRepository) Update(customer *entity.Customer) error {
  query := `
    UPDATE customers
    SET name = $1, phone = $2, email = $3, address = $4, id_card_number = $5, 
        type = $6, updated_at = CURRENT_TIMESTAMP
    WHERE id = $7
  `
  
  _, err := r.db.Exec(query, customer.Name, customer.Phone, customer.Email, 
                     customer.Address, customer.IDCardNumber, customer.Type, customer.ID)
  if err != nil {
    return fmt.Errorf("failed to update customer: %w", err)
  }
  
  return nil
}

func (r *customerRepository) Delete(id int) error {
  query := `UPDATE customers SET is_active = false WHERE id = $1`
  
  _, err := r.db.Exec(query, id)
  if err != nil {
    return fmt.Errorf("failed to delete customer: %w", err)
  }
  
  return nil
}

func (r *customerRepository) GenerateCustomerCode() (string, error) {
  var lastCode string
  query := `
    SELECT customer_code 
    FROM customers 
    WHERE customer_code LIKE 'CUST-%' 
    ORDER BY customer_code DESC 
    LIMIT 1
  `
  
  err := r.db.Get(&lastCode, query)
  if err != nil && err != sql.ErrNoRows {
    return "", fmt.Errorf("failed to get last customer code: %w", err)
  }
  
  nextNumber := 1
  if lastCode != "" {
    // Extract number from CUST-XXX format
    parts := strings.Split(lastCode, "-")
    if len(parts) == 2 {
      var num int
      fmt.Sscanf(parts[1], "%d", &num)
      nextNumber = num + 1
    }
  }
  
  return fmt.Sprintf("CUST-%03d", nextNumber), nil
}
