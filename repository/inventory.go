package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"time"

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
	CreateInventory(inventory *model.Inventory) error
	UpdateInventory(id int, inventory *model.Inventory) error
	DeleteInventory(id int) error
}

// constructor
func NewInventoryRepo(db database.PgxIface,
	log *zap.Logger) InventoryRepoInterface {
	return &InventoryRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *InventoryRepo) CreateInventory(inventory *model.Inventory) error {
	query := `INSERT INTO inventories ("name", "price", "stock", "category_inventory_id", created_at, updated_at)
VALUES
($1, $2, $3, $4,$5,$6) RETURNING inventory_id`

	now := time.Now()
	err := r.DB.QueryRow(context.Background(), query, inventory.Name, inventory.Price, inventory.Stock, inventory.Category_inventory_id, now, now).Scan(&inventory.ID)
	if err != nil {
		return err
	}
	inventory.CreatedAt = now
	inventory.UpdatedAt = now
	return nil
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

// update inventory
func (r *InventoryRepo) UpdateInventory(id int, inventory *model.Inventory) error {
	query := `UPDATE inventories
			SET name=$1,price=$2,stock=$3,category_inventory_id = $4, updated_at=$5 where inventory_id=$6`

	now := time.Now()
	_, err := r.DB.Exec(context.Background(), query, inventory.Name, inventory.Price, inventory.Stock, inventory.Category_inventory_id, now, id)
	if err != nil {
		return err
	}
	inventory.UpdatedAt = now
	return nil
}

// delete inventory
func (r *InventoryRepo) DeleteInventory(id int) error {
	query := `DELETE FROM inventories
			 where inventory_id = $1`

	_, err := r.DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
