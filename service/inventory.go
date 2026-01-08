package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"

	"go.uber.org/zap"
)

type InventoryService struct {
	repo   repository.Repo
	logger *zap.Logger
}
type InventoryServiceInterface interface {
	GetAllInventory(page, limit int) ([]model.Inventory, *dto.Pagination, error)
	CheckStock(page, limit int) ([]model.Inventory, *dto.Pagination, error)
	CreateInventory(inventory *model.Inventory) error
	UpdateInventory(id int, inventory *model.Inventory) error
	DeleteInventory(id int) error
}

// constructor
func NewInventoryService(repo repository.Repo, log *zap.Logger) InventoryServiceInterface {
	return &InventoryService{
		repo:   repo,
		logger: log,
	}
}

func (s *InventoryService) GetAllInventory(page, limit int) ([]model.Inventory, *dto.Pagination, error) {
	Inventories, total, err := s.repo.InventoryRepo.GetAllInventory(page, limit)
	if err != nil {
		s.logger.Error("failed to connect service to read list Inventory", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
	}
	return Inventories, &pagination, nil
}

func (s *InventoryService) CreateInventory(inventory *model.Inventory) error {
	err := s.repo.InventoryRepo.CreateInventory(inventory)
	if err != nil {
		s.logger.Error("failed created inventory on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *InventoryService) UpdateInventory(id int, inventory *model.Inventory) error {
	err := s.repo.InventoryRepo.UpdateInventory(id, inventory)
	if err != nil {
		s.logger.Error("failed created category on repository ",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (s *InventoryService) DeleteInventory(id int) error {
	err := s.repo.InventoryRepo.DeleteInventory(id)
	if err != nil {
		s.logger.Error("failed delete inventory on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *InventoryService) CheckStock(page, limit int) ([]model.Inventory, *dto.Pagination, error) {
	Inventories, total, err := s.repo.InventoryRepo.CheckStock(page, limit)
	if err != nil {
		s.logger.Error("failed to connect service to read list Inventory", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
	}
	return Inventories, &pagination, nil
}
