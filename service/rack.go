package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"
	"context"

	"go.uber.org/zap"
)

type RackService struct {
	repo   repository.RackRepoInterface
	logger *zap.Logger
}

type RackServiceInterface interface {
	GetAllRack(ctx context.Context, page, limit int) ([]model.Rack, *dto.Pagination, error)
	CreateRack(ctx context.Context, rack *model.Rack) error
	UpdateRack(ctx context.Context, id int, rack *model.Rack) error
	DeleteRack(ctx context.Context, id int) error
	GetRackById(ctx context.Context, id int) (*model.Rack, error)
}

func NewRackService(repo repository.RackRepoInterface, log *zap.Logger) RackServiceInterface {
	return &RackService{
		repo:   repo,
		logger: log,
	}
}

func (s *RackService) GetAllRack(ctx context.Context, page, limit int) ([]model.Rack, *dto.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	racks, total, err := s.repo.GetAllRack(ctx, page, limit)
	if err != nil {
		s.logger.Error("Failed to read rack list from repository", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
		TotalItems:  int64(total),
	}
	return racks, &pagination, nil
}

func (s *RackService) CreateRack(ctx context.Context, rack *model.Rack) error {
	err := s.repo.CreateRack(ctx, rack)
	if err != nil {
		s.logger.Error("Failed to create rack on repository", zap.Error(err))
		return err
	}
	return nil
}

func (s *RackService) UpdateRack(ctx context.Context, id int, rack *model.Rack) error {
	existing, err := s.repo.GetRackByID(ctx, id)
	if err != nil {
		return err
	}

	if rack.Name != "" {
		existing.Name = rack.Name
	}
	if rack.WarehouseInventoryId != 0 {
		existing.WarehouseInventoryId = rack.WarehouseInventoryId
	}

	err = s.repo.UpdateRack(ctx, id, existing)
	if err != nil {
		s.logger.Error("Failed to update rack on repository", zap.Error(err), zap.Int("id", id))
		return err
	}
	return nil
}

func (s *RackService) DeleteRack(ctx context.Context, id int) error {
	err := s.repo.DeleteRack(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete rack on repository", zap.Error(err), zap.Int("id", id))
		return err
	}
	return nil
}

func (s *RackService) GetRackById(ctx context.Context, id int) (*model.Rack, error) {
	rack, err := s.repo.GetRackByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to read rack by id", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	return rack, nil
}
