package repository

import (
	"app-inventory/database"

	"go.uber.org/zap"
)

type Repo struct {
	UserRepo      UserRepoInterface
	InventoryRepo InventoryRepoInterface
	CategoryRepo  CategoryRepoInterface
	Permission    PermissionIface
	Log           *zap.Logger
}

func AllRepo(db database.PgxIface, log *zap.Logger) Repo {
	return Repo{
		UserRepo:      NewUserRepo(db, log),
		InventoryRepo: NewInventoryRepo(db, log),
		CategoryRepo:  NewCategoryRepo(db, log),
		Permission:    NewPermissionRepository(db),
	}
}
