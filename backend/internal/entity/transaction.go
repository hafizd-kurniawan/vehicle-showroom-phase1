package entity

import (
  "time"
)

type PurchaseTransaction struct {
  ID                int        `json:"id" db:"id"`
  TransactionNumber string     `json:"transaction_number" db:"transaction_number"`
  InvoiceNumber     string     `json:"invoice_number" db:"invoice_number"`
  VehicleID         int        `json:"vehicle_id" db:"vehicle_id"`
  CustomerID        int        `json:"customer_id" db:"customer_id"`
  VehiclePrice      float64    `json:"vehicle_price" db:"vehicle_price"`
  TaxAmount         float64    `json:"tax_amount" db:"tax_amount"`
  TotalAmount       float64    `json:"total_amount" db:"total_amount"`
  PaymentMethod     string     `json:"payment_method" db:"payment_method"`
  PaymentReference  *string    `json:"payment_reference" db:"payment_reference"`
  TransactionDate   time.Time  `json:"transaction_date" db:"transaction_date"`
  CashierID         int        `json:"cashier_id" db:"cashier_id"`
  Status            string     `json:"status" db:"status"`
  Notes             *string    `json:"notes" db:"notes"`
  CreatedAt         time.Time  `json:"created_at" db:"created_at"`
  
  // Joined fields
  Vehicle  *Vehicle  `json:"vehicle,omitempty"`
  Customer *Customer `json:"customer,omitempty"`
  Cashier  *User     `json:"cashier,omitempty"`
}

type SalesTransaction struct {
  ID                int        `json:"id" db:"id"`
  TransactionNumber string     `json:"transaction_number" db:"transaction_number"`
  InvoiceNumber     string     `json:"invoice_number" db:"invoice_number"`
  VehicleID         int        `json:"vehicle_id" db:"vehicle_id"`
  CustomerID        int        `json:"customer_id" db:"customer_id"`
  VehiclePrice      float64    `json:"vehicle_price" db:"vehicle_price"`
  TaxAmount         float64    `json:"tax_amount" db:"tax_amount"`
  DiscountAmount    float64    `json:"discount_amount" db:"discount_amount"`
  TotalAmount       float64    `json:"total_amount" db:"total_amount"`
  PaymentMethod     string     `json:"payment_method" db:"payment_method"`
  PaymentReference  *string    `json:"payment_reference" db:"payment_reference"`
  TransactionDate   time.Time  `json:"transaction_date" db:"transaction_date"`
  CashierID         int        `json:"cashier_id" db:"cashier_id"`
  Status            string     `json:"status" db:"status"`
  Notes             *string    `json:"notes" db:"notes"`
  CreatedAt         time.Time  `json:"created_at" db:"created_at"`
  
  // Joined fields
  Vehicle  *Vehicle  `json:"vehicle,omitempty"`
  Customer *Customer `json:"customer,omitempty"`
  Cashier  *User     `json:"cashier,omitempty"`
}

type CreatePurchaseTransactionRequest struct {
  VehicleID        int     `json:"vehicle_id" binding:"required"`
  CustomerID       int     `json:"customer_id" binding:"required"`
  VehiclePrice     float64 `json:"vehicle_price" binding:"required,min=0"`
  TaxAmount        float64 `json:"tax_amount" binding:"min=0"`
  PaymentMethod    string  `json:"payment_method" binding:"required,oneof=cash transfer check"`
  PaymentReference *string `json:"payment_reference"`
  Notes            *string `json:"notes"`
}

type CreateSalesTransactionRequest struct {
  VehicleID        int     `json:"vehicle_id" binding:"required"`
  CustomerID       int     `json:"customer_id" binding:"required"`
  VehiclePrice     float64 `json:"vehicle_price" binding:"required,min=0"`
  TaxAmount        float64 `json:"tax_amount" binding:"min=0"`
  DiscountAmount   float64 `json:"discount_amount" binding:"min=0"`
  PaymentMethod    string  `json:"payment_method" binding:"required,oneof=cash transfer check credit"`
  PaymentReference *string `json:"payment_reference"`
  Notes            *string `json:"notes"`
}

type TransactionListResponse struct {
  Transactions interface{} `json:"transactions"`
  Total        int         `json:"total"`
  Page         int         `json:"page"`
  Limit        int         `json:"limit"`
}

type DashboardStats struct {
  TotalVehicles       int     `json:"total_vehicles"`
  VehiclesForSale     int     `json:"vehicles_for_sale"`
  VehiclesInRepair    int     `json:"vehicles_in_repair"`
  VehiclesSold        int     `json:"vehicles_sold"`
  TotalCustomers      int     `json:"total_customers"`
  TodayPurchases      int     `json:"today_purchases"`
  TodaySales          int     `json:"today_sales"`
  TodayRevenue        float64 `json:"today_revenue"`
  MonthlyRevenue      float64 `json:"monthly_revenue"`
  TotalProfit         float64 `json:"total_profit"`
}
