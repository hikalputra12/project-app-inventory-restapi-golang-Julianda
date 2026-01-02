package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"

	"go.uber.org/zap"
)

type userService struct {
	repo   repository.Repo
	logger *zap.Logger
}
type UserServiceInterface interface {
	GetAllUser(limit, offset int) ([]model.User, *dto.Pagination, error)
}

// constructor
func NewUserService(repo repository.Repo, log *zap.Logger) UserServiceInterface {
	return &userService{
		repo:   repo,
		logger: log,
	}
}

func (s *userService) GetAllUser(page, limit int) ([]model.User, *dto.Pagination, error) {
	users, total, err := s.repo.UserRepo.GetAllUser(page, limit)
	if err != nil {
		s.logger.Error("failed to connect service to read list user", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
	}
	return users, &pagination, nil
}
