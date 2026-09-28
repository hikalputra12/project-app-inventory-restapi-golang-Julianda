package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"
	"context"

	"go.uber.org/zap"
)

type WarehouseService struct {
	repo   repository.WarehouseRepoInterface
	logger *zap.Logger
}

type WarehouseServiceInterface interface {
	GetAllWarehouse(ctx context.Context, page, limit int) ([]model.Warehouse, *dto.Pagination, error)
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	UpdateWarehouse(ctx context.Context, id int, warehouse *model.Warehouse) error
	DeleteWarehouse(ctx context.Context, id int) error
	GetWarehouseById(ctx context.Context, id int) (*model.Warehouse, error)
}

func NewWarehouseService(repo repository.WarehouseRepoInterface, log *zap.Logger) WarehouseServiceInterface {
	return &WarehouseService{
		repo:   repo,
		logger: log,
	}
}

func (s *WarehouseService) GetAllWarehouse(ctx context.Context, page, limit int) ([]model.Warehouse, *dto.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	warehouses, total, err := s.repo.GetAllWarehouse(ctx, page, limit)
	if err != nil {
		s.logger.Error("Failed to read warehouse list from repository", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
		TotalItems:  int64(total),
	}
	return warehouses, &pagination, nil
}

func (s *WarehouseService) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	return s.repo.CreateWarehouse(ctx, warehouse)
}

func (s *WarehouseService) UpdateWarehouse(ctx context.Context, id int, warehouse *model.Warehouse) error {
	existing, err := s.repo.GetWarehouseByID(ctx, id)
	if err != nil {
		return err
	}

	if warehouse.Name != "" {
		existing.Name = warehouse.Name
	}
	if warehouse.Location != "" {
		existing.Location = warehouse.Location
	}

	return s.repo.UpdateWarehouse(ctx, id, existing)
}

func (s *WarehouseService) DeleteWarehouse(ctx context.Context, id int) error {
	return s.repo.DeleteWarehouse(ctx, id)
}

func (s *WarehouseService) GetWarehouseById(ctx context.Context, id int) (*model.Warehouse, error) {
	warehouse, err := s.repo.GetWarehouseByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to read warehouse by id", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	return warehouse, nil
}
