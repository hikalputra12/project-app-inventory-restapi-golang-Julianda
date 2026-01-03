package repository

import (
	"app-inventory/database"

	"go.uber.org/zap"
)

type Repo struct {
	UserRepo   UserRepoInterface
	Permission PermissionIface
	Log        *zap.Logger
}

func AllRepo(db database.PgxIface, log *zap.Logger) Repo {
	return Repo{
		UserRepo:   NewUserRepo(db, log),
		Permission: NewPermissionRepository(db),
	}
}
