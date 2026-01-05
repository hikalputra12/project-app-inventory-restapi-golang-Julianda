package dto

type RackListResponse struct {
	Name                   string `json:"name"`
	Warehouse_inventory_id int    `json:"warehouse_inventory_id"`
}

type CreateRackRequest struct {
	Name                   string `json:"name" validate:"required,min=6"`
	Warehouse_inventory_id int    `json:"warehouse_inventory_id" validate:"required,gte=0"`
}
type UpdateRackRequest struct {
	Name                   string `json:"name" validate:"required,min=6"`
	Warehouse_inventory_id int    `json:"warehouse_inventory_id" validate:"required,gte=0"`
}
