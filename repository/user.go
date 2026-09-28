package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type userRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

type UserRepoInterface interface {
	GetAllUser(ctx context.Context, page, limit int) ([]model.User, int, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, id int, user *model.User) error
	DeleteUser(ctx context.Context, id int) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
}

func NewUserRepo(db database.PgxIface, log *zap.Logger) UserRepoInterface {
	return &userRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *userRepo) UpdateUser(ctx context.Context, id int, user *model.User) error {
	query := `UPDATE users
			SET name = $1, email = $2, password_hash = $3, role_id = $4, updated_at = $5 
			WHERE user_id = $6 AND deleted_at IS NULL`

	now := time.Now()
	_, err := r.DB.Exec(ctx, query, user.Name, user.Email, user.Password, user.Role_id, now, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to update user",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		return err
	}
	user.UpdatedAt = now
	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id int) error {
	// Soft delete
	query := `UPDATE users SET deleted_at = NOW(), updated_at = NOW() WHERE user_id = $1 AND deleted_at IS NULL`

	cmdTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to soft delete user",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (name, email, password_hash, role_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id`

	now := time.Now()
	err := r.DB.QueryRow(ctx, query, user.Name, user.Email, user.Password, user.Role_id, now, now).Scan(&user.ID)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to create user",
			zap.Error(err),
		)
		return err
	}
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *userRepo) GetAllUser(ctx context.Context, page, limit int) ([]model.User, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("Failed to count users", zap.Error(err))
		return nil, 0, err
	}

	query := `SELECT 
		users.user_id,
		users.name, 
		users.email, 
		roles.name AS role_name,
		users.role_id
	FROM users
	JOIN roles ON users.role_id = roles.id
	WHERE users.deleted_at IS NULL
	ORDER BY users.user_id ASC
	LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to get user list",
			zap.Error(err),
		)
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var t model.User
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Role, &t.Role_id)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT u.user_id, u.created_at, u.updated_at, u.deleted_at, u.name, u.email, u.password_hash, r.name as role, u.role_id
        FROM users u
        JOIN roles r ON u.role_id = r.id
        WHERE u.email = $1 AND u.deleted_at IS NULL
	`
	var user model.User
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		&user.Name, &user.Email, &user.Password, &user.Role, &user.Role_id,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	query := `SELECT 
		users.user_id,
		users.name, 
		users.email, 
		roles.name AS role_name,
		users.role_id,
		users.password_hash
	FROM users
	JOIN roles ON users.role_id = roles.id
	WHERE users.user_id = $1 AND users.deleted_at IS NULL;`

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.Role_id,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
