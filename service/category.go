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
	GetCategoryById(id int) (*model.Category, error)
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
		s.logger.Error("failed get all list category on repository ",
			zap.Error(err),
		)
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
		s.logger.Error("failed created category on repository ",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *CategoryService) UpdateCategory(id int, Category *model.Category) error {
	getCategory, err := s.repo.CategoryRepo.GetCategoryByID(id)
	if err != nil {
		return err
	}

	if Category.Name != "" {
		getCategory.Name = Category.Name
	}
	if Category.Rack_inventory_id != 0 {
		getCategory.Rack_inventory_id = Category.Rack_inventory_id
	}

	NewCategoryUpdate := &model.Category{
		Name:              getCategory.Name,
		Rack_inventory_id: getCategory.Rack_inventory_id,
	}
	err = s.repo.CategoryRepo.UpdateCategory(id, NewCategoryUpdate)
	if err != nil {
		s.logger.Error("failed updated category on repository ",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (s *CategoryService) DeleteCategory(id int) error {
	err := s.repo.CategoryRepo.DeleteCategory(id)
	if err != nil {
		s.logger.Error("failed deleted category on repository ",
			zap.Error(err),
		)
		return err
	}
	return nil
}

// get category by id
func (s *CategoryService) GetCategoryById(id int) (*model.Category, error) {
	category, err := s.repo.CategoryRepo.GetCategoryByID(id)
	if err != nil {
		s.logger.Error("failed to connect service to read category by id", zap.Error(err))
		return nil, err
	}

	return category, nil
}
