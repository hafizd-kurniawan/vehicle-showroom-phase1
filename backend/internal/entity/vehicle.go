package entity

import (
  "time"
)

type Vehicle struct {
  ID                      int        `json:"id" db:"id"`
  VehicleCode             string     `json:"vehicle_code" db:"vehicle_code"`
  ChassisNumber           string     `json:"chassis_number" db:"chassis_number"`
  LicensePlate            *string    `json:"license_plate" db:"license_plate"`
  Brand                   string     `json:"brand" db:"brand"`
  Model                   string     `json:"model" db:"model"`
  Variant                 *string    `json:"variant" db:"variant"`
  Year                    int        `json:"year" db:"year"`
  Color                   *string    `json:"color" db:"color"`
  Mileage                 *int       `json:"mileage" db:"mileage"`
  FuelType                *string    `json:"fuel_type" db:"fuel_type"`
  Transmission            *string    `json:"transmission" db:"transmission"`
  PurchasePrice           *float64   `json:"purchase_price" db:"purchase_price"`
  TotalRepairCost         float64    `json:"total_repair_cost" db:"total_repair_cost"`
  SuggestedSellingPrice   *float64   `json:"suggested_selling_price" db:"suggested_selling_price"`
  ApprovedSellingPrice    *float64   `json:"approved_selling_price" db:"approved_selling_price"`
  FinalSellingPrice       *float64   `json:"final_selling_price" db:"final_selling_price"`
  Status                  string     `json:"status" db:"status"`
  PurchasedFromCustomerID *int       `json:"purchased_from_customer_id" db:"purchased_from_customer_id"`
  SoldToCustomerID        *int       `json:"sold_to_customer_id" db:"sold_to_customer_id"`
  PurchasedByCashier      *int       `json:"purchased_by_cashier" db:"purchased_by_cashier"`
  SoldByCashier           *int       `json:"sold_by_cashier" db:"sold_by_cashier"`
  PriceApprovedByAdmin    *int       `json:"price_approved_by_admin" db:"price_approved_by_admin"`
  PurchasedAt             *time.Time `json:"purchased_at" db:"purchased_at"`
  SoldAt                  *time.Time `json:"sold_at" db:"sold_at"`
  CreatedAt               time.Time  `json:"created_at" db:"created_at"`
  UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`
  PurchaseNotes           *string    `json:"purchase_notes" db:"purchase_notes"`
  ConditionNotes          *string    `json:"condition_notes" db:"condition_notes"`
  
  // Joined fields
  PurchasedFromCustomer *Customer `json:"purchased_from_customer,omitempty"`
  SoldToCustomer        *Customer `json:"sold_to_customer,omitempty"`
}

type CreateVehicleRequest struct {
  ChassisNumber           string   `json:"chassis_number" binding:"required"`
  LicensePlate            *string  `json:"license_plate"`
  Brand                   string   `json:"brand" binding:"required"`
  Model                   string   `json:"model" binding:"required"`
  Variant                 *string  `json:"variant"`
  Year                    int      `json:"year" binding:"required,min=1900,max=2030"`
  Color                   *string  `json:"color"`
  Mileage                 *int     `json:"mileage"`
  FuelType                *string  `json:"fuel_type"`
  Transmission            *string  `json:"transmission"`
  PurchasePrice           *float64 `json:"purchase_price"`
  PurchasedFromCustomerID *int     `json:"purchased_from_customer_id"`
  PurchaseNotes           *string  `json:"purchase_notes"`
  ConditionNotes          *string  `json:"condition_notes"`
}

type UpdateVehicleRequest struct {
  LicensePlate            *string  `json:"license_plate"`
  Brand                   string   `json:"brand" binding:"required"`
  Model                   string   `json:"model" binding:"required"`
  Variant                 *string  `json:"variant"`
  Year                    int      `json:"year" binding:"required,min=1900,max=2030"`
  Color                   *string  `json:"color"`
  Mileage                 *int     `json:"mileage"`
  FuelType                *string  `json:"fuel_type"`
  Transmission            *string  `json:"transmission"`
  SuggestedSellingPrice   *float64 `json:"suggested_selling_price"`
  PurchaseNotes           *string  `json:"purchase_notes"`
  ConditionNotes          *string  `json:"condition_notes"`
}

type UpdateVehicleStatusRequest struct {
  Status string `json:"status" binding:"required,oneof=purchased in_repair ready_to_sell reserved sold"`
}

type VehicleListResponse struct {
  Vehicles []Vehicle `json:"vehicles"`
  Total    int       `json:"total"`
  Page     int       `json:"page"`
  Limit    int       `json:"limit"`
}
