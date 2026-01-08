package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"

	"go.uber.org/zap"
)

type RackService struct {
	repo   repository.Repo
	logger *zap.Logger
}
type RackServiceInterface interface {
	GetAllRack(page, limit int) ([]model.Rack, *dto.Pagination, error)
	CreateRack(Rack *model.Rack) error
	UpdateRack(id int, Rack *model.Rack) error
	DeleteRack(id int) error
}

// constructor
func NewRackService(repo repository.Repo, log *zap.Logger) RackServiceInterface {
	return &RackService{
		repo:   repo,
		logger: log,
	}
}

func (s *RackService) GetAllRack(page, limit int) ([]model.Rack, *dto.Pagination, error) {
	Categories, total, err := s.repo.RackRepo.GetAllRack(page, limit)
	if err != nil {
		s.logger.Error("failed to connect service to read list Rack", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
	}
	return Categories, &pagination, nil
}

func (s *RackService) CreateRack(Rack *model.Rack) error {
	err := s.repo.RackRepo.CreateRack(Rack)
	if err != nil {
		s.logger.Error("failed create rack on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *RackService) UpdateRack(id int, Rack *model.Rack) error {
	err := s.repo.RackRepo.UpdateRack(id, Rack)
	if err != nil {
		s.logger.Error("failed upadate rack on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (s *RackService) DeleteRack(id int) error {
	err := s.repo.RackRepo.DeleteRack(id)
	if err != nil {
		s.logger.Error("failed delete rack on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}
