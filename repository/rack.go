package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"time"

	"go.uber.org/zap"
)

type RackRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}
type RackRepoInterface interface {
	GetAllRack(page, limit int) ([]model.Rack, int, error)
	CreateRack(Rack *model.Rack) error
	UpdateRack(id int, Rack *model.Rack) error
	DeleteRack(id int) error
}

// constructor
func NewRackRepo(db database.PgxIface,
	log *zap.Logger) RackRepoInterface {
	return &RackRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *RackRepo) CreateRack(Rack *model.Rack) error {
	query := `INSERT INTO Rack_inventory ("name", "warehouse_inventory_id", created_at, updated_at)
VALUES ($1, $2, $3, $4) RETURNING rack_inventory_id`

	now := time.Now()
	err := r.DB.QueryRow(context.Background(), query, Rack.Name, Rack.Warehouse_inventory_id, now, now).Scan(&Rack.ID)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal create rak",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	Rack.CreatedAt = now
	Rack.UpdatedAt = now
	return nil
}

// untuk membaca Rack yang ada
func (r *RackRepo) GetAllRack(page, limit int) ([]model.Rack, int, error) {

	//menghitung offset
	offset := (page - 1) * limit
	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM rack_inventory WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	query := `SELECT name,warehouse_inventory_id FROM rack_inventory
	ORDER BY rack_inventory_id ASC
LIMIT $1 OFFSET $2;`
	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal mendapatkan semua rak",
			zap.Error(err),
			zap.String("query", query),
		)
		return nil, 0, err
	}
	defer rows.Close()
	var rack []model.Rack
	for rows.Next() {
		var t model.Rack
		err := rows.Scan(&t.Name, &t.Warehouse_inventory_id)
		if err != nil {
			return nil, 0, err
		}
		rack = append(rack, t)
	}
	return rack, total, nil
}

// update Rack
func (r *RackRepo) UpdateRack(id int, Rack *model.Rack) error {
	query := `UPDATE rack_inventory
			SET name=$1,warehouse_inventory_id=$2,updated_at=$3 WHERE rack_inventory_id=$4`
	now := time.Now()
	_, err := r.DB.Exec(context.Background(), query, Rack.Name, Rack.Warehouse_inventory_id, now, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal update rak",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	Rack.UpdatedAt = now
	return nil
}

// delete Rack
func (r *RackRepo) DeleteRack(id int) error {
	query := `DELETE FROM rack_inventory
			 where rack_inventory_id = $1`

	_, err := r.DB.Exec(context.Background(), query, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal delete rak",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	return nil
}
