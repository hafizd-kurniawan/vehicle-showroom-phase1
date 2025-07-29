package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"vehicle-showroom/internal/entity"
)

type ReportRepository interface {
	GetVehicleProfitabilityReport(startDate, endDate time.Time) ([]entity.VehicleProfitability, error)
	GetSalesReport(startDate, endDate time.Time) ([]entity.SalesTransaction, error)
	GetPurchaseReport(startDate, endDate time.Time) ([]entity.PurchaseTransaction, error)
}

type reportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetVehicleProfitabilityReport(startDate, endDate time.Time) ([]entity.VehicleProfitability, error) {
	var report []entity.VehicleProfitability
	query := `
        SELECT
            v.id AS vehicle_id,
            v.vehicle_code,
            v.brand,
            v.model,
            v.year,
            COALESCE(v.purchase_price, 0) AS purchase_price,
            COALESCE(v.total_repair_cost, 0) AS total_repair_cost,
            COALESCE(v.final_selling_price, 0) AS final_selling_price,
            (COALESCE(v.final_selling_price, 0) - COALESCE(v.purchase_price, 0) - COALESCE(v.total_repair_cost, 0)) AS profit,
            v.sold_at
        FROM vehicles v
        WHERE v.status = 'sold' AND v.sold_at BETWEEN $1 AND $2
        ORDER BY v.sold_at DESC
    `
	err := r.db.Select(&report, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle profitability report: %w", err)
	}
	return report, nil
}

func (r *reportRepository) GetSalesReport(startDate, endDate time.Time) ([]entity.SalesTransaction, error) {
	var transactions []entity.SalesTransaction
	query := `
        SELECT
            st.id, st.transaction_number, st.invoice_number, st.vehicle_id, st.customer_id,
            st.vehicle_price, st.tax_amount, st.discount_amount, st.total_amount,
            st.payment_method, st.payment_reference, st.transaction_date, st.cashier_id,
            st.status, st.notes, st.created_at
        FROM sales_transactions st
        WHERE st.transaction_date BETWEEN $1 AND $2
        ORDER BY st.transaction_date DESC
    `
	err := r.db.Select(&transactions, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales report: %w", err)
	}

	// Load related data for each transaction
	for i := range transactions {
		r.loadSalesRelatedData(&transactions[i])
	}

	return transactions, nil
}

func (r *reportRepository) GetPurchaseReport(startDate, endDate time.Time) ([]entity.PurchaseTransaction, error) {
	var transactions []entity.PurchaseTransaction
	query := `
        SELECT
            pt.id, pt.transaction_number, pt.invoice_number, pt.vehicle_id, pt.customer_id,
            pt.vehicle_price, pt.tax_amount, pt.total_amount, pt.payment_method,
            pt.payment_reference, pt.transaction_date, pt.cashier_id, pt.status,
            pt.notes, pt.created_at
        FROM purchase_transactions pt
        WHERE pt.transaction_date BETWEEN $1 AND $2
        ORDER BY pt.transaction_date DESC
    `
	err := r.db.Select(&transactions, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase report: %w", err)
	}

	// Load related data for each transaction
	for i := range transactions {
		r.loadPurchaseRelatedData(&transactions[i])
	}

	return transactions, nil
}

// Helper functions to load related data (can be moved to a shared place if needed)
func (r *reportRepository) loadPurchaseRelatedData(tx *entity.PurchaseTransaction) {
	// Load vehicle
	vehicle := &entity.Vehicle{}
	vehicleQuery := `SELECT id, vehicle_code, brand, model FROM vehicles WHERE id = $1`
	if err := r.db.Get(vehicle, vehicleQuery, tx.VehicleID); err == nil {
		tx.Vehicle = vehicle
	}

	// Load customer
	customer := &entity.Customer{}
	customerQuery := `SELECT id, name FROM customers WHERE id = $1`
	if err := r.db.Get(customer, customerQuery, tx.CustomerID); err == nil {
		tx.Customer = customer
	}

	// Load cashier
	cashier := &entity.User{}
	cashierQuery := `SELECT id, full_name FROM users WHERE id = $1`
	if err := r.db.Get(cashier, cashierQuery, tx.CashierID); err == nil {
		tx.Cashier = cashier
	}
}

func (r *reportRepository) loadSalesRelatedData(tx *entity.SalesTransaction) {
	// Load vehicle
	vehicle := &entity.Vehicle{}
	vehicleQuery := `SELECT id, vehicle_code, brand, model FROM vehicles WHERE id = $1`
	if err := r.db.Get(vehicle, vehicleQuery, tx.VehicleID); err == nil {
		tx.Vehicle = vehicle
	}

	// Load customer
	customer := &entity.Customer{}
	customerQuery := `SELECT id, name FROM customers WHERE id = $1`
	if err := r.db.Get(customer, customerQuery, tx.CustomerID); err == nil {
		tx.Customer = customer
	}

	// Load cashier
	cashier := &entity.User{}
	cashierQuery := `SELECT id, full_name FROM users WHERE id = $1`
	if err := r.db.Get(cashier, cashierQuery, tx.CashierID); err == nil {
		tx.Cashier = cashier
	}
}
