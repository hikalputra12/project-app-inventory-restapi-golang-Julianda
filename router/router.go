package router

import (
	"app-inventory/handler"
	"app-inventory/service"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(handler handler.Handler, service service.Service, log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	// mw := mCostume.NewMiddlewareCustome(service, log)

	r.Mount("/api/v1", Apiv1(handler))

	return r
}

func Apiv1(handler handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//authentication
	r.Post("/login", handler.Auth.Login)
	r.Post("/logout", handler.Auth.Logout)

	r.Route("/user", func(r chi.Router) {
		r.Get("/", handler.User.ListUser)
	})

	return r
}
