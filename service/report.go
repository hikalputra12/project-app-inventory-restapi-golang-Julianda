package service

import (
	"app-inventory/model"
	"app-inventory/repository"

	"go.uber.org/zap"
)

type ReportService struct {
	repo   repository.Repo
	logger *zap.Logger
}
type ReportServiceInterface interface {
	Report() (*model.Report, error)
}

// constructor
func NewReportService(repo repository.Repo, log *zap.Logger) ReportServiceInterface {
	return &ReportService{
		repo:   repo,
		logger: log,
	}
}

func (s *ReportService) Report() (*model.Report, error) {
	report, err := s.repo.ReportRepo.Report()
	if err != nil {
		return nil, err
	}
	return report, nil
}
