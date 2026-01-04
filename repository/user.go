package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
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
	CreateUser(*model.User) error
	UpdateUser(id int, user *model.User) error
	DeleteUser(id int) error
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

// update user
func (r *userRepo) UpdateUser(id int, user *model.User) error {
	query := `UPDATE users
			SET name=$1,email=$2,password_hash=$3,role_id=$4,updated_at=$5 where user_id = $6`

	now := time.Now()
	_, err := r.DB.Exec(context.Background(), query, user.Name, user.Email, user.Password, user.Role_id, now, id)
	if err != nil {
		return err
	}
	user.UpdatedAt = now
	return nil
}

// delete user
func (r *userRepo) DeleteUser(id int) error {
	query := `DELETE FROM users
			 where user_id = $1`

	_, err := r.DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}

// create user
func (r *userRepo) CreateUser(user *model.User) error {
	query := `INSERT INTO "users" ("name", "email", "password_hash", "role_id", created_at, updated_at)
VALUES
($1, $2, $3, $4,$5,$6) RETURNING user_id`

	now := time.Now()
	err := r.DB.QueryRow(context.Background(), query, user.Name, user.Email, user.Password, user.Role_id, now, now).Scan(&user.ID)
	if err != nil {
		return err
	}
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
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
	query := `SELECT 
	users.user_id,
    users.name, 
    users.email, 
    roles.name AS role_name
FROM users
JOIN roles ON users.role_id = roles.id
WHERE users.deleted_at IS NULL
ORDER BY users.user_id ASC
LIMIT $1 OFFSET $2;`
	rows, err := r.DB.Query(context.Background(), query, limit, offset)
	if err != nil {

		return nil, 0, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var t model.User
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Role)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, t)
	}
	return users, total, nil
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	query := `
		SELECT u.user_id, u.created_at, u.updated_at, u.deleted_at, u.name, u.email, u.password_hash, r.name as role
        FROM users u
        JOIN roles r ON u.role_id = r.id
        WHERE u.email = $1 AND u.deleted_at IS NULL
	`
	var user model.User
	err := r.DB.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		&user.Name, &user.Email, &user.Password, &user.Role,
	)

	if err == pgx.ErrNoRows {
		return nil, err // menandakan  tidak ditemukan
	}

	return &user, err
}
