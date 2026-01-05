package service

import (
	"app-inventory/repository"

	"go.uber.org/zap"
)

type Service struct {
	UserService      UserServiceInterface
	AuthService      AuthServiceInterface
	Permission       PermissionIface
	InventoryService InventoryServiceInterface
	CategoryService  CategoryServiceInterface
	RackService      RackServiceInterface
	WarehouseService WarehouseServiceInterface
	log              *zap.Logger
}

func AllService(repo repository.Repo, log *zap.Logger) Service {
	return Service{
		UserService:      NewUserService(repo, log),
		InventoryService: NewInventoryService(repo, log),
		CategoryService:  NewCategoryService(repo, log),
		RackService:      NewRackService(repo, log),
		WarehouseService: NewWarehouseService(repo, log),
		AuthService:      NewAuthService(repo, log),
		Permission:       NewPermissionService(repo),
	}
}
