package repository

import (
  "database/sql"
  "fmt"
  "strings"

  "github.com/jmoiron/sqlx"
  "vehicle-showroom/internal/entity"
)

type VehicleRepository interface {
  Create(vehicle *entity.Vehicle) error
  GetByID(id int) (*entity.Vehicle, error)
  GetByCode(code string) (*entity.Vehicle, error)
  List(page, limit int, search, status string) ([]entity.Vehicle, int, error)
  Update(vehicle *entity.Vehicle) error
  UpdateStatus(id int, status string) error
  Delete(id int) error
  GenerateVehicleCode() (string, error)
}

type vehicleRepository struct {
  db *sqlx.DB
}

func NewVehicleRepository(db *sqlx.DB) VehicleRepository {
  return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(vehicle *entity.Vehicle) error {
  query := `
    INSERT INTO vehicles (
      vehicle_code, chassis_number, license_plate, brand, model, variant, year, color, mileage,
      fuel_type, transmission, purchase_price, status, purchased_from_customer_id,
      purchased_by_cashier, purchased_at, purchase_notes, condition_notes
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
    RETURNING id, created_at, updated_at
  `
  
  err := r.db.QueryRow(
    query,
    vehicle.VehicleCode,
    vehicle.ChassisNumber,
    vehicle.LicensePlate,
    vehicle.Brand,
    vehicle.Model,
    vehicle.Variant,
    vehicle.Year,
    vehicle.Color,
    vehicle.Mileage,
    vehicle.FuelType,
    vehicle.Transmission,
    vehicle.PurchasePrice,
    vehicle.Status,
    vehicle.PurchasedFromCustomerID,
    vehicle.PurchasedByCashier,
    vehicle.PurchasedAt,
    vehicle.PurchaseNotes,
    vehicle.ConditionNotes,
  ).Scan(&vehicle.ID, &vehicle.CreatedAt, &vehicle.UpdatedAt)
  
  if err != nil {
    return fmt.Errorf("failed to create vehicle: %w", err)
  }
  
  return nil
}

func (r *vehicleRepository) GetByID(id int) (*entity.Vehicle, error) {
  vehicle := &entity.Vehicle{}
  query := `
    SELECT v.id, v.vehicle_code, v.chassis_number, v.license_plate, v.brand, v.model, v.variant,
           v.year, v.color, v.mileage, v.fuel_type, v.transmission, v.purchase_price,
           v.total_repair_cost, v.suggested_selling_price, v.approved_selling_price,
           v.final_selling_price, v.status, v.purchased_from_customer_id, v.sold_to_customer_id,
           v.purchased_by_cashier, v.sold_by_cashier, v.price_approved_by_admin,
           v.purchased_at, v.sold_at, v.created_at, v.updated_at, v.purchase_notes, v.condition_notes
    FROM vehicles v
    WHERE v.id = $1
  `
  
  err := r.db.Get(vehicle, query, id)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("failed to get vehicle by id: %w", err)
  }
  
  // Load related customer data
  if vehicle.PurchasedFromCustomerID != nil {
    customer := &entity.Customer{}
    customerQuery := `
      SELECT id, customer_code, name, phone, email, address, id_card_number, type,
             created_at, updated_at, created_by, is_active
      FROM customers
      WHERE id = $1
    `
    err = r.db.Get(customer, customerQuery, *vehicle.PurchasedFromCustomerID)
    if err == nil {
      vehicle.PurchasedFromCustomer = customer
    }
  }
  
  if vehicle.SoldToCustomerID != nil {
    customer := &entity.Customer{}
    customerQuery := `
      SELECT id, customer_code, name, phone, email, address, id_card_number, type,
             created_at, updated_at, created_by, is_active
      FROM customers
      WHERE id = $1
    `
    err = r.db.Get(customer, customerQuery, *vehicle.SoldToCustomerID)
    if err == nil {
      vehicle.SoldToCustomer = customer
    }
  }
  
  return vehicle, nil
}

func (r *vehicleRepository) GetByCode(code string) (*entity.Vehicle, error) {
  vehicle := &entity.Vehicle{}
  query := `
    SELECT id, vehicle_code, chassis_number, license_plate, brand, model, variant,
           year, color, mileage, fuel_type, transmission, purchase_price,
           total_repair_cost, suggested_selling_price, approved_selling_price,
           final_selling_price, status, purchased_from_customer_id, sold_to_customer_id,
           purchased_by_cashier, sold_by_cashier, price_approved_by_admin,
           purchased_at, sold_at, created_at, updated_at, purchase_notes, condition_notes
    FROM vehicles
    WHERE vehicle_code = $1
  `
  
  err := r.db.Get(vehicle, query, code)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("failed to get vehicle by code: %w", err)
  }
  
  return vehicle, nil
}

func (r *vehicleRepository) List(page, limit int, search, status string) ([]entity.Vehicle, int, error) {
  offset := (page - 1) * limit
  
  whereClause := "WHERE 1=1"
  args := []interface{}{}
  argIndex := 1
  
  if search != "" {
    whereClause += fmt.Sprintf(" AND (v.brand ILIKE $%d OR v.model ILIKE $%d OR v.vehicle_code ILIKE $%d OR v.chassis_number ILIKE $%d OR v.license_plate ILIKE $%d)", 
                              argIndex, argIndex+1, argIndex+2, argIndex+3, argIndex+4)
    searchPattern := "%" + search + "%"
    args = append(args, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
    argIndex += 5
  }
  
  if status != "" {
    whereClause += fmt.Sprintf(" AND v.status = $%d", argIndex)
    args = append(args, status)
    argIndex++
  }
  
  // Get total count
  countQuery := fmt.Sprintf("SELECT COUNT(*) FROM vehicles v %s", whereClause)
  var total int
  err := r.db.Get(&total, countQuery, args...)
  if err != nil {
    return nil, 0, fmt.Errorf("failed to get vehicle count: %w", err)
  }
  
  // Get vehicles with customer data
  query := fmt.Sprintf(`
    SELECT v.id, v.vehicle_code, v.chassis_number, v.license_plate, v.brand, v.model, v.variant,
           v.year, v.color, v.mileage, v.fuel_type, v.transmission, v.purchase_price,
           v.total_repair_cost, v.suggested_selling_price, v.approved_selling_price,
           v.final_selling_price, v.status, v.purchased_from_customer_id, v.sold_to_customer_id,
           v.purchased_by_cashier, v.sold_by_cashier, v.price_approved_by_admin,
           v.purchased_at, v.sold_at, v.created_at, v.updated_at, v.purchase_notes, v.condition_notes
    FROM vehicles v
    %s
    ORDER BY v.created_at DESC
    LIMIT $%d OFFSET $%d
  `, whereClause, argIndex, argIndex+1)
  
  args = append(args, limit, offset)
  
  var vehicles []entity.Vehicle
  err = r.db.Select(&vehicles, query, args...)
  if err != nil {
    return nil, 0, fmt.Errorf("failed to list vehicles: %w", err)
  }
  
  // Load customer data for each vehicle
  for i := range vehicles {
    if vehicles[i].PurchasedFromCustomerID != nil {
      customer := &entity.Customer{}
      customerQuery := `
        SELECT id, customer_code, name, phone, email, address, id_card_number, type,
               created_at, updated_at, created_by, is_active
        FROM customers
        WHERE id = $1
      `
      err = r.db.Get(customer, customerQuery, *vehicles[i].PurchasedFromCustomerID)
      if err == nil {
        vehicles[i].PurchasedFromCustomer = customer
      }
    }
  }
  
  return vehicles, total, nil
}

func (r *vehicleRepository) Update(vehicle *entity.Vehicle) error {
  query := `
    UPDATE vehicles
    SET license_plate = $1, brand = $2, model = $3, variant = $4, year = $5, color = $6,
        mileage = $7, fuel_type = $8, transmission = $9, suggested_selling_price = $10,
        purchase_notes = $11, condition_notes = $12, updated_at = CURRENT_TIMESTAMP
    WHERE id = $13
  `
  
  _, err := r.db.Exec(query, vehicle.LicensePlate, vehicle.Brand, vehicle.Model, vehicle.Variant,
                     vehicle.Year, vehicle.Color, vehicle.Mileage, vehicle.FuelType,
                     vehicle.Transmission, vehicle.SuggestedSellingPrice, vehicle.PurchaseNotes,
                     vehicle.ConditionNotes, vehicle.ID)
  if err != nil {
    return fmt.Errorf("failed to update vehicle: %w", err)
  }
  
  return nil
}

func (r *vehicleRepository) UpdateStatus(id int, status string) error {
  query := `UPDATE vehicles SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
  
  _, err := r.db.Exec(query, status, id)
  if err != nil {
    return fmt.Errorf("failed to update vehicle status: %w", err)
  }
  
  return nil
}

func (r *vehicleRepository) Delete(id int) error {
  query := `DELETE FROM vehicles WHERE id = $1`
  
  _, err := r.db.Exec(query, id)
  if err != nil {
    return fmt.Errorf("failed to delete vehicle: %w", err)
  }
  
  return nil
}

func (r *vehicleRepository) GenerateVehicleCode() (string, error) {
  var lastCode string
  query := `
    SELECT vehicle_code 
    FROM vehicles 
    WHERE vehicle_code LIKE 'VEH-%' 
    ORDER BY vehicle_code DESC 
    LIMIT 1
  `
  
  err := r.db.Get(&lastCode, query)
  if err != nil && err != sql.ErrNoRows {
    return "", fmt.Errorf("failed to get last vehicle code: %w", err)
  }
  
  nextNumber := 1
  if lastCode != "" {
    // Extract number from VEH-XXX format
    parts := strings.Split(lastCode, "-")
    if len(parts) == 2 {
      var num int
      fmt.Sscanf(parts[1], "%d", &num)
      nextNumber = num + 1
    }
  }
  
  return fmt.Sprintf("VEH-%03d", nextNumber), nil
}
