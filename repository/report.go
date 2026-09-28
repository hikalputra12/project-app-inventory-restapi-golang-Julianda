package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ReportRepo struct {
	DB     database.PgxIface
	Redis  *redis.Client
	Logger *zap.Logger
}

type ReportRepoInterface interface {
	Report(ctx context.Context) (*model.Report, error)
	InvalidateReportCache(ctx context.Context)
}

func NewReportRepo(db database.PgxIface, rdb *redis.Client, log *zap.Logger) ReportRepoInterface {
	return &ReportRepo{
		DB:     db,
		Redis:  rdb,
		Logger: log,
	}
}

func (r *ReportRepo) Report(ctx context.Context) (*model.Report, error) {
	cacheKey := "inventory:report:summary"

	// Try reading from Redis cache first
	if r.Redis != nil {
		cachedData, err := r.Redis.Get(ctx, cacheKey).Result()
		if err == nil && cachedData != "" {
			var cachedReport model.Report
			if err := json.Unmarshal([]byte(cachedData), &cachedReport); err == nil {
				return &cachedReport, nil
			}
		}
	}

	query := `
        SELECT 
            COUNT(sales_item_id),
            COALESCE(SUM(quantity), 0),          
            COALESCE(SUM(quantity * price), 0)  
        FROM sales_item
        WHERE deleted_at IS NULL
    `

	var report model.Report
	err := r.DB.QueryRow(ctx, query).Scan(
		&report.TotalTransactions,
		&report.TotalItemsSold,
		&report.TotalRevenue,
	)

	if err != nil {
		r.Logger.Error("Database Query Error: Failed to generate report",
			zap.Error(err),
		)
		return nil, err
	}

	// Cache in Redis for 3 minutes
	if r.Redis != nil {
		if data, err := json.Marshal(report); err == nil {
			_ = r.Redis.Set(ctx, cacheKey, data, 3*time.Minute).Err()
		}
	}

	return &report, nil
}

func (r *ReportRepo) InvalidateReportCache(ctx context.Context) {
	if r.Redis != nil {
		_ = r.Redis.Del(ctx, "inventory:report:summary").Err()
	}
}
