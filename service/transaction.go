package service

import (
	"app-inventory/model"
	"app-inventory/repository"

	"go.uber.org/zap"
)

type TransactionService struct {
	repo   repository.Repo
	logger *zap.Logger
}
type TransactionServiceInterface interface {
	CreateTransaction(Transaction *model.Transaction) error
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
		return err
	}
	return nil
}
