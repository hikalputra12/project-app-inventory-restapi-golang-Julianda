package router

import (
	"app-inventory/handler"
	"app-inventory/middleware"
	"app-inventory/service"
	"app-inventory/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter initializes Gin engine with enterprise middleware and routes
func NewRouter(h handler.Handler, s service.Service, jwtConfig utils.JWTConfig, log *zap.Logger, debug bool) *gin.Engine {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global Middlewares
	r.Use(gin.Recovery()) // Protects against panics and 500 crashes
	r.Use(middleware.Logging(log))
	r.Use(corsMiddleware())

	// Health Check Endpoint
	r.GET("/health", func(c *gin.Context) {
		utils.ResponseSuccess(c, http.StatusOK, "System is healthy", gin.H{
			"status": "up",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API v1 Routing
	api := r.Group("/api/v1")
	RegisterAPIV1Routes(api, h, s, jwtConfig, log)

	return r
}

func RegisterAPIV1Routes(rg *gin.RouterGroup, h handler.Handler, s service.Service, jwtConfig utils.JWTConfig, log *zap.Logger) {
	mw := middleware.NewCustomMiddleware(s, jwtConfig, log)

	// Public Authentication Routes
	rg.POST("/login", h.Auth.Login)
	rg.POST("/logout", h.Auth.Logout)

	// Protected Routes Group (Protected by JWT Authentication)
	protected := rg.Group("")
	protected.Use(mw.JWTAuth())
	{
		// User Management
		user := protected.Group("/user")
		{
			user.GET("", mw.RequirePermission("user:view"), h.User.ListUser)
			user.GET("/:id", mw.RequirePermission("user:view"), h.User.UserById)
			user.POST("", mw.RequirePermission("user:manage"), h.User.CreateUser)
			user.PATCH("/:id", mw.RequirePermission("user:manage"), h.User.UpdateUser)
			user.PUT("/:id", mw.RequirePermission("user:manage"), h.User.UpdateUser)
			user.DELETE("/:id", mw.RequirePermission("user:manage"), h.User.DeleteUser)
		}

		// Inventory Management
		inventory := protected.Group("/inventory")
		{
			inventory.GET("", mw.RequirePermission("inventory:view"), h.Inventory.ListInventory)
			inventory.GET("/:id", mw.RequirePermission("inventory:view"), h.Inventory.GetInventoryById)
			inventory.POST("", mw.RequirePermission("inventory:create"), h.Inventory.CreateInventory)
			inventory.PATCH("/:id", mw.RequirePermission("inventory:edit"), h.Inventory.UpdateInventory)
			inventory.PUT("/:id", mw.RequirePermission("inventory:edit"), h.Inventory.UpdateInventory)
			inventory.DELETE("/:id", mw.RequirePermission("inventory:delete"), h.Inventory.DeleteInventory)
		}

		// Category Management
		category := protected.Group("/category")
		{
			category.GET("", mw.RequirePermission("category:view"), h.Category.ListCategory)
			category.GET("/:id", mw.RequirePermission("category:view"), h.Category.GetCategoryById)
			category.POST("", mw.RequirePermission("category:manage"), h.Category.CreateCategory)
			category.PATCH("/:id", mw.RequirePermission("category:manage"), h.Category.UpdateCategory)
			category.PUT("/:id", mw.RequirePermission("category:manage"), h.Category.UpdateCategory)
			category.DELETE("/:id", mw.RequirePermission("category:manage"), h.Category.DeleteCategory)
		}

		// Storage Rack Management
		rack := protected.Group("/rack")
		{
			rack.GET("", mw.RequirePermission("rack:view"), h.Rack.ListRack)
			rack.GET("/:id", mw.RequirePermission("rack:view"), h.Rack.GetRackById)
			rack.POST("", mw.RequirePermission("rack:manage"), h.Rack.CreateRack)
			rack.PATCH("/:id", mw.RequirePermission("rack:manage"), h.Rack.UpdateRack)
			rack.PUT("/:id", mw.RequirePermission("rack:manage"), h.Rack.UpdateRack)
			rack.DELETE("/:id", mw.RequirePermission("rack:manage"), h.Rack.DeleteRack)
		}

		// Warehouse Management
		warehouse := protected.Group("/warehouse")
		{
			warehouse.GET("", mw.RequirePermission("warehouse:view"), h.Warehouse.ListWarehouse)
			warehouse.GET("/:id", mw.RequirePermission("warehouse:view"), h.Warehouse.GetWarehouseById)
			warehouse.POST("", mw.RequirePermission("warehouse:manage"), h.Warehouse.CreateWarehouse)
			warehouse.PATCH("/:id", mw.RequirePermission("warehouse:manage"), h.Warehouse.UpdateWarehouse)
			warehouse.PUT("/:id", mw.RequirePermission("warehouse:manage"), h.Warehouse.UpdateWarehouse)
			warehouse.DELETE("/:id", mw.RequirePermission("warehouse:manage"), h.Warehouse.DeleteWarehouse)
		}

		// Sales Transactions
		transaction := protected.Group("/transaction")
		{
			transaction.GET("", mw.RequirePermission("transaction:manage"), h.Transaction.ListTransaction)
			transaction.GET("/:id", mw.RequirePermission("transaction:manage"), h.Transaction.GetTransactionById)
			transaction.POST("", mw.RequirePermission("transaction:manage"), h.Transaction.CreateTransaction)
			transaction.PATCH("/:id", mw.RequirePermission("transaction:manage"), h.Transaction.UpdateTransaction)
			transaction.PUT("/:id", mw.RequirePermission("transaction:manage"), h.Transaction.UpdateTransaction)
			transaction.DELETE("/:id", mw.RequirePermission("transaction:manage"), h.Transaction.DeleteTransaction)
		}

		// Reporting & Stock Alert
		protected.GET("/report", mw.RequirePermission("report:view"), h.Report.Report)
		protected.GET("/stock", mw.RequirePermission("stock:view"), h.Inventory.CheckStock)
	}
}

// corsMiddleware configures Cross-Origin Resource Sharing for Gin
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
