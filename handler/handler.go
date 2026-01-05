package handler

import (
	"app-inventory/service"

	"go.uber.org/zap"
)

type Handler struct {
	User      UserHandler
	Inventory InventoryHandler
	Category  CategoryHandler
	Rack      RackHandler
	Warehouse WarehouseHandler
	Auth      AuthHandler
	log       *zap.Logger
}

func AllHandler(service service.Service, log *zap.Logger) Handler {
	return Handler{
		User:      NewUserHandler(service.UserService, log),
		Inventory: NewInventoryHandler(service.InventoryService, log),
		Category:  NewCategoryHandler(service.CategoryService, log),
		Rack:      NewRackHandler(service.RackService, log),
		Warehouse: NewWarehouseHandler(service.WarehouseService, log),
		Auth:      NewAuthHandler(service.AuthService, log),
	}
}
