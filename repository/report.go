package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"

	"go.uber.org/zap"
)

type ReportRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}
type ReportRepoInterface interface {
	Report() (*model.Report, error)
}

// constructor
func NewReportRepo(db database.PgxIface,
	log *zap.Logger) ReportRepoInterface {
	return &ReportRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *ReportRepo) Report() (*model.Report, error) {
	query := `
        SELECT 
            COUNT(sales_item_id),
            COALESCE(SUM(quantity), 0),          
            COALESCE(SUM(quantity * price), 0)  
        FROM sales_item
        WHERE deleted_at IS NULL
    `

	var report model.Report
	err := r.DB.QueryRow(context.Background(), query).Scan(
		&report.TotalTransactions,
		&report.TotalItemsSold,
		&report.TotalRevenue,
	)

	if err != nil {
		return nil, err
	}

	return &report, nil
}
