package main

import (
	"app-inventory/database"
	"app-inventory/handler"
	"app-inventory/repository"
	"app-inventory/router"
	"app-inventory/service"
	"app-inventory/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 1. Load Configuration
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatalf("Fatal: Failed to load configuration: %v", err)
	}

	// 2. Initialize Logger
	logger, err := utils.InitLogger(config.PathLogging+"app-", config.Debug)
	if err != nil {
		log.Fatalf("Fatal: Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Inventory Management REST API...",
		zap.String("app_name", config.AppName),
		zap.String("port", config.Port),
	)

	// 3. Initialize PostgreSQL Database Connection Pool
	db, err := database.InitDB(config.DB)
	if err != nil {
		logger.Fatal("Fatal: Failed to connect to PostgreSQL", zap.Error(err))
	}
	defer db.Close()
	logger.Info("PostgreSQL connection pool established successfully")

	// 4. Initialize Redis Cache (Graceful Fallback)
	rdb := database.InitRedis(config.Redis, logger)
	if rdb != nil {
		defer rdb.Close()
	}

	// 5. Dependency Injection (Clean Architecture)
	repo := repository.AllRepo(db, rdb, logger)
	svc := service.AllService(repo, config.JWT, logger)
	h := handler.AllHandler(svc, logger)

	// 6. Router Setup (Gin Framework)
	r := router.NewRouter(h, svc, config.JWT, logger, config.Debug)

	serverAddr := fmt.Sprintf(":%s", config.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 7. Start HTTP Server in a separate goroutine
	go func() {
		logger.Info("HTTP Server is listening", zap.String("address", serverAddr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("HTTP server failed to listen", zap.Error(err))
		}
	}()

	// 8. Graceful Shutdown on SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down HTTP server gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("HTTP server cleanly stopped")
	}
}
