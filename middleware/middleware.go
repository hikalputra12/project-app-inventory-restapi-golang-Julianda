package middleware

import (
	"app-inventory/service"
	"app-inventory/utils"

	"go.uber.org/zap"
)

type CustomMiddleware struct {
	Service   service.Service
	JWTConfig utils.JWTConfig
	Log       *zap.Logger
}

// Backwards compatibility alias
type MiddlewareCostume = CustomMiddleware

func NewCustomMiddleware(service service.Service, jwtConfig utils.JWTConfig, log *zap.Logger) CustomMiddleware {
	return CustomMiddleware{
		Service:   service,
		JWTConfig: jwtConfig,
		Log:       log,
	}
}

// Backwards compatibility constructor
func NewMiddlewareCustome(service service.Service, log *zap.Logger) CustomMiddleware {
	return CustomMiddleware{
		Service: service,
		Log:     log,
	}
}
