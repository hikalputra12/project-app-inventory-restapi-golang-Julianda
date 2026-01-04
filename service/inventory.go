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
