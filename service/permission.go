package service

import (
	"app-inventory/repository"
)

type PermissionIface interface {
	Allowed(userID int, code string) (bool, error)
}

type permissionService struct {
	Repo repository.PermissionIface
}

func NewPermissionService(repo repository.PermissionIface) *permissionService {
	return &permissionService{Repo: repo}
}

func (permissionService *permissionService) Allowed(userID int, code string) (bool, error) {
	allowed, err := permissionService.Repo.Allowed(userID, code)
	if err != nil {
		return false, err
	}

	return allowed, nil
}
