package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"
	"context"

	"go.uber.org/zap"
)

type InventoryService struct {
	repo       repository.InventoryRepoInterface
	reportRepo repository.ReportRepoInterface
	logger     *zap.Logger
}

type InventoryServiceInterface interface {
	GetAllInventory(ctx context.Context, page, limit int) ([]model.Inventory, *dto.Pagination, error)
	CheckStock(ctx context.Context, page, limit int) ([]model.Inventory, *dto.Pagination, error)
	CreateInventory(ctx context.Context, inventory *model.Inventory) error
	UpdateInventory(ctx context.Context, id int, req *dto.UpdateInventoryRequest) error
	DeleteInventory(ctx context.Context, id int) error
	GetInventoryById(ctx context.Context, id int) (*model.Inventory, error)
}

func NewInventoryService(repo repository.InventoryRepoInterface, reportRepo repository.ReportRepoInterface, log *zap.Logger) InventoryServiceInterface {
	return &InventoryService{
		repo:       repo,
		reportRepo: reportRepo,
		logger:     log,
	}
}

func (s *InventoryService) GetAllInventory(ctx context.Context, page, limit int) ([]model.Inventory, *dto.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	inventories, total, err := s.repo.GetAllInventory(ctx, page, limit)
	if err != nil {
		s.logger.Error("Failed to fetch inventory list from repository", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
		TotalItems:  int64(total),
	}
	return inventories, &pagination, nil
}

func (s *InventoryService) CreateInventory(ctx context.Context, inventory *model.Inventory) error {
	err := s.repo.CreateInventory(ctx, inventory)
	if err != nil {
		s.logger.Error("Failed to create inventory in repository", zap.Error(err))
		return err
	}
	s.reportRepo.InvalidateReportCache(ctx)
	return nil
}

func (s *InventoryService) UpdateInventory(ctx context.Context, id int, req *dto.UpdateInventoryRequest) error {
	existing, err := s.repo.GetInventoryByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	// Fixed: Pointers allow setting price or stock to 0
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.Stock != nil {
		existing.Stock = *req.Stock
	}
	if req.Category_id != nil && *req.Category_id > 0 {
		existing.Category_inventory_id = *req.Category_id
	}

	err = s.repo.UpdateInventory(ctx, id, existing)
	if err != nil {
		s.logger.Error("Failed to update inventory in repository", zap.Error(err))
		return err
	}

	s.reportRepo.InvalidateReportCache(ctx)
	return nil
}

func (s *InventoryService) DeleteInventory(ctx context.Context, id int) error {
	err := s.repo.DeleteInventory(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete inventory in repository", zap.Error(err))
		return err
	}
	s.reportRepo.InvalidateReportCache(ctx)
	return nil
}

func (s *InventoryService) CheckStock(ctx context.Context, page, limit int) ([]model.Inventory, *dto.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	inventories, total, err := s.repo.CheckStock(ctx, page, limit)
	if err != nil {
		s.logger.Error("Failed to check low stock inventory", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
		TotalItems:  int64(total),
	}
	return inventories, &pagination, nil
}

func (s *InventoryService) GetInventoryById(ctx context.Context, id int) (*model.Inventory, error) {
	inventory, err := s.repo.GetInventoryByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch inventory by id", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	return inventory, nil
}
