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
	GetAllUser(page, limit int) ([]model.User, *dto.Pagination, error)
	CreateUser(*model.User) error
	UpdateUser(id int, user *model.User) error
	DeleteUser(id int) error
	GetUserById(id int) (*model.User, error)
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

// get user by id
func (s *userService) GetUserById(id int) (*model.User, error) {
	users, err := s.repo.UserRepo.GetUserByID(id)
	if err != nil {
		s.logger.Error("failed to connect service to read user by id", zap.Error(err))
		return nil, err
	}

	return users, nil
}

func (s *userService) CreateUser(user *model.User) error {
	passwordHash := utils.HashPassword(user.Password)

	NewUser := &model.User{
		Name:     user.Name,
		Password: passwordHash,
		Email:    user.Email,
		Role_id:  user.Role_id,
	}
	err := s.repo.UserRepo.CreateUser(NewUser)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUser(id int, user *model.User) error {
	passwordHash := utils.HashPassword(user.Password)
	getUser, err := s.repo.UserRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	if user.Name != "" {
		getUser.Name = user.Name
	}
	if user.Email != "" {
		getUser.Email = user.Email
	}
	if user.Password != "" {
		getUser.Password = passwordHash
	}
	if user.Role_id != 0 {
		getUser.Role_id = user.Role_id
	}

	NewUserUpdate := &model.User{
		Name:     getUser.Name,
		Password: passwordHash,
		Email:    getUser.Email,
		Role_id:  getUser.Role_id,
	}
	err = s.repo.UserRepo.UpdateUser(id, NewUserUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) DeleteUser(id int) error {
	err := s.repo.UserRepo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
