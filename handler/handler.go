package handler

import (
	"app-inventory/service"

	"go.uber.org/zap"
)

type Handler struct {
	User      UserHandler
	Inventory InventoryHandler
	Auth      AuthHandler
	log       *zap.Logger
}

func AllHandler(service service.Service, log *zap.Logger) Handler {
	return Handler{
		User:      NewUserHandler(service.UserService, log),
		Inventory: NewInventoryHandler(service.InventoryService, log),
		Auth:      NewAuthHandler(service.AuthService, log),
	}
}
