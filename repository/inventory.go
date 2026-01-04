package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"

	"go.uber.org/zap"
)

//untuk mengelola Inventory

//untuk super admin

// buat struct
type InventoryRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}
type InventoryRepoInterface interface {
	GetAllInventory(page, limit int) ([]model.Inventory, int, error)
}

// constructor
func NewInventoryRepo(db database.PgxIface,
	log *zap.Logger) InventoryRepoInterface {
	return &InventoryRepo{
		DB:     db,
		Logger: log,
	}
}

// untuk membaca Inventory yang ada
func (r *InventoryRepo) GetAllInventory(page, limit int) ([]model.Inventory, int, error) {

	//menghitung offset
	offset := (page - 1) * limit
	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM inventories WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	query := `SELECT 
    i.name AS product_name, 
    i.price, 
    i.stock, 
    c.name AS category_name, 
    r.name AS rack_name,     
    w.name AS warehouse_name
FROM 
    inventories i
JOIN 
    category_inventory c ON i.category_inventory_id = c.category_inventory_id
JOIN 
    rack_inventory r ON c.rack_inventory_id = r.rack_inventory_id
JOIN 
    warehouse_inventory w ON r.warehouse_inventory_id = w.warehouse_inventory_id 
ORDER BY 
    i.inventory_id ASC
LIMIT $1 OFFSET $2;`
	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {

		return nil, 0, err
	}
	defer rows.Close()
	var Inventories []model.Inventory
	for rows.Next() {
		var t model.Inventory
		err := rows.Scan(&t.Name, &t.Price, &t.Stock, &t.Category, &t.Rack, &t.Warehouse)
		if err != nil {
			return nil, 0, err
		}
		Inventories = append(Inventories, t)
	}
	return Inventories, total, nil
}
