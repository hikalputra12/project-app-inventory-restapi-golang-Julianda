package dto

type WarehouseListResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type CreateWarehouseRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Location string `json:"location" validate:"required,min=2"`
}

type UpdateWarehouseRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2"`
	Location string `json:"location" validate:"omitempty,min=2"`
}

type WarehouseByIdResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}
