package dto

type CategoryListResponse struct {
	Name              string `json:"name"`
	Rack_inventory_id int    `json:"rack_inventory_id "`
}

type CreateCategoryRequest struct {
	Name              string `json:"name" validate:"required,min=6"`
	Rack_inventory_id int    `json:"rack_inventory_id" validate:"required,gte=0"`
}
type UpdateCategoryRequest struct {
	Name              string `json:"name" validate:"required,min=6"`
	Rack_inventory_id int    `json:"rack_inventory_id" validate:"required,gte=0"`
}
type CategoryByIdResponse struct {
	Name            string `json:"name"`
	RackInventory   string `json:"rack_inventory"`
	RackInventoryId int    `json:"-"`
}
