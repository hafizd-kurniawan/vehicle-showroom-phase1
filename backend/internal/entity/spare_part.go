package entity

import "time"

type SparePart struct {
	ID             int       `json:"id" db:"id"`
	PartCode       string    `json:"part_code" db:"part_code"`
	Name           string    `json:"name" db:"name"`
	Description    *string   `json:"description" db:"description"`
	Brand          *string   `json:"brand" db:"brand"`
	CostPrice      float64   `json:"cost_price" db:"cost_price"`
	SellingPrice   float64   `json:"selling_price" db:"selling_price"`
	StockQuantity  int       `json:"stock_quantity" db:"stock_quantity"`
	MinStockLevel  int       `json:"min_stock_level" db:"min_stock_level"`
	UnitMeasure    *string   `json:"unit_measure" db:"unit_measure"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	IsActive       bool      `json:"is_active" db:"is_active"`
}

type CreateSparePartRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   *string  `json:"description"`
	Brand         *string  `json:"brand"`
	CostPrice     float64  `json:"cost_price" binding:"required,min=0"`
	SellingPrice  float64  `json:"selling_price" binding:"required,min=0"`
	StockQuantity int      `json:"stock_quantity" binding:"min=0"`
	MinStockLevel int      `json:"min_stock_level" binding:"min=0"`
	UnitMeasure   *string  `json:"unit_measure"`
}

type UpdateSparePartRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   *string  `json:"description"`
	Brand         *string  `json:"brand"`
	CostPrice     float64  `json:"cost_price" binding:"required,min=0"`
	SellingPrice  float64  `json:"selling_price" binding:"required,min=0"`
	MinStockLevel int      `json:"min_stock_level" binding:"min=0"`
	UnitMeasure   *string  `json:"unit_measure"`
	IsActive      *bool    `json:"is_active"`
}
