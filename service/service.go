package service

import (
	"app-inventory/repository"

	"go.uber.org/zap"
)

type Service struct {
	UserService UserServiceInterface
	AuthService AuthServiceInterface
}

func AllService(repo repository.Repo, log *zap.Logger) Service {
	return Service{
		UserService: NewUserService(repo, log),
		AuthService: NewAuthService(repo, log),
	}
}
