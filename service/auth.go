package service

import (
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"
	"errors"

	"go.uber.org/zap"
)

type AuthServiceInterface interface {
	Login(email, password string) (*model.User, error)
}

type authService struct {
	Repo   repository.UserRepoInterface
	logger *zap.Logger
}

func NewAuthService(repo repository.UserRepoInterface, log *zap.Logger) AuthServiceInterface {
	return &authService{Repo: repo,
		logger: log}
}

// untuk auth login
func (s *authService) Login(email, password string) (*model.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		s.logger.Warn("Login attempt failed: email not found",
			zap.String("email", email))
		return nil, errors.New("user not found")
	}

	if !utils.CompareHashAndPassword([]byte(user.Password), []byte(password)) {
		s.logger.Warn("Login attempt failed: password is incorrect",
			zap.String("password", password),
		)
		return nil, errors.New("incorrect password")
	}

	return user, nil
}
