package model

import "time"

type Session struct {
	SessionID string
	UserID    int
	ExpiredAt time.Time
	RevokedAt time.Time
	CreatedAt time.Time
}
