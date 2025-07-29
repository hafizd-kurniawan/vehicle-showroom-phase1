package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"vehicle-showroom/internal/entity"
)

type SparePartRepository interface {
	Create(sparePart *entity.SparePart) error
	GetByID(id int) (*entity.SparePart, error)
	GetByCode(code string) (*entity.SparePart, error)
	List(page, limit int, search string) ([]entity.SparePart, int, error)
	Update(sparePart *entity.SparePart) error
	Delete(id int) error
	UpdateStock(id int, quantity int) error
	GeneratePartCode() (string, error)
}

type sparePartRepository struct {
	db *sqlx.DB
}

func NewSparePartRepository(db *sqlx.DB) SparePartRepository {
	return &sparePartRepository{db: db}
}

func (r *sparePartRepository) Create(sparePart *entity.SparePart) error {
	query := `
		INSERT INTO spare_parts (part_code, name, description, brand, cost_price, selling_price, 
		                        stock_quantity, min_stock_level, unit_measure, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		sparePart.PartCode,
		sparePart.Name,
		sparePart.Description,
		sparePart.Brand,
		sparePart.CostPrice,
		sparePart.SellingPrice,
		sparePart.StockQuantity,
		sparePart.MinStockLevel,
		sparePart.UnitMeasure,
		sparePart.IsActive,
	).Scan(&sparePart.ID, &sparePart.CreatedAt, &sparePart.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create spare part: %w", err)
	}

	return nil
}

func (r *sparePartRepository) GetByID(id int) (*entity.SparePart, error) {
	sparePart := &entity.SparePart{}
	query := `
		SELECT id, part_code, name, description, brand, cost_price, selling_price,
		       stock_quantity, min_stock_level, unit_measure, created_at, updated_at, is_active
		FROM spare_parts
		WHERE id = $1 AND is_active = true
	`

	err := r.db.Get(sparePart, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get spare part by id: %w", err)
	}

	return sparePart, nil
}

func (r *sparePartRepository) GetByCode(code string) (*entity.SparePart, error) {
	sparePart := &entity.SparePart{}
	query := `
		SELECT id, part_code, name, description, brand, cost_price, selling_price,
		       stock_quantity, min_stock_level, unit_measure, created_at, updated_at, is_active
		FROM spare_parts
		WHERE part_code = $1 AND is_active = true
	`

	err := r.db.Get(sparePart, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get spare part by code: %w", err)
	}

	return sparePart, nil
}

func (r *sparePartRepository) List(page, limit int, search string) ([]entity.SparePart, int, error) {
	offset := (page - 1) * limit

	whereClause := "WHERE is_active = true"
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR part_code ILIKE $%d OR brand ILIKE $%d OR description ILIKE $%d)",
			argIndex, argIndex+1, argIndex+2, argIndex+3)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
		argIndex += 4
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM spare_parts %s", whereClause)
	var total int
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get spare part count: %w", err)
	}

	// Get spare parts
	query := fmt.Sprintf(`
		SELECT id, part_code, name, description, brand, cost_price, selling_price,
		       stock_quantity, min_stock_level, unit_measure, created_at, updated_at, is_active
		FROM spare_parts
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	var spareParts []entity.SparePart
	err = r.db.Select(&spareParts, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list spare parts: %w", err)
	}

	return spareParts, total, nil
}

func (r *sparePartRepository) Update(sparePart *entity.SparePart) error {
	query := `
		UPDATE spare_parts
		SET name = $1, description = $2, brand = $3, cost_price = $4, selling_price = $5,
		    min_stock_level = $6, unit_measure = $7, is_active = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
	`

	_, err := r.db.Exec(query, sparePart.Name, sparePart.Description, sparePart.Brand,
		sparePart.CostPrice, sparePart.SellingPrice, sparePart.MinStockLevel,
		sparePart.UnitMeasure, sparePart.IsActive, sparePart.ID)
	if err != nil {
		return fmt.Errorf("failed to update spare part: %w", err)
	}

	return nil
}

func (r *sparePartRepository) Delete(id int) error {
	query := `UPDATE spare_parts SET is_active = false WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete spare part: %w", err)
	}

	return nil
}

func (r *sparePartRepository) UpdateStock(id int, quantity int) error {
	query := `UPDATE spare_parts SET stock_quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	_, err := r.db.Exec(query, quantity, id)
	if err != nil {
		return fmt.Errorf("failed to update spare part stock: %w", err)
	}

	return nil
}

func (r *sparePartRepository) GeneratePartCode() (string, error) {
	var lastCode string
	query := `
		SELECT part_code 
		FROM spare_parts 
		WHERE part_code LIKE 'PART-%' 
		ORDER BY part_code DESC 
		LIMIT 1
	`

	err := r.db.Get(&lastCode, query)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get last part code: %w", err)
	}

	nextNumber := 1
	if lastCode != "" {
		// Extract number from PART-XXX format
		parts := strings.Split(lastCode, "-")
		if len(parts) == 2 {
			var num int
			fmt.Sscanf(parts[1], "%d", &num)
			nextNumber = num + 1
		}
	}

	return fmt.Sprintf("PART-%03d", nextNumber), nil
}
