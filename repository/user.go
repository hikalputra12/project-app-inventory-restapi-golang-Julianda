package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

//untuk mengelola user

//untuk super admin

// buat struct
type userRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}
type UserRepoInterface interface {
	GetAllUser(page, limit int) ([]model.User, int, error)
	FindByEmail(email string) (*model.User, error)
}

// constructor
func NewUserRepo(db database.PgxIface,
	log *zap.Logger) UserRepoInterface {
	return &userRepo{
		DB:     db,
		Logger: log,
	}
}

// untuk membaca user yang ada
func (r *userRepo) GetAllUser(page, limit int) ([]model.User, int, error) {

	//menghitung offset
	offset := (page - 1) * limit
	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(context.Background(), countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query findall repo ", zap.Error(err))
		return nil, 0, err
	}
	query := `SELECT name,email,role 
	FROM users
	WHERE deleted_at IS NULL
	ORDER BY user_id ASC
	LIMIT $1 OFFSET $2`
	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {

		return nil, 0, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var t model.User
		err := rows.Scan(&t.Name, &t.Email, &t.Role)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, t)
	}
	return users, total, nil
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at, name, email, password, role
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	var user model.User
	err := r.DB.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		&user.Name, &user.Email, &user.Password, &user.Role,
	)

	if err == sql.ErrNoRows {
		return nil, nil // menandakan  tidak ditemukan
	}

	return &user, err
}
