package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type InventoryRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

type InventoryRepoInterface interface {
	GetAllInventory(ctx context.Context, page, limit int) ([]model.Inventory, int, error)
	CreateInventory(ctx context.Context, inventory *model.Inventory) error
	UpdateInventory(ctx context.Context, id int, inventory *model.Inventory) error
	DeleteInventory(ctx context.Context, id int) error
	CheckStock(ctx context.Context, page, limit int) ([]model.Inventory, int, error)
	GetInventoryByID(ctx context.Context, id int) (*model.Inventory, error)
}

func NewInventoryRepo(db database.PgxIface, log *zap.Logger) InventoryRepoInterface {
	return &InventoryRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *InventoryRepo) CreateInventory(ctx context.Context, inventory *model.Inventory) error {
	query := `INSERT INTO inventories (name, price, stock, category_inventory_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING inventory_id`

	now := time.Now()
	err := r.DB.QueryRow(ctx, query, inventory.Name, inventory.Price, inventory.Stock, inventory.Category_inventory_id, now, now).Scan(&inventory.ID)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to create inventory",
			zap.Error(err),
		)
		return err
	}
	inventory.CreatedAt = now
	inventory.UpdatedAt = now
	return nil
}

func (r *InventoryRepo) GetAllInventory(ctx context.Context, page, limit int) ([]model.Inventory, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM inventories WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT 
		i.inventory_id,
		i.name AS product_name, 
		i.price, 
		i.stock, 
		c.name AS category_name, 
		r.name AS rack_name,     
		w.name AS warehouse_name,
		i.category_inventory_id
	FROM inventories i
	JOIN category_inventory c ON i.category_inventory_id = c.category_inventory_id
	JOIN rack_inventory r ON c.rack_inventory_id = r.rack_inventory_id
	JOIN warehouse_inventory w ON r.warehouse_inventory_id = w.warehouse_inventory_id 
	WHERE i.deleted_at IS NULL
	ORDER BY i.inventory_id ASC
	LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to get all inventory",
			zap.Error(err),
		)
		return nil, 0, err
	}
	defer rows.Close()

	var inventories []model.Inventory
	for rows.Next() {
		var t model.Inventory
		err := rows.Scan(&t.ID, &t.Name, &t.Price, &t.Stock, &t.Category, &t.Rack, &t.Warehouse, &t.Category_inventory_id)
		if err != nil {
			return nil, 0, err
		}
		inventories = append(inventories, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}

func (r *InventoryRepo) UpdateInventory(ctx context.Context, id int, inventory *model.Inventory) error {
	query := `UPDATE inventories
			SET name = $1, price = $2, stock = $3, category_inventory_id = $4, updated_at = $5 
			WHERE inventory_id = $6 AND deleted_at IS NULL`

	now := time.Now()
	cmdTag, err := r.DB.Exec(ctx, query, inventory.Name, inventory.Price, inventory.Stock, inventory.Category_inventory_id, now, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to update inventory",
			zap.Error(err),
			zap.Int("inventory_id", id),
		)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	inventory.UpdatedAt = now
	return nil
}

func (r *InventoryRepo) DeleteInventory(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE inventories SET deleted_at = NOW(), updated_at = NOW() WHERE inventory_id = $1 AND deleted_at IS NULL`

	cmdTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to delete inventory",
			zap.Error(err),
			zap.Int("inventory_id", id),
		)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *InventoryRepo) CheckStock(ctx context.Context, page, limit int) ([]model.Inventory, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int
	// Fixed: Count only items where stock <= 5
	countQuery := `SELECT COUNT(*) FROM inventories WHERE deleted_at IS NULL AND stock <= 5`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT 
		i.inventory_id,
		i.name AS product_name, 
		i.price, 
		i.stock, 
		c.name AS category_name, 
		r.name AS rack_name,     
		w.name AS warehouse_name,
		i.category_inventory_id
	FROM inventories i
	JOIN category_inventory c ON i.category_inventory_id = c.category_inventory_id
	JOIN rack_inventory r ON c.rack_inventory_id = r.rack_inventory_id
	JOIN warehouse_inventory w ON r.warehouse_inventory_id = w.warehouse_inventory_id 
	WHERE i.deleted_at IS NULL AND i.stock <= 5 
	ORDER BY i.inventory_id ASC
	LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		r.Logger.Error("Failed to check low stock inventory", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	var inventories []model.Inventory
	for rows.Next() {
		var t model.Inventory
		err := rows.Scan(&t.ID, &t.Name, &t.Price, &t.Stock, &t.Category, &t.Rack, &t.Warehouse, &t.Category_inventory_id)
		if err != nil {
			return nil, 0, err
		}
		inventories = append(inventories, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}

func (r *InventoryRepo) GetInventoryByID(ctx context.Context, id int) (*model.Inventory, error) {
	var inventory model.Inventory
	query := `SELECT 
		i.inventory_id,
		i.name AS product_name, 
		i.price, 
		i.stock,
		i.category_inventory_id, 
		c.name AS category_name, 
		r.name AS rack_name,     
		w.name AS warehouse_name
	FROM inventories i
	JOIN category_inventory c ON i.category_inventory_id = c.category_inventory_id
	JOIN rack_inventory r ON c.rack_inventory_id = r.rack_inventory_id
	JOIN warehouse_inventory w ON r.warehouse_inventory_id = w.warehouse_inventory_id 
	WHERE i.inventory_id = $1 AND i.deleted_at IS NULL`

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&inventory.ID,
		&inventory.Name,
		&inventory.Price,
		&inventory.Stock,
		&inventory.Category_inventory_id,
		&inventory.Category,
		&inventory.Rack,
		&inventory.Warehouse,
	)
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}
