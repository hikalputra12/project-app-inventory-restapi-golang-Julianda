package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type RackRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

type RackRepoInterface interface {
	GetAllRack(ctx context.Context, page, limit int) ([]model.Rack, int, error)
	CreateRack(ctx context.Context, rack *model.Rack) error
	UpdateRack(ctx context.Context, id int, rack *model.Rack) error
	DeleteRack(ctx context.Context, id int) error
	GetRackByID(ctx context.Context, id int) (*model.Rack, error)
}

func NewRackRepo(db database.PgxIface, log *zap.Logger) RackRepoInterface {
	return &RackRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *RackRepo) CreateRack(ctx context.Context, rack *model.Rack) error {
	query := `INSERT INTO rack_inventory (name, warehouse_inventory_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4) RETURNING rack_inventory_id`

	now := time.Now()
	err := r.DB.QueryRow(ctx, query, rack.Name, rack.WarehouseInventoryId, now, now).Scan(&rack.ID)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to create rack",
			zap.Error(err),
		)
		return err
	}
	rack.CreatedAt = now
	rack.UpdatedAt = now
	return nil
}

func (r *RackRepo) GetAllRack(ctx context.Context, page, limit int) ([]model.Rack, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM rack_inventory WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT 
		r.rack_inventory_id,
		r.name,
		r.warehouse_inventory_id,
		COALESCE(w.name, '') AS warehouse_name
	FROM rack_inventory r
	LEFT JOIN warehouse_inventory w ON r.warehouse_inventory_id = w.warehouse_inventory_id
	WHERE r.deleted_at IS NULL
	ORDER BY r.rack_inventory_id ASC
	LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to get all racks",
			zap.Error(err),
		)
		return nil, 0, err
	}
	defer rows.Close()

	var racks []model.Rack
	for rows.Next() {
		var t model.Rack
		err := rows.Scan(&t.ID, &t.Name, &t.WarehouseInventoryId, &t.WarehouseInventory)
		if err != nil {
			return nil, 0, err
		}
		racks = append(racks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return racks, total, nil
}

func (r *RackRepo) UpdateRack(ctx context.Context, id int, rack *model.Rack) error {
	query := `UPDATE rack_inventory
			SET name = $1, warehouse_inventory_id = $2, updated_at = $3 
			WHERE rack_inventory_id = $4 AND deleted_at IS NULL`

	now := time.Now()
	cmdTag, err := r.DB.Exec(ctx, query, rack.Name, rack.WarehouseInventoryId, now, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to update rack",
			zap.Error(err),
			zap.Int("rack_id", id),
		)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	rack.UpdatedAt = now
	return nil
}

func (r *RackRepo) DeleteRack(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE rack_inventory SET deleted_at = NOW(), updated_at = NOW() WHERE rack_inventory_id = $1 AND deleted_at IS NULL`

	cmdTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to delete rack",
			zap.Error(err),
			zap.Int("rack_id", id),
		)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *RackRepo) GetRackByID(ctx context.Context, id int) (*model.Rack, error) {
	var rack model.Rack
	query := `SELECT 
		r.rack_inventory_id,
		r.name,
		r.warehouse_inventory_id,
		COALESCE(w.name, '') as warehouse_inventory 
	FROM rack_inventory r
	LEFT JOIN warehouse_inventory w ON r.warehouse_inventory_id = w.warehouse_inventory_id
	WHERE r.rack_inventory_id = $1 AND r.deleted_at IS NULL;`

	err := r.DB.QueryRow(ctx, query, id).Scan(&rack.ID, &rack.Name, &rack.WarehouseInventoryId, &rack.WarehouseInventory)
	if err != nil {
		return nil, err
	}
	return &rack, nil
}
