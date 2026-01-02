package middleware

import (
	"app-inventory/service"

	"go.uber.org/zap"
)

type MiddlewareCostume struct {
	Service service.Service
	Log     *zap.Logger
}

func NewMiddlewareCustome(service service.Service, log *zap.Logger) MiddlewareCostume {
	return MiddlewareCostume{
		Service: service,
		Log:     log,
	}
}
