package model

import "time"

type Session struct {
	SessionID  string     `json:"session_id"`
	UserID     int        `json:"user_id"`
	ExpiredAt  time.Time  `json:"expired_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	LastActive time.Time  `json:"last_active"`
	CreatedAt  time.Time  `json:"created_at"`
}
