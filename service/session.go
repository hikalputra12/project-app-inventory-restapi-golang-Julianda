package service

import (
	"app-inventory/model"
	"app-inventory/repository"
	"context"

	"go.uber.org/zap"
)

type SessionService struct {
	repo   repository.SessionRepoInterface
	logger *zap.Logger
}

type SessionServiceInterface interface {
	CreateSession(ctx context.Context, session *model.Session) error
	RevokeSession(ctx context.Context, sessionID string) error
	ExtendSession(ctx context.Context, sessionID string) error
	IsValid(ctx context.Context, sessionID string) (bool, error)
	GetUserIDBySession(ctx context.Context, sessionID string) (int, error)
}

func NewSessionService(repo repository.SessionRepoInterface, log *zap.Logger) SessionServiceInterface {
	return &SessionService{
		repo:   repo,
		logger: log,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, session *model.Session) error {
	return s.repo.CreateSession(ctx, session)
}

func (s *SessionService) RevokeSession(ctx context.Context, sessionID string) error {
	return s.repo.RevokeSession(ctx, sessionID)
}

func (s *SessionService) ExtendSession(ctx context.Context, sessionID string) error {
	return s.repo.ExtendSession(ctx, sessionID)
}

func (s *SessionService) IsValid(ctx context.Context, sessionID string) (bool, error) {
	return s.repo.IsValid(ctx, sessionID)
}

func (s *SessionService) GetUserIDBySession(ctx context.Context, sessionID string) (int, error) {
	return s.repo.GetUserIDBySession(ctx, sessionID)
}
