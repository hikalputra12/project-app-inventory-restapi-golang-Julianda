package dto

type RackListResponse struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	Warehouse_inventory_id int    `json:"warehouse_inventory_id"`
	WarehouseName          string `json:"warehouse_name,omitempty"`
}

type CreateRackRequest struct {
	Name                   string `json:"name" validate:"required,min=2"`
	Warehouse_inventory_id int    `json:"warehouse_inventory_id" validate:"required,gte=1"`
}

type UpdateRackRequest struct {
	Name                   string `json:"name" validate:"omitempty,min=2"`
	Warehouse_inventory_id int    `json:"warehouse_inventory_id" validate:"omitempty,gte=1"`
}

type RackByIdResponse struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	WarehouseInventoryId int    `json:"warehouse_inventory_id"`
	WarehouseInventory   string `json:"warehouse_inventory"`
}
