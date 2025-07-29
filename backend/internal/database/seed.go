package database

import (
  "log"

  "github.com/jmoiron/sqlx"
  "golang.org/x/crypto/bcrypt"
)

func SeedDemoData(db *sqlx.DB) error {
  // Check if demo users already exist
  var count int
  err := db.Get(&count, "SELECT COUNT(*) FROM users")
  if err != nil {
    return err
  }
  
  if count > 0 {
    log.Println("Demo data already exists, skipping seed")
    return nil
  }

  log.Println("Seeding demo data...")

  // Create demo users
  users := []struct {
    username, email, password, fullName, phone, role string
  }{
    {"admin", "admin@showroom.com", "admin123", "Admin User", "081234567890", "admin"},
    {"cashier", "cashier@showroom.com", "cashier123", "Cashier User", "081234567891", "cashier"},
    {"mechanic", "mechanic@showroom.com", "mechanic123", "Mechanic User", "081234567892", "mechanic"},
  }

  for _, user := range users {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
    if err != nil {
      return err
    }

    _, err = db.Exec(`
      INSERT INTO users (username, email, password_hash, full_name, phone, role, is_active)
      VALUES ($1, $2, $3, $4, $5, $6, true)
    `, user.username, user.email, string(hashedPassword), user.fullName, user.phone, user.role)
    
    if err != nil {
      return err
    }
  }

  // Create demo customers
  customers := []struct {
    code, name, phone, email, address, idCard, customerType string
  }{
    {"CUST-001", "John Doe", "081234567893", "john@email.com", "Jl. Sudirman No. 123, Jakarta", "3171234567890001", "individual"},
    {"CUST-002", "Jane Smith", "081234567894", "jane@email.com", "Jl. Thamrin No. 456, Jakarta", "3171234567890002", "individual"},
    {"CUST-003", "PT. Maju Jaya", "081234567895", "info@majujaya.com", "Jl. Gatot Subroto No. 789, Jakarta", "021234567890", "corporate"},
    {"CUST-004", "Ahmad Rahman", "081234567896", "ahmad@email.com", "Jl. Kuningan No. 321, Jakarta", "3171234567890003", "individual"},
    {"CUST-005", "CV. Berkah Motor", "081234567897", "info@berkahmotor.com", "Jl. Casablanca No. 654, Jakarta", "021234567891", "corporate"},
  }

  for _, customer := range customers {
    _, err = db.Exec(`
      INSERT INTO customers (customer_code, name, phone, email, address, id_card_number, type, created_by, is_active)
      VALUES ($1, $2, $3, $4, $5, $6, $7, 1, true)
    `, customer.code, customer.name, customer.phone, customer.email, customer.address, customer.idCard, customer.customerType)
    
    if err != nil {
      return err
    }
  }

  // Create demo vehicles
  vehicles := []struct {
    code, chassis, plate, brand, model, variant string
    year, mileage int
    color, fuelType, transmission string
    purchasePrice float64
    status string
    purchasedFromCustomerID int
  }{
    {"VEH-001", "MHKA1234567890123", "B 1234 ABC", "Honda", "Civic", "RS", 2022, 15000, "White", "gasoline", "manual", 350000000, "purchased", 1},
    {"VEH-002", "WBAVA1234567890123", "B 5678 DEF", "BMW", "320i", "Sport", 2021, 25000, "Black", "gasoline", "automatic", 650000000, "in_repair", 2},
    {"VEH-003", "JTDKN1234567890123", "B 9012 GHI", "Toyota", "Camry", "Hybrid", 2023, 8000, "Silver", "hybrid", "cvt", 550000000, "ready_to_sell", 3},
    {"VEH-004", "KMHJ1234567890123", "B 3456 JKL", "Hyundai", "Tucson", "GLS", 2022, 20000, "Red", "gasoline", "automatic", 450000000, "purchased", 4},
    {"VEH-005", "JN1AZ1234567890123", "B 7890 MNO", "Nissan", "X-Trail", "XT", 2021, 35000, "Blue", "gasoline", "cvt", 400000000, "ready_to_sell", 5},
  }

  for _, vehicle := range vehicles {
    _, err = db.Exec(`
      INSERT INTO vehicles (
        vehicle_code, chassis_number, license_plate, brand, model, variant, 
        year, mileage, color, fuel_type, transmission, purchase_price, 
        status, purchased_from_customer_id, purchased_by_cashier, purchased_at,
        purchase_notes, condition_notes
      )
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, 2, CURRENT_TIMESTAMP, 
              'Vehicle purchased in good condition', 'Minor scratches on rear bumper')
    `, vehicle.code, vehicle.chassis, vehicle.plate, vehicle.brand, vehicle.model, vehicle.variant,
       vehicle.year, vehicle.mileage, vehicle.color, vehicle.fuelType, vehicle.transmission,
       vehicle.purchasePrice, vehicle.status, vehicle.purchasedFromCustomerID)
    
    if err != nil {
      return err
    }
  }

  log.Println("Demo data seeded successfully!")
  return nil
}
