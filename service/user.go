package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"
	"context"

	"go.uber.org/zap"
)

type userService struct {
	repo   repository.UserRepoInterface
	logger *zap.Logger
}

type UserServiceInterface interface {
	GetAllUser(ctx context.Context, page, limit int) ([]model.User, *dto.Pagination, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, id int, user *model.User) error
	DeleteUser(ctx context.Context, id int) error
	GetUserById(ctx context.Context, id int) (*model.User, error)
}

func NewUserService(repo repository.UserRepoInterface, log *zap.Logger) UserServiceInterface {
	return &userService{
		repo:   repo,
		logger: log,
	}
}

func (s *userService) GetAllUser(ctx context.Context, page, limit int) ([]model.User, *dto.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	users, total, err := s.repo.GetAllUser(ctx, page, limit)
	if err != nil {
		s.logger.Error("Failed to fetch user list from repository", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
		TotalItems:  int64(total),
	}
	return users, &pagination, nil
}

func (s *userService) GetUserById(ctx context.Context, id int) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch user by id from repository", zap.Error(err), zap.Int("user_id", id))
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	passwordHash, err := utils.HashPassword(user.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return err
	}

	newUser := &model.User{
		Name:     user.Name,
		Password: passwordHash,
		Email:    user.Email,
		Role_id:  user.Role_id,
	}
	return s.repo.CreateUser(ctx, newUser)
}

func (s *userService) UpdateUser(ctx context.Context, id int, user *model.User) error {
	existingUser, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	// Fixed: Only update password if a new password was provided!
	if user.Password != "" {
		newHash, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		existingUser.Password = newHash
	}
	if user.Role_id != 0 {
		existingUser.Role_id = user.Role_id
	}

	return s.repo.UpdateUser(ctx, id, existingUser)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}
