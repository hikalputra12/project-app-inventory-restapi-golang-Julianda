package service

import (
	"app-inventory/repository"

	"go.uber.org/zap"
)

type Service struct {
	UserService        UserServiceInterface
	AuthService        AuthServiceInterface
	Permission         PermissionIface
	InventoryService   InventoryServiceInterface
	CategoryService    CategoryServiceInterface
	TransactionService TransactionServiceInterface
	RackService        RackServiceInterface
	WarehouseService   WarehouseServiceInterface
	ReportService      ReportServiceInterface
	log                *zap.Logger
}

func AllService(repo repository.Repo, log *zap.Logger) Service {
	return Service{
		UserService:        NewUserService(repo, log),
		InventoryService:   NewInventoryService(repo, log),
		CategoryService:    NewCategoryService(repo, log),
		RackService:        NewRackService(repo, log),
		WarehouseService:   NewWarehouseService(repo, log),
		TransactionService: NewTransactionService(repo, log),
		ReportService:      NewReportService(repo, log),
		AuthService:        NewAuthService(repo, log),
		Permission:         NewPermissionService(repo),
	}
}
