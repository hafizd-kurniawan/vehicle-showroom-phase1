package entity

import "time"

type VehicleProfitability struct {
	VehicleID         int       `json:"vehicle_id" db:"vehicle_id"`
	VehicleCode       string    `json:"vehicle_code" db:"vehicle_code"`
	Brand             string    `json:"brand" db:"brand"`
	Model             string    `json:"model" db:"model"`
	Year              int       `json:"year" db:"year"`
	PurchasePrice     float64   `json:"purchase_price" db:"purchase_price"`
	TotalRepairCost   float64   `json:"total_repair_cost" db:"total_repair_cost"`
	FinalSellingPrice float64   `json:"final_selling_price" db:"final_selling_price"`
	Profit            float64   `json:"profit" db:"profit"`
	SoldAt            time.Time `json:"sold_at" db:"sold_at"`
}

type DateRangeRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}
