package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"

	"go.uber.org/zap"
)

type CategoryService struct {
	repo   repository.Repo
	logger *zap.Logger
}
type CategoryServiceInterface interface {
	GetAllCategory(page, limit int) ([]model.Category, *dto.Pagination, error)
	CreateCategory(category *model.Category) error
	UpdateCategory(id int, category *model.Category) error
	DeleteCategory(id int) error
}

// constructor
func NewCategoryService(repo repository.Repo, log *zap.Logger) CategoryServiceInterface {
	return &CategoryService{
		repo:   repo,
		logger: log,
	}
}

func (s *CategoryService) GetAllCategory(page, limit int) ([]model.Category, *dto.Pagination, error) {
	Categories, total, err := s.repo.CategoryRepo.GetAllCategory(page, limit)
	if err != nil {
		s.logger.Error("failed to connect service to read list Category", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
	}
	return Categories, &pagination, nil
}

func (s *CategoryService) CreateCategory(Category *model.Category) error {
	err := s.repo.CategoryRepo.CreateCategory(Category)
	if err != nil {
		return err
	}
	return nil
}

func (s *CategoryService) UpdateCategory(id int, Category *model.Category) error {
	err := s.repo.CategoryRepo.UpdateCategory(id, Category)
	if err != nil {
		return err
	}
	return nil
}
func (s *CategoryService) DeleteCategory(id int) error {
	err := s.repo.CategoryRepo.DeleteCategory(id)
	if err != nil {
		return err
	}
	return nil
}
