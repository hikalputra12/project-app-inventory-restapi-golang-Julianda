package service

import (
	"app-inventory/dto"
	"app-inventory/model"
	"app-inventory/repository"
	"app-inventory/utils"
	"context"

	"go.uber.org/zap"
)

type TransactionService struct {
	repo       repository.TransactionRepoInterface
	reportRepo repository.ReportRepoInterface
	logger     *zap.Logger
}

type TransactionServiceInterface interface {
	GetAllTransaction(ctx context.Context, page, limit int) ([]model.Transaction, *dto.Pagination, error)
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	UpdateTransaction(ctx context.Context, id int, transaction *model.Transaction) error
	DeleteTransaction(ctx context.Context, id int) error
	GetTransactionById(ctx context.Context, id int) (*model.Transaction, error)
}

func NewTransactionService(repo repository.TransactionRepoInterface, reportRepo repository.ReportRepoInterface, log *zap.Logger) TransactionServiceInterface {
	return &TransactionService{
		repo:       repo,
		reportRepo: reportRepo,
		logger:     log,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	err := s.repo.CreateTransaction(ctx, transaction)
	if err != nil {
		s.logger.Error("Failed to create transaction in repository", zap.Error(err))
		return err
	}
	s.reportRepo.InvalidateReportCache(ctx)
	return nil
}

func (s *TransactionService) GetAllTransaction(ctx context.Context, page, limit int) ([]model.Transaction, *dto.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	transactions, total, err := s.repo.GetAllTransaction(ctx, page, limit)
	if err != nil {
		s.logger.Error("Failed to read transaction list from repository", zap.Error(err))
		return nil, nil, err
	}
	pagination := dto.Pagination{
		CurrentPage: page,
		Limit:       limit,
		TotalPages:  utils.TotalPage(limit, int64(total)),
		TotalItems:  int64(total),
	}
	return transactions, &pagination, nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, id int, transaction *model.Transaction) error {
	err := s.repo.UpdateTransaction(ctx, id, transaction)
	if err != nil {
		s.logger.Error("Failed to update transaction in repository", zap.Error(err), zap.Int("id", id))
		return err
	}
	s.reportRepo.InvalidateReportCache(ctx)
	return nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, id int) error {
	err := s.repo.DeleteTransaction(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete transaction in repository", zap.Error(err), zap.Int("id", id))
		return err
	}
	s.reportRepo.InvalidateReportCache(ctx)
	return nil
}

func (s *TransactionService) GetTransactionById(ctx context.Context, id int) (*model.Transaction, error) {
	transaction, err := s.repo.GetTransactionById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to read transaction by id from repository", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	return transaction, nil
}
