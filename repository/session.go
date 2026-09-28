package repository

import (
	"app-inventory/database"
	"app-inventory/model"
	"context"
	"time"

	"go.uber.org/zap"
)

type SessionRepo struct {
	DB     database.PgxIface
	Logger *zap.Logger
}

type SessionRepoInterface interface {
	CreateSession(ctx context.Context, session *model.Session) error
	RevokeSession(ctx context.Context, sessionID string) error
	ExtendSession(ctx context.Context, sessionID string) error
	IsValid(ctx context.Context, sessionID string) (bool, error)
	GetUserIDBySession(ctx context.Context, sessionID string) (int, error)
}

func NewSessionRepo(db database.PgxIface, log *zap.Logger) SessionRepoInterface {
	return &SessionRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *SessionRepo) CreateSession(ctx context.Context, session *model.Session) error {
	query := `INSERT INTO sessions (session_id, user_id, expired_at, created_at, last_active)
		VALUES ($1, $2, $3, $4, $5)`

	now := time.Now()
	expired := now.Add(24 * time.Hour)
	_, err := r.DB.Exec(ctx, query, session.SessionID, session.UserID, expired, now, now)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to insert session",
			zap.Error(err),
		)
		return err
	}
	session.ExpiredAt = expired
	session.CreatedAt = now
	session.LastActive = now
	return nil
}

func (r *SessionRepo) RevokeSession(ctx context.Context, sessionID string) error {
	query := `UPDATE sessions SET revoked_at = NOW() WHERE session_id = $1 AND revoked_at IS NULL`
	_, err := r.DB.Exec(ctx, query, sessionID)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to revoke session",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (r *SessionRepo) ExtendSession(ctx context.Context, sessionID string) error {
	query := `UPDATE sessions 
			SET expired_at = $1, last_active = NOW() 
			WHERE session_id = $2 AND revoked_at IS NULL AND expired_at > NOW()`
	expired := time.Now().Add(24 * time.Hour)
	_, err := r.DB.Exec(ctx, query, expired, sessionID)
	if err != nil {
		r.Logger.Error("Database Query Error: Failed to extend session",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (r *SessionRepo) IsValid(ctx context.Context, sessionID string) (bool, error) {
	query := `SELECT EXISTS(
			  SELECT 1 FROM sessions WHERE session_id = $1 AND revoked_at IS NULL AND expired_at > NOW())`
	var valid bool
	err := r.DB.QueryRow(ctx, query, sessionID).Scan(&valid)
	return valid, err
}

func (r *SessionRepo) GetUserIDBySession(ctx context.Context, sessionID string) (int, error) {
	var userID int
	query := `SELECT user_id FROM sessions WHERE session_id = $1 AND revoked_at IS NULL AND expired_at > NOW()`
	err := r.DB.QueryRow(ctx, query, sessionID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
