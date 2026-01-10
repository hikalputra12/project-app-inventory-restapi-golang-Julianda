package service

import (
	"app-inventory/model"
	"app-inventory/repository"

	"go.uber.org/zap"
)

type SessionService struct {
	repo   repository.Repo
	logger *zap.Logger
}
type SessionServiceInterface interface {
	CreateSession(session *model.Session) error
	RevokeSession(session *model.Session) error
	ExtendSession(id int, session *model.Session) error
	IsValid(id int, session *model.Session) (bool, error)
}

// constructor
func NewSessionService(repo repository.Repo, log *zap.Logger) SessionServiceInterface {
	return &SessionService{
		repo:   repo,
		logger: log,
	}
}

func (s *SessionService) CreateSession(session *model.Session) error {
	err := s.repo.SessionRepo.CreateSession(session)
	if err != nil {
		s.logger.Error("failed create Session on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *SessionService) RevokeSession(session *model.Session) error {
	err := s.repo.SessionRepo.RevokeSession(session)
	if err != nil {
		s.logger.Error("failed revoke Session on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (s *SessionService) ExtendSession(id int, session *model.Session) error {
	err := s.repo.SessionRepo.ExtendSession(id, session)
	if err != nil {
		s.logger.Error("failed extend Session on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (s *SessionService) IsValid(id int, session *model.Session) (bool, error) {
	valid, err := s.repo.SessionRepo.IsValid(id, session)
	if err != nil {
		s.logger.Error("failed check valid session Session on repository",
			zap.Error(err),
		)
		return false, err
	}
	return valid, err
}
