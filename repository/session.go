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
	CreateSession(session *model.Session) error
	RevokeSession(session *model.Session) error
	ExtendSession(id int, session *model.Session) error
	IsValid(id int, session *model.Session) (bool, error)
	GetUserIDBySession(session *model.Session) (int, error)
}

// constructor
func NewSessionRepo(db database.PgxIface,
	log *zap.Logger) SessionRepoInterface {
	return &SessionRepo{
		DB:     db,
		Logger: log,
	}
}

func (r *SessionRepo) CreateSession(session *model.Session) error {
	query := `INSERT INTO sessions ("session_id", "user_id", expired_at, created_at)
VALUES ($1, $2, $3, $4)`

	now := time.Now()
	expired := time.Now().Add(24 * time.Hour)
	_, err := r.DB.Exec(context.Background(), query, session.SessionID, session.UserID, expired, now)
	if err != nil {
		r.Logger.Error("Database Query Error: failed insert session uuid to database",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	session.ExpiredAt = expired
	session.CreatedAt = now
	return nil
}

// revoke Session ( pencabutan session oleh admin) dan pelajari lagi konsepnya
func (r *SessionRepo) RevokeSession(session *model.Session) error {
	query := `UPDATE sessions
			  SET revoked_at=NOW()
			  WHERE session_id=$1 AND revoked_at is NULL`
	_, err := r.DB.Exec(context.Background(), query, session.SessionID)
	if err != nil {
		r.Logger.Error("Database Query Error: failed revoke session on database",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	return nil
}
func (r *SessionRepo) ExtendSession(id int, session *model.Session) error {
	query := `UPDATE sessions 
			SET session_id=$1 WHERE user_id=$2 AND revoked_at is NULL `
	_, err := r.DB.Exec(context.Background(), query, session.SessionID, session.UserID)
	if err != nil {
		r.Logger.Error("Database Query Error: failed update session on database",
			zap.Error(err),
			zap.String("query", query),
		)
		return err
	}
	return nil
}

// pengecekan valid atau tidak validnya session
func (r *SessionRepo) IsValid(id int, session *model.Session) (bool, error) {
	query := `SELECT EXIST(
			  SELECT 1 FROM sessions WHERE session_id=$1 AND revoked_at is NULL AND expired_at > NOW() )`
	var valid bool
	err := r.DB.QueryRow(context.Background(), query, session.SessionID).Scan(&valid)
	return valid, err
}

func (r *SessionRepo) GetUserIDBySession(session *model.Session) (int, error) {
	var userID int
	query := `SELECT user_id FROM sessions WHERE session_id = $1 AND revoked_at IS NULL`

	err := r.DB.QueryRow(context.Background(), query, session.SessionID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
