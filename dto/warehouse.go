package dto

type WarehouseListResponse struct {
	Name string `json:"name"`
}

type CreateWarehouseRequest struct {
	Name string `json:"name" validate:"required,min=6"`
}
type UpdateWarehouseRequest struct {
	Name string `json:"name" validate:"required,min=6"`
}
