package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type WarehouseRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

type WarehouseRepoInterface interface {
	GetAllWarehouse(ctx context.Context, page, limit int) ([]model.Warehouse, int, error)
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	UpdateWarehouse(ctx context.Context, id int, warehouse *model.Warehouse) error
	DeleteWarehouse(ctx context.Context, id int) error
	GetWarehouseByID(ctx context.Context, id int) (*model.Warehouse, error)
}

func NewWarehouseRepo(db database.PgxIface, log *zap.Logger) WarehouseRepoInterface {
	return &WarehouseRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *WarehouseRepo) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	query := `INSERT INTO warehouse_inventory (name, location, created_at, updated_at)
			VALUES ($1, $2, $3, $4) RETURNING warehouse_inventory_id`

	now := time.Now()
	err := r.DB.QueryRow(ctx, query, warehouse.Name, warehouse.Location, now, now).Scan(&warehouse.ID)
	if err != nil {
		r.Logger.Error("Failed to insert warehouse",
			zap.Error(err),
			zap.String("name", warehouse.Name),
		)
		return err
	}
	warehouse.CreatedAt = now
	warehouse.UpdatedAt = now
	return nil
}

func (r *WarehouseRepo) GetAllWarehouse(ctx context.Context, page, limit int) ([]model.Warehouse, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM warehouse_inventory WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT warehouse_inventory_id, name, location 
		FROM warehouse_inventory
		WHERE deleted_at IS NULL
		ORDER BY warehouse_inventory_id ASC
		LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to get warehouse list",
			zap.Error(err),
		)
		return nil, 0, err
	}
	defer rows.Close()

	var warehouses []model.Warehouse
	for rows.Next() {
		var t model.Warehouse
		err := rows.Scan(&t.ID, &t.Name, &t.Location)
		if err != nil {
			return nil, 0, err
		}
		warehouses = append(warehouses, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return warehouses, total, nil
}

func (r *WarehouseRepo) UpdateWarehouse(ctx context.Context, id int, warehouse *model.Warehouse) error {
	query := `UPDATE warehouse_inventory
			SET name = $1, location = $2, updated_at = $3 
			WHERE warehouse_inventory_id = $4 AND deleted_at IS NULL`

	now := time.Now()
	cmdTag, err := r.DB.Exec(ctx, query, warehouse.Name, warehouse.Location, now, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to update warehouse",
			zap.Error(err),
			zap.Int("warehouse_id", id),
		)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	warehouse.UpdatedAt = now
	return nil
}

func (r *WarehouseRepo) DeleteWarehouse(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE warehouse_inventory SET deleted_at = NOW(), updated_at = NOW() WHERE warehouse_inventory_id = $1 AND deleted_at IS NULL`

	cmdTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to delete warehouse",
			zap.Error(err),
			zap.Int("warehouse_id", id),
		)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *WarehouseRepo) GetWarehouseByID(ctx context.Context, id int) (*model.Warehouse, error) {
	var warehouse model.Warehouse
	query := `SELECT warehouse_inventory_id, name, location 
		FROM warehouse_inventory 
		WHERE warehouse_inventory_id = $1 AND deleted_at IS NULL;`

	err := r.DB.QueryRow(ctx, query, id).Scan(&warehouse.ID, &warehouse.Name, &warehouse.Location)
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}
