package repository

import (
	"app-inventory/database"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type PermissionIface interface {
	Allowed(ctx context.Context, userID int, code string) (bool, error)
}

type permissionRepo struct {
	db    database.PgxIface
	redis *redis.Client
}

func NewPermissionRepository(db database.PgxIface, rdb *redis.Client) PermissionIface {
	return &permissionRepo{db: db, redis: rdb}
}

func (r *permissionRepo) Allowed(ctx context.Context, userID int, code string) (bool, error) {
	cacheKey := fmt.Sprintf("perm:%d:%s", userID, code)

	if r.redis != nil {
		val, err := r.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			return val == "1", nil
		}
	}

	const qAllowed = `
    WITH perm AS (
        SELECT id FROM permissions WHERE code = $2
    )
    SELECT
    CASE
        WHEN EXISTS (
            SELECT 1 FROM user_permissions up, perm
            WHERE up.user_id = $1 AND up.permission_id = perm.id AND up.effect='deny'
        ) THEN FALSE
        
        WHEN EXISTS (
            SELECT 1 FROM user_permissions up, perm
            WHERE up.user_id = $1 AND up.permission_id = perm.id AND up.effect='allow'
        ) THEN TRUE
        
        WHEN EXISTS (
            SELECT 1
            FROM users u
            JOIN role_permissions rp ON rp.role_id = u.role_id
            JOIN perm ON perm.id = rp.permission_id
            WHERE u.user_id = $1 AND u.deleted_at IS NULL
        ) THEN TRUE
        
        ELSE FALSE
    END AS allowed;
    `

	var allowed bool
	err := r.db.QueryRow(ctx, qAllowed, userID, code).Scan(&allowed)
	if err != nil {
		return false, err
	}

	if r.redis != nil {
		val := "0"
		if allowed {
			val = "1"
		}
		_ = r.redis.Set(ctx, cacheKey, val, 5*time.Minute).Err()
	}

	return allowed, nil
}
