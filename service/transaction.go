package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"

	"go.uber.org/zap"
)

type TransactionService struct {
	repo   repository.Repo
	logger *zap.Logger
}
type TransactionServiceInterface interface {
	GetAllTransaction(page, limit int) ([]model.Transaction, *dto.Pagination, error)
	CreateTransaction(Transaction *model.Transaction) error
	UpdateTransaction(id int, Transaction *model.Transaction) error
	DeleteTransaction(id int) error
}

// constructor
func NewTransactionService(repo repository.Repo, log *zap.Logger) TransactionServiceInterface {
	return &TransactionService{
		repo:   repo,
		logger: log,
	}
}

func (s *TransactionService) CreateTransaction(Transaction *model.Transaction) error {
	err := s.repo.TransactionRepo.CreateTransaction(Transaction)
	if err != nil {
		s.logger.Error("failed create transaction on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *TransactionService) GetAllTransaction(page, limit int) ([]model.Transaction, *dto.Pagination, error) {
	Transaction, total, err := s.repo.TransactionRepo.GetAllTransaction(page, limit)
	if err != nil {
		s.logger.Error("failed to connect service to read list Transaction", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
	}
	return Transaction, &pagination, nil
}

func (s *TransactionService) UpdateTransaction(id int, Transaction *model.Transaction) error {
	err := s.repo.TransactionRepo.UpdateTransaction(id, Transaction)
	if err != nil {
		s.logger.Error("failed update transaction on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (s *TransactionService) DeleteTransaction(id int) error {
	err := s.repo.TransactionRepo.DeleteTransaction(id)
	if err != nil {
		s.logger.Error("failed delete transaction on repository",
			zap.Error(err),
		)
		return err
	}
	return nil
}
