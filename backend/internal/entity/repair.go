package entity

import "time"

type Repair struct {
	ID              int          `json:"id" db:"id"`
	RepairNumber    string       `json:"repair_number" db:"repair_number"`
	VehicleID       int          `json:"vehicle_id" db:"vehicle_id"`
	Title           string       `json:"title" db:"title"`
	Description     *string      `json:"description" db:"description"`
	LaborCost       float64      `json:"labor_cost" db:"labor_cost"`
	TotalPartsCost  float64      `json:"total_parts_cost" db:"total_parts_cost"`
	TotalCost       float64      `json:"total_cost" db:"total_cost"`
	Status          string       `json:"status" db:"status"`
	MechanicID      *int         `json:"mechanic_id" db:"mechanic_id"`
	StartedAt       *time.Time   `json:"started_at" db:"started_at"`
	CompletedAt     *time.Time   `json:"completed_at" db:"completed_at"`
	CreatedAt       time.Time    `json:"created_at" db:"created_at"`
	WorkNotes       *string      `json:"work_notes" db:"work_notes"`
	Vehicle         *Vehicle     `json:"vehicle,omitempty"`
	Mechanic        *User        `json:"mechanic,omitempty"`
	RepairParts     []RepairPart `json:"repair_parts,omitempty"`
}

type RepairPart struct {
	ID          int       `json:"id" db:"id"`
	RepairID    int       `json:"repair_id" db:"repair_id"`
	SparePartID int       `json:"spare_part_id" db:"spare_part_id"`
	QuantityUsed int      `json:"quantity_used" db:"quantity_used"`
	UnitCost    float64   `json:"unit_cost" db:"unit_cost"`
	TotalCost   float64   `json:"total_cost" db:"total_cost"`
	UsedAt      time.Time `json:"used_at" db:"used_at"`
	Notes       *string   `json:"notes" db:"notes"`
	SparePart   *SparePart `json:"spare_part,omitempty"`
}

type CreateRepairRequest struct {
	VehicleID   int     `json:"vehicle_id" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	MechanicID  *int    `json:"mechanic_id"`
}

type UpdateRepairRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description *string  `json:"description"`
	LaborCost   *float64 `json:"labor_cost"`
	MechanicID  *int     `json:"mechanic_id"`
	WorkNotes   *string  `json:"work_notes"`
}

type UpdateRepairStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending in_progress completed cancelled"`
}

type AddPartToRepairRequest struct {
	SparePartID int `json:"spare_part_id" binding:"required"`
	Quantity    int `json:"quantity" binding:"required,min=1"`
}

type RepairListResponse struct {
	Repairs []Repair `json:"repairs"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
}
