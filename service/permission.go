package service

// import (
// 		"app-inventory/repository"
// )

// type PermissionIface interface {
// 	Allowed(userID int, code string) (bool, error)
// }

// type permissionService struct {
// 	Repo repository.Repo
// }

// func NewPermissionService(repo repository.Repo) *permissionService {
// 	return &permissionService{Repo: repo}
// }

// func (permissionService *permissionService) Allowed(userID int, code string) (bool, error) {
// 	allowed, err := permissionService.Repo.PermissionRepository.Allowed(userID, code)
// 	if err != nil {
// 		return false, err
// 	}

// 	return allowed, nil
// }
