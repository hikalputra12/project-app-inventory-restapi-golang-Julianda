package service

import (
	"app-inventory/model"
	"app-inventory/repository"
	"errors"

	"go.uber.org/zap"
)

type AuthServiceInterface interface {
	Login(email, password string) (*model.User, error)
}

type authService struct {
	Repo   repository.Repo
	logger *zap.Logger
}

func NewAuthService(repo repository.Repo, log *zap.Logger) AuthServiceInterface {
	return &authService{Repo: repo,
		logger: log}
}

// untuk auth login
func (s *authService) Login(email, password string) (*model.User, error) {
	user, err := s.Repo.UserRepo.FindByEmail(email)
	if err != nil {
		s.logger.Warn("Login attempt failed: email not found",
			zap.String("email", email))
		return nil, errors.New("user not found")
	}

	if user.Password != password {
		s.logger.Warn("Login attempt failed: password not found",
			zap.String("password", password),
		)
		return nil, errors.New("incorrect password")
	}

	return user, nil
}
