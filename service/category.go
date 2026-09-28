package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"
	"context"

	"go.uber.org/zap"
)

type CategoryService struct {
	repo   repository.CategoryRepoInterface
	logger *zap.Logger
}

type CategoryServiceInterface interface {
	GetAllCategory(ctx context.Context, page, limit int) ([]model.Category, *dto.Pagination, error)
	CreateCategory(ctx context.Context, category *model.Category) error
	UpdateCategory(ctx context.Context, id int, category *model.Category) error
	DeleteCategory(ctx context.Context, id int) error
	GetCategoryById(ctx context.Context, id int) (*model.Category, error)
}

func NewCategoryService(repo repository.CategoryRepoInterface, log *zap.Logger) CategoryServiceInterface {
	return &CategoryService{
		repo:   repo,
		logger: log,
	}
}

func (s *CategoryService) GetAllCategory(ctx context.Context, page, limit int) ([]model.Category, *dto.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	categories, total, err := s.repo.GetAllCategory(ctx, page, limit)
	if err != nil {
		s.logger.Error("Failed to get all category from repository", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
		TotalItems:  int64(total),
	}
	return categories, &pagination, nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *model.Category) error {
	err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		s.logger.Error("Failed to create category on repository", zap.Error(err))
		return err
	}
	return nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id int, category *model.Category) error {
	existing, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return err
	}

	if category.Name != "" {
		existing.Name = category.Name
	}
	if category.Rack_inventory_id != 0 {
		existing.Rack_inventory_id = category.Rack_inventory_id
	}

	err = s.repo.UpdateCategory(ctx, id, existing)
	if err != nil {
		s.logger.Error("Failed to update category on repository", zap.Error(err))
		return err
	}
	return nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete category on repository", zap.Error(err))
		return err
	}
	return nil
}

func (s *CategoryService) GetCategoryById(ctx context.Context, id int) (*model.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to read category by id", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	return category, nil
}
