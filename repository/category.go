package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type CategoryRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

type CategoryRepoInterface interface {
	GetAllCategory(ctx context.Context, page, limit int) ([]model.Category, int, error)
	CreateCategory(ctx context.Context, category *model.Category) error
	UpdateCategory(ctx context.Context, id int, category *model.Category) error
	DeleteCategory(ctx context.Context, id int) error
	GetCategoryByID(ctx context.Context, id int) (*model.Category, error)
}

func NewCategoryRepo(db database.PgxIface, log *zap.Logger) CategoryRepoInterface {
	return &CategoryRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *CategoryRepo) CreateCategory(ctx context.Context, category *model.Category) error {
	query := `INSERT INTO category_inventory (name, rack_inventory_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4) RETURNING category_inventory_id`

	now := time.Now()
	err := r.DB.QueryRow(ctx, query, category.Name, category.Rack_inventory_id, now, now).Scan(&category.ID)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to insert category",
			zap.Error(err),
		)
		return err
	}
	category.CreatedAt = now
	category.UpdatedAt = now
	return nil
}

func (r *CategoryRepo) GetAllCategory(ctx context.Context, page, limit int) ([]model.Category, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM category_inventory WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT 
		c.category_inventory_id,
		c.name,
		c.rack_inventory_id,
		COALESCE(r.name, '') AS rack_name
	FROM category_inventory c
	LEFT JOIN rack_inventory r ON c.rack_inventory_id = r.rack_inventory_id
	WHERE c.deleted_at IS NULL
	ORDER BY c.category_inventory_id ASC
	LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to select all categories",
			zap.Error(err),
		)
		return nil, 0, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var t model.Category
		err := rows.Scan(&t.ID, &t.Name, &t.Rack_inventory_id, &t.RackInventory)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *CategoryRepo) UpdateCategory(ctx context.Context, id int, category *model.Category) error {
	query := `UPDATE category_inventory
			SET name = $1, rack_inventory_id = $2, updated_at = $3 
			WHERE category_inventory_id = $4 AND deleted_at IS NULL`

	now := time.Now()
	cmdTag, err := r.DB.Exec(ctx, query, category.Name, category.Rack_inventory_id, now, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to update category",
			zap.Error(err),
			zap.Int("category_id", id),
		)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	category.UpdatedAt = now
	return nil
}

func (r *CategoryRepo) DeleteCategory(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE category_inventory SET deleted_at = NOW(), updated_at = NOW() WHERE category_inventory_id = $1 AND deleted_at IS NULL`

	cmdTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to soft delete category",
			zap.Error(err),
			zap.Int("category_id", id),
		)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *CategoryRepo) GetCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	var category model.Category
	query := `SELECT 
		c.category_inventory_id,
		c.name,
		c.rack_inventory_id,
		COALESCE(r.name, '') as rack_inventory 
	FROM category_inventory c
	LEFT JOIN rack_inventory r ON c.rack_inventory_id = r.rack_inventory_id
	WHERE c.category_inventory_id = $1 AND c.deleted_at IS NULL;`

	err := r.DB.QueryRow(ctx, query, id).Scan(&category.ID, &category.Name, &category.Rack_inventory_id, &category.RackInventory)
	if err != nil {
		return nil, err
	}
	return &category, nil
}
