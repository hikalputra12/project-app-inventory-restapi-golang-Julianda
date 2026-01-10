package router

import (
	"app-inventory/handler"
	mCostume "app-inventory/middleware"
	"app-inventory/service"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(handler handler.Handler, service service.Service, log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/api/v1", Apiv1(handler, service, log))

	return r
}

func Apiv1(handler handler.Handler, service service.Service, log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	// middleware
	mw := mCostume.NewMiddlewareCustome(service, log)
	r.Use(mCostume.Logging(log))
	//authentication
	r.Post("/login", handler.Auth.Login)
	r.Post("/logout", handler.Auth.Logout)

	r.Route("/user", func(r chi.Router) {
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("user:view")).Get("/", handler.User.ListUser)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("user:manage")).Post("/", handler.User.CreateUser)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("user:manage")).Put("/{id}", handler.User.UpdateUser)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("user:manage")).Delete("/{id}", handler.User.DeleteUser)
	})
	r.Route("/inventory", func(r chi.Router) {
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("inventory:view")).Get("/", handler.Inventory.ListInventory)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("inventory:create")).Post("/", handler.Inventory.CreateInventory)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("inventory:edit")).Put("/{id}", handler.Inventory.UpdateInventory)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("inventory:delete")).Delete("/{id}", handler.Inventory.DeleteInventory)

	})
	r.Route("/category", func(r chi.Router) {
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("category:view")).Get("/", handler.Category.ListCategory)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("category:manage")).Post("/", handler.Category.CreateCategory)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("category:manage")).Put("/{id}", handler.Category.UpdateCategory)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("category:manage")).Delete("/{id}", handler.Category.DeleteCategory)

	})
	r.Route("/rack", func(r chi.Router) {
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("rack:view")).Get("/", handler.Rack.ListRack)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("rack:manage")).Post("/", handler.Rack.CreateRack)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("rack:manage")).Put("/{id}", handler.Rack.UpdateRack)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("rack:manage")).Delete("/{id}", handler.Rack.DeleteRack)

	})
	r.Route("/warehouse", func(r chi.Router) {
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("warehouse:view")).Get("/", handler.Warehouse.ListWarehouse)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("warehouse:manage")).Post("/", handler.Warehouse.CreateWarehouse)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("warehouse:manage")).Put("/{id}", handler.Warehouse.UpdateWarehouse)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("warehouse:manage")).Delete("/{id}", handler.Warehouse.DeleteWarehouse)

	})
	r.Route("/transaction", func(r chi.Router) {
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("transaction:manage")).Get("/", handler.Transaction.ListTransaction)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("transaction:manage")).Post("/", handler.Transaction.CreateTransaction)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("transaction:manage")).Put("/{id}", handler.Transaction.UpdateTransaction)
		r.With(mw.ValidAndExtendSession(), mw.RequirePermission("transaction:manage")).Delete("/{id}", handler.Transaction.DeleteTransaction)

	})
	r.With(mw.ValidAndExtendSession(), mw.RequirePermission("report:view")).Get("/report", handler.Report.Report)
	r.With(mw.ValidAndExtendSession(), mw.RequirePermission("stock:view")).Get("/stock", handler.Inventory.CheckStock)

	return r
}
