package dto

type CategoryListResponse struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Rack_inventory_id int    `json:"rack_inventory_id"`
	RackName          string `json:"rack_name,omitempty"`
}

type CreateCategoryRequest struct {
	Name              string `json:"name" validate:"required,min=2"`
	Rack_inventory_id int    `json:"rack_inventory_id" validate:"required,gte=1"`
}

type UpdateCategoryRequest struct {
	Name              string `json:"name" validate:"omitempty,min=2"`
	Rack_inventory_id int    `json:"rack_inventory_id" validate:"omitempty,gte=1"`
}

type CategoryByIdResponse struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	RackInventory   string `json:"rack_inventory"`
	RackInventoryId int    `json:"rack_inventory_id"`
}
