package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"

	"go.uber.org/zap"
)

type WarehouseService struct {
	repo   repository.Repo
	logger *zap.Logger
}
type WarehouseServiceInterface interface {
	GetAllWarehouse(page, limit int) ([]model.Warehouse, *dto.Pagination, error)
	CreateWarehouse(Warehouse *model.Warehouse) error
	UpdateWarehouse(id int, Warehouse *model.Warehouse) error
	DeleteWarehouse(id int) error
}

// constructor
func NewWarehouseService(repo repository.Repo, log *zap.Logger) WarehouseServiceInterface {
	return &WarehouseService{
		repo:   repo,
		logger: log,
	}
}

func (s *WarehouseService) GetAllWarehouse(page, limit int) ([]model.Warehouse, *dto.Pagination, error) {
	warehouse, total, err := s.repo.WarehouseRepo.GetAllWarehouse(page, limit)
	if err != nil {
		s.logger.Error("failed to connect service to read list Warehouse", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
	}
	return warehouse, &pagination, nil
}

func (s *WarehouseService) CreateWarehouse(Warehouse *model.Warehouse) error {
	err := s.repo.WarehouseRepo.CreateWarehouse(Warehouse)
	if err != nil {
		return err
	}
	return nil
}

func (s *WarehouseService) UpdateWarehouse(id int, Warehouse *model.Warehouse) error {
	err := s.repo.WarehouseRepo.UpdateWarehouse(id, Warehouse)
	if err != nil {
		return err
	}
	return nil
}
func (s *WarehouseService) DeleteWarehouse(id int) error {
	err := s.repo.WarehouseRepo.DeleteWarehouse(id)
	if err != nil {
		return err
	}
	return nil
}
