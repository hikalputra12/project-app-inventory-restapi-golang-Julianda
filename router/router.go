package router

import (
	"app-inventory/handler"
	mCostume "app-inventory/middleware"
	"app-inventory/service"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(handler handler.Handler, service service.Service, log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	// mw := mCostume.NewMiddlewareCustome(service, log)

	r.Mount("/api/v1", Apiv1(handler, service, log))

	return r
}

func Apiv1(handler handler.Handler, service service.Service, log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// middleware
	mw := mCostume.NewMiddlewareCustome(service, log)
	//authentication
	r.Post("/login", handler.Auth.Login)
	r.Post("/logout", handler.Auth.Logout)

	r.Route("/user", func(r chi.Router) {
		r.With(mw.RequirePermission("user:view")).Get("/", handler.User.ListUser)
		r.With(mw.RequirePermission("user:manage")).Post("/create", handler.User.CreateUser)
		r.With(mw.RequirePermission("user:manage")).Put("/update", handler.User.UpdateUser)
		r.With(mw.RequirePermission("user:manage")).Delete("/delete", handler.User.DeleteUser)
	})
	r.Route("/inventory", func(r chi.Router) {
		r.With(mw.RequirePermission("inventory:view")).Get("/", handler.Inventory.ListInventory)
		r.With(mw.RequirePermission("inventory:create")).Post("/create", handler.Inventory.CreateInventory)
		r.With(mw.RequirePermission("inventory:edit")).Put("/update", handler.Inventory.UpdateInventory)
		r.With(mw.RequirePermission("inventory:delete")).Delete("/delete", handler.Inventory.DeleteInventory)

	})
	r.Route("/category", func(r chi.Router) {
		r.With(mw.RequirePermission("category:view")).Get("/", handler.Category.ListCategory)
		r.With(mw.RequirePermission("category:manage")).Post("/create", handler.Category.CreateCategory)
		r.With(mw.RequirePermission("category:manage")).Put("/update", handler.Category.UpdateCategory)
		r.With(mw.RequirePermission("category:manage")).Delete("/delete", handler.Category.DeleteCategory)

	})
	r.Route("/rack", func(r chi.Router) {
		r.With(mw.RequirePermission("rack:view")).Get("/", handler.Rack.ListRack)
		r.With(mw.RequirePermission("rack:manage")).Post("/create", handler.Rack.CreateRack)
		r.With(mw.RequirePermission("rack:manage")).Put("/update", handler.Rack.UpdateRack)
		r.With(mw.RequirePermission("rack:manage")).Delete("/delete", handler.Rack.DeleteRack)

	})
	r.Route("/warehouse", func(r chi.Router) {
		r.With(mw.RequirePermission("warehouse:view")).Get("/", handler.Warehouse.ListWarehouse)
		r.With(mw.RequirePermission("warehouse:manage")).Post("/create", handler.Warehouse.CreateWarehouse)
		r.With(mw.RequirePermission("warehouse:manage")).Put("/update", handler.Warehouse.UpdateWarehouse)
		r.With(mw.RequirePermission("warehouse:manage")).Delete("/delete", handler.Warehouse.DeleteWarehouse)

	})
	r.Route("/transaction", func(r chi.Router) {
		r.With(mw.RequirePermission("transaction:manage")).Post("/create", handler.Transaction.CreateTransaction)

	})
	return r
}
