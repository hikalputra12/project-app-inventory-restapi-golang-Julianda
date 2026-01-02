package handler

import (
	"app-inventory/service"

	"go.uber.org/zap"
)

type Handler struct {
	User UserHandler
	Auth AuthHandler
	log  *zap.Logger
}

func AllHandler(service service.Service, log *zap.Logger) Handler {
	return Handler{
		User: NewUserHandler(service.UserService, log),
		Auth: NewAuthHandler(service.AuthService, log),
	}
}
