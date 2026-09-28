package service

import (
	"app-inventory/model"
	"app-inventory/repository"
	"context"

	"go.uber.org/zap"
)

type ReportService struct {
	repo   repository.ReportRepoInterface
	logger *zap.Logger
}

type ReportServiceInterface interface {
	Report(ctx context.Context) (*model.Report, error)
}

func NewReportService(repo repository.ReportRepoInterface, log *zap.Logger) ReportServiceInterface {
	return &ReportService{
		repo:   repo,
		logger: log,
	}
}

func (s *ReportService) Report(ctx context.Context) (*model.Report, error) {
	report, err := s.repo.Report(ctx)
	if err != nil {
		s.logger.Error("Failed to generate report on service", zap.Error(err))
		return nil, err
	}
	return report, nil
}
