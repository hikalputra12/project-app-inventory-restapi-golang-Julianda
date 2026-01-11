package dto

type WarehouseListResponse struct {
	Name string `json:"name"`
}

type CreateWarehouseRequest struct {
	Name     string `json:"name" validate:"required,min=6"`
	Location string `json:"location" validate:"required,min=6"`
}
type UpdateWarehouseRequest struct {
	Name     string `json:"name" validate:"required,min=6"`
	Location string `json:"location" validate:"required,min=6"`
}

type WarehouseByIdResponse struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}
