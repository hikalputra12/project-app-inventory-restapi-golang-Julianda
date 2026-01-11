package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"time"

	"go.uber.org/zap"
)

type WarehouseRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}
type WarehouseRepoInterface interface {
	GetAllWarehouse(page, limit int) ([]model.Warehouse, int, error)
	CreateWarehouse(Warehouse *model.Warehouse) error
	UpdateWarehouse(id int, Warehouse *model.Warehouse) error
	DeleteWarehouse(id int) error
	GetWarehouseByID(id int) (*model.Warehouse, error)
}

// constructor
func NewWarehouseRepo(db database.PgxIface,
	log *zap.Logger) WarehouseRepoInterface {
	return &WarehouseRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *WarehouseRepo) CreateWarehouse(Warehouse *model.Warehouse) error {
	query := `INSERT INTO warehouse_inventory ("name", "location", created_at, updated_at)
VALUES ($1, $2, $3, $4) RETURNING warehouse_inventory_id`

	now := time.Now()
	err := r.DB.QueryRow(context.Background(), query, Warehouse.Name, Warehouse.Location, now, now).Scan(&Warehouse.ID)
	if err != nil {
		r.Logger.Error("failed to insert a new warehouse to database",
			zap.Error(err),
			zap.String("warehouse_name: %s", Warehouse.Name),
		)
		return err
	}
	Warehouse.CreatedAt = now
	Warehouse.UpdatedAt = now
	return nil
}

// untuk membaca Warehouse yang ada
func (r *WarehouseRepo) GetAllWarehouse(page, limit int) ([]model.Warehouse, int, error) {

	//menghitung offset
	offset := (page - 1) * limit
	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM warehouse_inventory WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	query := `SELECT name FROM warehouse_inventory
	ORDER BY warehouse_inventory_id ASC
LIMIT $1 OFFSET $2;`
	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal mendapatkan data jenis gudang",
			zap.Error(err),
			zap.String("query", query),
		)
		return nil, 0, err
	}
	defer rows.Close()
	var Warehouse []model.Warehouse
	for rows.Next() {
		var t model.Warehouse
		err := rows.Scan(&t.Name)
		if err != nil {
			return nil, 0, err
		}
		Warehouse = append(Warehouse, t)
	}
	return Warehouse, total, nil
}

// update Warehouse
func (r *WarehouseRepo) UpdateWarehouse(id int, Warehouse *model.Warehouse) error {
	query := `UPDATE Warehouse_inventory
			SET name=$1,location=$2,updated_at=$3 WHERE warehouse_inventory_id=$4`
	now := time.Now()
	_, err := r.DB.Exec(context.Background(), query, Warehouse.Name, Warehouse.Location, now, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal mengubah jenis gudang",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	Warehouse.UpdatedAt = now
	return nil
}

// delete Warehouse
func (r *WarehouseRepo) DeleteWarehouse(id int) error {
	query := `DELETE FROM warehouse_inventory
			 where warehouse_inventory_id = $1`

	_, err := r.DB.Exec(context.Background(), query, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal menghapus gudang",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	return nil
}

// untuk membaca warehouse berdarsaskan id
func (r *WarehouseRepo) GetWarehouseByID(id int) (*model.Warehouse, error) {
	var warehouse model.Warehouse
	query := `SELECT name, location FROM warehouse_inventory 
WHERE warehouse_inventory_id = $1;`
	err := r.DB.QueryRow(context.Background(), query, id).Scan(&warehouse.Name, &warehouse.Location)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal mendapatkan rak berdarsaskan id",
			zap.Error(err),
			zap.String("query", query),
		)
		return nil, err
	}
	return &warehouse, nil
}
