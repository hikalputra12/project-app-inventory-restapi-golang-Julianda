package dto

type WarehouseListResponse struct {
	Name string `json:"name"`
}

type CreateWarehouseRequest struct {
	Name string `json:"name" validate:"required,min=6"`
}
type UpdateWarehouseRequest struct {
	Name                   string `json:"name" validate:"required,min=6"`
	Warehouse_inventory_id int    `json:"warehouse_inventory_id" validate:"required,gte=0"`
}
type DeleteWarehouseRequest struct {
	Warehouse_Inventory_id int `json:"Warehouse_inventory_id" validate:"required,gte=0"`
}
