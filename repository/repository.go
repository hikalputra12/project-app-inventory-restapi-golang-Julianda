package repository

import (
	"app-inventory/database"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Repo struct {
	UserRepo        UserRepoInterface
	InventoryRepo   InventoryRepoInterface
	CategoryRepo    CategoryRepoInterface
	RackRepo        RackRepoInterface
	WarehouseRepo   WarehouseRepoInterface
	TransactionRepo TransactionRepoInterface
	ReportRepo      ReportRepoInterface
	SessionRepo     SessionRepoInterface
	Permission      PermissionIface
	Redis           *redis.Client
	Log             *zap.Logger
}

func AllRepo(db database.PgxIface, rdb *redis.Client, log *zap.Logger) Repo {
	return Repo{
		UserRepo:        NewUserRepo(db, log),
		InventoryRepo:   NewInventoryRepo(db, log),
		CategoryRepo:    NewCategoryRepo(db, log),
		RackRepo:        NewRackRepo(db, log),
		WarehouseRepo:   NewWarehouseRepo(db, log),
		TransactionRepo: NewTransactionRepo(db, log),
		ReportRepo:      NewReportRepo(db, rdb, log),
		SessionRepo:     NewSessionRepo(db, log),
		Permission:      NewPermissionRepository(db, rdb),
		Redis:           rdb,
		Log:             log,
	}
}
