package service

import (
	"app-inventory/model"
	"app-inventory/repository"

	"go.uber.org/zap"
)

type ReportService struct {
	repo   repository.ReportRepoInterface
	logger *zap.Logger
}
type ReportServiceInterface interface {
	Report() (*model.Report, error)
}

// constructor
func NewReportService(repo repository.ReportRepoInterface, log *zap.Logger) ReportServiceInterface {
	return &ReportService{
		repo:   repo,
		logger: log,
	}
}

func (s *ReportService) Report() (*model.Report, error) {
	report, err := s.repo.Report()
	if err != nil {
		s.logger.Error("failed create report on repository",
			zap.Error(err),
		)
		return nil, err
	}
	return report, nil
}
