package database

import (
  "github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) error {
  migrations := []string{
    createUsersTable,
    createUserSessionsTable,
    createCustomersTable,
    createVehiclesTable,
    createVehicleImagesTable,
    createSparePartsTable,
    createPurchaseTransactionsTable,
    createSalesTransactionsTable,
    createRepairsTable,
    createRepairPartsTable,
    createStockMovementsTable,
  }

  for _, migration := range migrations {
    if _, err := db.Exec(migration); err != nil {
      return err
    }
  }

  return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  full_name VARCHAR(100) NOT NULL,
  phone VARCHAR(20),
  role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'mechanic', 'cashier')),
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createUserSessionsTable = `
CREATE TABLE IF NOT EXISTS user_sessions (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  session_token VARCHAR(255) UNIQUE NOT NULL,
  login_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  logout_at TIMESTAMP,
  ip_address VARCHAR(45),
  is_active BOOLEAN DEFAULT true
);
`

const createCustomersTable = `
CREATE TABLE IF NOT EXISTS customers (
  id SERIAL PRIMARY KEY,
  customer_code VARCHAR(20) UNIQUE NOT NULL,
  name VARCHAR(100) NOT NULL,
  phone VARCHAR(20),
  email VARCHAR(100),
  address TEXT,
  id_card_number VARCHAR(50),
  type VARCHAR(20) NOT NULL CHECK (type IN ('individual', 'corporate')),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by INTEGER REFERENCES users(id),
  is_active BOOLEAN DEFAULT true
);
`

const createVehiclesTable = `
CREATE TABLE IF NOT EXISTS vehicles (
  id SERIAL PRIMARY KEY,
  vehicle_code VARCHAR(20) UNIQUE NOT NULL,
  chassis_number VARCHAR(50) UNIQUE NOT NULL,
  license_plate VARCHAR(20),
  brand VARCHAR(50) NOT NULL,
  model VARCHAR(50) NOT NULL,
  variant VARCHAR(50),
  year INTEGER NOT NULL,
  color VARCHAR(30),
  mileage INTEGER,
  fuel_type VARCHAR(20) CHECK (fuel_type IN ('gasoline', 'diesel', 'electric', 'hybrid')),
  transmission VARCHAR(20) CHECK (transmission IN ('manual', 'automatic', 'cvt')),
  purchase_price DECIMAL(15,2),
  total_repair_cost DECIMAL(15,2) DEFAULT 0,
  suggested_selling_price DECIMAL(15,2),
  approved_selling_price DECIMAL(15,2),
  final_selling_price DECIMAL(15,2),
  status VARCHAR(20) DEFAULT 'purchased' CHECK (status IN ('purchased', 'in_repair', 'ready_to_sell', 'reserved', 'sold')),
  purchased_from_customer_id INTEGER REFERENCES customers(id),
  sold_to_customer_id INTEGER REFERENCES customers(id),
  purchased_by_cashier INTEGER REFERENCES users(id),
  sold_by_cashier INTEGER REFERENCES users(id),
  price_approved_by_admin INTEGER REFERENCES users(id),
  purchased_at TIMESTAMP,
  sold_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  purchase_notes TEXT,
  condition_notes TEXT
);
`

const createVehicleImagesTable = `
CREATE TABLE IF NOT EXISTS vehicle_images (
  id SERIAL PRIMARY KEY,
  vehicle_id INTEGER REFERENCES vehicles(id) ON DELETE CASCADE,
  image_path VARCHAR(255) NOT NULL,
  image_type VARCHAR(20) CHECK (image_type IN ('front', 'back', 'left', 'right', 'interior', 'engine', 'dashboard', 'damage', 'other')),
  description TEXT,
  is_primary BOOLEAN DEFAULT false,
  uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  uploaded_by INTEGER REFERENCES users(id)
);
`

const createSparePartsTable = `
CREATE TABLE IF NOT EXISTS spare_parts (
  id SERIAL PRIMARY KEY,
  part_code VARCHAR(20) UNIQUE NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  brand VARCHAR(50),
  cost_price DECIMAL(15,2) NOT NULL,
  selling_price DECIMAL(15,2) NOT NULL,
  stock_quantity INTEGER DEFAULT 0,
  min_stock_level INTEGER DEFAULT 0,
  unit_measure VARCHAR(20),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_active BOOLEAN DEFAULT true
);
`

const createPurchaseTransactionsTable = `
CREATE TABLE IF NOT EXISTS purchase_transactions (
  id SERIAL PRIMARY KEY,
  transaction_number VARCHAR(30) UNIQUE NOT NULL,
  invoice_number VARCHAR(30) UNIQUE NOT NULL,
  vehicle_id INTEGER REFERENCES vehicles(id),
  customer_id INTEGER REFERENCES customers(id),
  vehicle_price DECIMAL(15,2) NOT NULL,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) NOT NULL,
  payment_method VARCHAR(20) CHECK (payment_method IN ('cash', 'transfer', 'check')),
  payment_reference VARCHAR(100),
  transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  cashier_id INTEGER REFERENCES users(id),
  status VARCHAR(20) DEFAULT 'completed' CHECK (status IN ('completed', 'cancelled')),
  notes TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createSalesTransactionsTable = `
CREATE TABLE IF NOT EXISTS sales_transactions (
  id SERIAL PRIMARY KEY,
  transaction_number VARCHAR(30) UNIQUE NOT NULL,
  invoice_number VARCHAR(30) UNIQUE NOT NULL,
  vehicle_id INTEGER REFERENCES vehicles(id),
  customer_id INTEGER REFERENCES customers(id),
  vehicle_price DECIMAL(15,2) NOT NULL,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) NOT NULL,
  payment_method VARCHAR(20) CHECK (payment_method IN ('cash', 'transfer', 'check', 'credit')),
  payment_reference VARCHAR(100),
  transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  cashier_id INTEGER REFERENCES users(id),
  status VARCHAR(20) DEFAULT 'completed' CHECK (status IN ('completed', 'cancelled')),
  notes TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createRepairsTable = `
CREATE TABLE IF NOT EXISTS repairs (
  id SERIAL PRIMARY KEY,
  repair_number VARCHAR(30) UNIQUE NOT NULL,
  vehicle_id INTEGER REFERENCES vehicles(id),
  title VARCHAR(100) NOT NULL,
  description TEXT,
  labor_cost DECIMAL(15,2) DEFAULT 0,
  total_parts_cost DECIMAL(15,2) DEFAULT 0,
  total_cost DECIMAL(15,2) DEFAULT 0,
  status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled')),
  mechanic_id INTEGER REFERENCES users(id),
  started_at TIMESTAMP,
  completed_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  work_notes TEXT
);
`

const createRepairPartsTable = `
CREATE TABLE IF NOT EXISTS repair_parts (
  id SERIAL PRIMARY KEY,
  repair_id INTEGER REFERENCES repairs(id) ON DELETE CASCADE,
  spare_part_id INTEGER REFERENCES spare_parts(id),
  quantity_used INTEGER NOT NULL,
  unit_cost DECIMAL(15,2) NOT NULL,
  total_cost DECIMAL(15,2) NOT NULL,
  used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  notes TEXT
);
`

const createStockMovementsTable = `
CREATE TABLE IF NOT EXISTS stock_movements (
  id SERIAL PRIMARY KEY,
  spare_part_id INTEGER REFERENCES spare_parts(id),
  movement_type VARCHAR(20) CHECK (movement_type IN ('in', 'out', 'adjustment')),
  reference_type VARCHAR(20) CHECK (reference_type IN ('repair', 'purchase', 'sales', 'adjustment')),
  reference_id INTEGER,
  quantity_before INTEGER NOT NULL,
  quantity_moved INTEGER NOT NULL,
  quantity_after INTEGER NOT NULL,
  movement_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  processed_by INTEGER REFERENCES users(id),
  notes TEXT
);
`
