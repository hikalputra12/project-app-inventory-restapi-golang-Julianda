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
	SessionService     SessionServiceInterface
	log                *zap.Logger
}

func AllService(repo repository.Repo, log *zap.Logger) Service {
	return Service{
		UserService:        NewUserService(repo.UserRepo, log),
		InventoryService:   NewInventoryService(repo.InventoryRepo, log),
		CategoryService:    NewCategoryService(repo.CategoryRepo, log),
		RackService:        NewRackService(repo.RackRepo, log),
		WarehouseService:   NewWarehouseService(repo.WarehouseRepo, log),
		TransactionService: NewTransactionService(repo.TransactionRepo, log),
		ReportService:      NewReportService(repo.ReportRepo, log),
		SessionService:     NewSessionService(repo.SessionRepo, log),
		AuthService:        NewAuthService(repo.UserRepo, log),
		Permission:         NewPermissionService(repo.Permission),
	}
}
