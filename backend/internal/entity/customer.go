package entity

import (
  "time"
)

type Customer struct {
  ID           int       `json:"id" db:"id"`
  CustomerCode string    `json:"customer_code" db:"customer_code"`
  Name         string    `json:"name" db:"name"`
  Phone        *string   `json:"phone" db:"phone"`
  Email        *string   `json:"email" db:"email"`
  Address      *string   `json:"address" db:"address"`
  IDCardNumber *string   `json:"id_card_number" db:"id_card_number"`
  Type         string    `json:"type" db:"type"`
  CreatedAt    time.Time `json:"created_at" db:"created_at"`
  UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
  CreatedBy    *int      `json:"created_by" db:"created_by"`
  IsActive     bool      `json:"is_active" db:"is_active"`
}

type CreateCustomerRequest struct {
  Name         string  `json:"name" binding:"required"`
  Phone        *string `json:"phone"`
  Email        *string `json:"email"`
  Address      *string `json:"address"`
  IDCardNumber *string `json:"id_card_number"`
  Type         string  `json:"type" binding:"required,oneof=individual corporate"`
}

type UpdateCustomerRequest struct {
  Name         string  `json:"name" binding:"required"`
  Phone        *string `json:"phone"`
  Email        *string `json:"email"`
  Address      *string `json:"address"`
  IDCardNumber *string `json:"id_card_number"`
  Type         string  `json:"type" binding:"required,oneof=individual corporate"`
}

type CustomerListResponse struct {
  Customers []Customer `json:"customers"`
  Total     int        `json:"total"`
  Page      int        `json:"page"`
  Limit     int        `json:"limit"`
}
