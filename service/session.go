package service

import (
	"app-inventory/model"
	"app-inventory/repository"

	"go.uber.org/zap"
)

type SessionService struct {
	repo   repository.SessionRepoInterface
	logger *zap.Logger
}
type SessionServiceInterface interface {
	CreateSession(session *model.Session) error
	RevokeSession(session *model.Session) error
	ExtendSession(session *model.Session) error
	IsValid(session *model.Session) (bool, error)
	GetUserIDBySession(session *model.Session) (int, error)
}

// constructor
func NewSessionService(repo repository.SessionRepoInterface, log *zap.Logger) SessionServiceInterface {
	return &SessionService{
		repo:   repo,
		logger: log,
	}
}

func (s *SessionService) CreateSession(session *model.Session) error {
	err := s.repo.CreateSession(session)
	if err != nil {
		s.logger.Error("failed create Session on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *SessionService) RevokeSession(session *model.Session) error {
	err := s.repo.RevokeSession(session)
	if err != nil {
		s.logger.Error("failed revoke Session on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (s *SessionService) ExtendSession(session *model.Session) error {
	err := s.repo.ExtendSession(session)
	if err != nil {
		s.logger.Error("failed extend Session on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (s *SessionService) IsValid(session *model.Session) (bool, error) {
	valid, err := s.repo.IsValid(session)
	if err != nil {
		s.logger.Error("failed check valid session Session on repository",
			zap.Error(err),
		)
		return false, err
	}
	return valid, err
}
func (s *SessionService) GetUserIDBySession(session *model.Session) (int, error) {
	userID, err := s.repo.GetUserIDBySession(session)
	if err != nil {
		s.logger.Error("failed extend Session on repository",
			zap.Error(err),
		)
		return 0, err
	}
	return userID, nil
}
