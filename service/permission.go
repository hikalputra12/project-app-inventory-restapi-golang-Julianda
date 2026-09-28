package service

import (
	"app-inventory/repository"
	"context"
)

type PermissionIface interface {
	Allowed(ctx context.Context, userID int, code string) (bool, error)
}

type permissionService struct {
	Repo repository.PermissionIface
}

func NewPermissionService(repo repository.PermissionIface) PermissionIface {
	return &permissionService{Repo: repo}
}

func (s *permissionService) Allowed(ctx context.Context, userID int, code string) (bool, error) {
	return s.Repo.Allowed(ctx, userID, code)
}
