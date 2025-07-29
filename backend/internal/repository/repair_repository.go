package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"vehicle-showroom/internal/entity"
)

type RepairRepository interface {
	Create(repair *entity.Repair) error
	GetByID(id int) (*entity.Repair, error)
	List(page, limit int, search, status string) ([]entity.Repair, int, error)
	Update(repair *entity.Repair) error
	UpdateStatus(id int, status string) error
	AddPart(repairPart *entity.RepairPart) error
	RemovePart(repairId, partId int) error
	GetRepairParts(repairId int) ([]entity.RepairPart, error)
	UpdateRepairCosts(repairId int) error
	GenerateRepairNumber() (string, error)
}

type repairRepository struct {
	db *sqlx.DB
}

func NewRepairRepository(db *sqlx.DB) RepairRepository {
	return &repairRepository{db: db}
}

func (r *repairRepository) Create(repair *entity.Repair) error {
	query := `
		INSERT INTO repairs (repair_number, vehicle_id, title, description, labor_cost,
		                    total_parts_cost, total_cost, status, mechanic_id, work_notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		repair.RepairNumber,
		repair.VehicleID,
		repair.Title,
		repair.Description,
		repair.LaborCost,
		repair.TotalPartsCost,
		repair.TotalCost,
		repair.Status,
		repair.MechanicID,
		repair.WorkNotes,
	).Scan(&repair.ID, &repair.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create repair: %w", err)
	}

	return nil
}

func (r *repairRepository) GetByID(id int) (*entity.Repair, error) {
	repair := &entity.Repair{}
	query := `
		SELECT r.id, r.repair_number, r.vehicle_id, r.title, r.description, r.labor_cost,
		       r.total_parts_cost, r.total_cost, r.status, r.mechanic_id, r.started_at,
		       r.completed_at, r.created_at, r.work_notes
		FROM repairs r
		WHERE r.id = $1
	`

	err := r.db.Get(repair, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get repair by id: %w", err)
	}

	// Load related data
	r.loadRepairRelatedData(repair)

	return repair, nil
}

func (r *repairRepository) List(page, limit int, search, status string) ([]entity.Repair, int, error) {
	offset := (page - 1) * limit

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		whereClause += fmt.Sprintf(" AND (r.repair_number ILIKE $%d OR r.title ILIKE $%d OR v.brand ILIKE $%d OR v.model ILIKE $%d)",
			argIndex, argIndex+1, argIndex+2, argIndex+3)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
		argIndex += 4
	}

	if status != "" {
		whereClause += fmt.Sprintf(" AND r.status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM repairs r
		LEFT JOIN vehicles v ON r.vehicle_id = v.id
		%s
	`, whereClause)

	var total int
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get repair count: %w", err)
	}

	// Get repairs
	query := fmt.Sprintf(`
		SELECT r.id, r.repair_number, r.vehicle_id, r.title, r.description, r.labor_cost,
		       r.total_parts_cost, r.total_cost, r.status, r.mechanic_id, r.started_at,
		       r.completed_at, r.created_at, r.work_notes
		FROM repairs r
		LEFT JOIN vehicles v ON r.vehicle_id = v.id
		%s
		ORDER BY r.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	var repairs []entity.Repair
	err = r.db.Select(&repairs, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list repairs: %w", err)
	}

	// Load related data for each repair
	for i := range repairs {
		r.loadRepairRelatedData(&repairs[i])
	}

	return repairs, total, nil
}

func (r *repairRepository) Update(repair *entity.Repair) error {
	query := `
		UPDATE repairs
		SET title = $1, description = $2, labor_cost = $3, mechanic_id = $4, work_notes = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(query, repair.Title, repair.Description, repair.LaborCost,
		repair.MechanicID, repair.WorkNotes, repair.ID)
	if err != nil {
		return fmt.Errorf("failed to update repair: %w", err)
	}

	return nil
}

func (r *repairRepository) UpdateStatus(id int, status string) error {
	var query string
	var args []interface{}

	if status == "in_progress" {
		query = `UPDATE repairs SET status = $1, started_at = CURRENT_TIMESTAMP WHERE id = $2`
		args = []interface{}{status, id}
	} else if status == "completed" {
		query = `UPDATE repairs SET status = $1, completed_at = CURRENT_TIMESTAMP WHERE id = $2`
		args = []interface{}{status, id}
	} else {
		query = `UPDATE repairs SET status = $1 WHERE id = $2`
		args = []interface{}{status, id}
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update repair status: %w", err)
	}

	return nil
}

func (r *repairRepository) AddPart(repairPart *entity.RepairPart) error {
	query := `
		INSERT INTO repair_parts (repair_id, spare_part_id, quantity_used, unit_cost, total_cost, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, used_at
	`

	err := r.db.QueryRow(
		query,
		repairPart.RepairID,
		repairPart.SparePartID,
		repairPart.QuantityUsed,
		repairPart.UnitCost,
		repairPart.TotalCost,
		repairPart.Notes,
	).Scan(&repairPart.ID, &repairPart.UsedAt)

	if err != nil {
		return fmt.Errorf("failed to add part to repair: %w", err)
	}

	return nil
}

func (r *repairRepository) RemovePart(repairId, partId int) error {
	query := `DELETE FROM repair_parts WHERE repair_id = $1 AND id = $2`

	_, err := r.db.Exec(query, repairId, partId)
	if err != nil {
		return fmt.Errorf("failed to remove part from repair: %w", err)
	}

	return nil
}

func (r *repairRepository) GetRepairParts(repairId int) ([]entity.RepairPart, error) {
	var parts []entity.RepairPart
	query := `
		SELECT rp.id, rp.repair_id, rp.spare_part_id, rp.quantity_used, rp.unit_cost,
		       rp.total_cost, rp.used_at, rp.notes
		FROM repair_parts rp
		WHERE rp.repair_id = $1
		ORDER BY rp.used_at DESC
	`

	err := r.db.Select(&parts, query, repairId)
	if err != nil {
		return nil, fmt.Errorf("failed to get repair parts: %w", err)
	}

	// Load spare part details for each part
	for i := range parts {
		sparePart := &entity.SparePart{}
		sparePartQuery := `
			SELECT id, part_code, name, description, brand, unit_measure
			FROM spare_parts WHERE id = $1
		`
		if err := r.db.Get(sparePart, sparePartQuery, parts[i].SparePartID); err == nil {
			parts[i].SparePart = sparePart
		}
	}

	return parts, nil
}

func (r *repairRepository) UpdateRepairCosts(repairId int) error {
	// Calculate total parts cost
	var totalPartsCost sql.NullFloat64
	partsQuery := `SELECT COALESCE(SUM(total_cost), 0) FROM repair_parts WHERE repair_id = $1`
	err := r.db.Get(&totalPartsCost, partsQuery, repairId)
	if err != nil {
		return fmt.Errorf("failed to calculate total parts cost: %w", err)
	}

	// Get labor cost
	var laborCost sql.NullFloat64
	laborQuery := `SELECT COALESCE(labor_cost, 0) FROM repairs WHERE id = $1`
	err = r.db.Get(&laborCost, laborQuery, repairId)
	if err != nil {
		return fmt.Errorf("failed to get labor cost: %w", err)
	}

	// Update repair costs
	totalCost := laborCost.Float64 + totalPartsCost.Float64
	updateQuery := `
		UPDATE repairs 
		SET total_parts_cost = $1, total_cost = $2 
		WHERE id = $3
	`

	_, err = r.db.Exec(updateQuery, totalPartsCost.Float64, totalCost, repairId)
	if err != nil {
		return fmt.Errorf("failed to update repair costs: %w", err)
	}

	return nil
}

func (r *repairRepository) GenerateRepairNumber() (string, error) {
	today := time.Now().Format("20060102")
	var lastNumber string
	query := `
		SELECT repair_number 
		FROM repairs 
		WHERE repair_number LIKE 'REP-' || $1 || '-%' 
		ORDER BY repair_number DESC 
		LIMIT 1
	`

	err := r.db.Get(&lastNumber, query, today)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get last repair number: %w", err)
	}

	nextNumber := 1
	if lastNumber != "" {
		parts := strings.Split(lastNumber, "-")
		if len(parts) == 3 {
			var num int
			fmt.Sscanf(parts[2], "%d", &num)
			nextNumber = num + 1
		}
	}

	return fmt.Sprintf("REP-%s-%03d", today, nextNumber), nil
}

func (r *repairRepository) loadRepairRelatedData(repair *entity.Repair) {
	// Load vehicle
	vehicle := &entity.Vehicle{}
	vehicleQuery := `
		SELECT id, vehicle_code, brand, model, variant, year
		FROM vehicles WHERE id = $1
	`
	if err := r.db.Get(vehicle, vehicleQuery, repair.VehicleID); err == nil {
		repair.Vehicle = vehicle
	}

	// Load mechanic
	if repair.MechanicID != nil {
		mechanic := &entity.User{}
		mechanicQuery := `
			SELECT id, username, full_name, role
			FROM users WHERE id = $1
		`
		if err := r.db.Get(mechanic, mechanicQuery, *repair.MechanicID); err == nil {
			repair.Mechanic = mechanic
		}
	}

	// Load repair parts
	parts, err := r.GetRepairParts(repair.ID)
	if err == nil {
		repair.RepairParts = parts
	}
}
