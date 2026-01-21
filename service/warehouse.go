package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"

	"go.uber.org/zap"
)

type WarehouseService struct {
	repo   repository.WarehouseRepoInterface
	logger *zap.Logger
}
type WarehouseServiceInterface interface {
	GetAllWarehouse(page, limit int) ([]model.Warehouse, *dto.Pagination, error)
	CreateWarehouse(Warehouse *model.Warehouse) error
	UpdateWarehouse(id int, Warehouse *model.Warehouse) error
	DeleteWarehouse(id int) error
	GetWarehouseById(id int) (*model.Warehouse, error)
}

// constructor
func NewWarehouseService(repo repository.WarehouseRepoInterface, log *zap.Logger) WarehouseServiceInterface {
	return &WarehouseService{
		repo:   repo,
		logger: log,
	}
}

func (s *WarehouseService) GetAllWarehouse(page, limit int) ([]model.Warehouse, *dto.Pagination, error) {
	warehouse, total, err := s.repo.GetAllWarehouse(page, limit)
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
	err := s.repo.CreateWarehouse(Warehouse)
	if err != nil {
		return err
	}
	return nil
}

func (s *WarehouseService) UpdateWarehouse(id int, Warehouse *model.Warehouse) error {
	getWarehouse, err := s.repo.GetWarehouseByID(id)
	if err != nil {
		return err
	}

	if Warehouse.Name != "" {
		getWarehouse.Name = Warehouse.Name
	}
	if Warehouse.Location != "" {
		getWarehouse.Location = Warehouse.Location
	}

	NewWarehouseUpdate := &model.Warehouse{
		Name:     getWarehouse.Name,
		Location: getWarehouse.Location,
	}
	err = s.repo.UpdateWarehouse(id, NewWarehouseUpdate)
	if err != nil {
		return err
	}
	return nil
}
func (s *WarehouseService) DeleteWarehouse(id int) error {
	err := s.repo.DeleteWarehouse(id)
	if err != nil {
		return err
	}
	return nil
}

// get warehouse by id
func (s *WarehouseService) GetWarehouseById(id int) (*model.Warehouse, error) {
	warehouse, err := s.repo.GetWarehouseByID(id)
	if err != nil {
		s.logger.Error("failed to connect service to read warehouse by id", zap.Error(err))
		return nil, err
	}

	return warehouse, nil
}
