package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"time"

	"go.uber.org/zap"
)

type CategoryRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}
type CategoryRepoInterface interface {
	GetAllCategory(page, limit int) ([]model.Category, int, error)
	CreateCategory(category *model.Category) error
	UpdateCategory(id int, category *model.Category) error
	DeleteCategory(id int) error
}

// constructor
func NewCategoryRepo(db database.PgxIface,
	log *zap.Logger) CategoryRepoInterface {
	return &CategoryRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *CategoryRepo) CreateCategory(category *model.Category) error {
	query := `INSERT INTO category_inventory ("name", "rack_inventory_id", created_at, updated_at)
VALUES ($1, $2, $3, $4) RETURNING category_inventory_id`

	now := time.Now()
	err := r.DB.QueryRow(context.Background(), query, category.Name, category.Rack_inventory_id, now, now).Scan(&category.ID)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal Insert Category",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	category.CreatedAt = now
	category.UpdatedAt = now
	return nil
}

// untuk membaca Category yang ada
func (r *CategoryRepo) GetAllCategory(page, limit int) ([]model.Category, int, error) {

	//menghitung offset
	offset := (page - 1) * limit
	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM category_inventory WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	query := `SELECT name,rack_inventory_id FROM category_inventory
	ORDER BY category_inventory_id ASC
LIMIT $1 OFFSET $2;`
	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal select all Category",
			zap.Error(err),
			zap.String("query", query),
		)
		return nil, 0, err
	}
	defer rows.Close()
	var Categories []model.Category
	for rows.Next() {
		var t model.Category
		err := rows.Scan(&t.Name, &t.Rack_inventory_id)
		if err != nil {
			return nil, 0, err
		}
		Categories = append(Categories, t)
	}
	return Categories, total, nil
}

// update category
func (r *CategoryRepo) UpdateCategory(id int, category *model.Category) error {
	query := `UPDATE category_inventory
			SET name=$1,rack_inventory_id=$2,updated_at=$3 WHERE category_inventory_id=$4`
	now := time.Now()
	_, err := r.DB.Exec(context.Background(), query, category.Name, category.Rack_inventory_id, now, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal update Category",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	category.UpdatedAt = now
	return nil
}

// delete category
func (r *CategoryRepo) DeleteCategory(id int) error {
	query := `DELETE FROM category_inventory
			 where category_inventory_id = $1`

	_, err := r.DB.Exec(context.Background(), query, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Gagal delete Category",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	return nil
}
