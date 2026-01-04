package dto

type InventoryListResponse struct {
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Stock     int    `json:"stock"`
	Category  string `json:"category"`
	Rack      string `json:"rack"`
	Warehouse string `json:"warehouse"`
}

type CreateInventoryRequest struct {
	Name        string `json:"name" validate:"required,min=6"`
	Price       int    `json:"price" validate:"required,gte=0"`
	Stock       int    `json:"stock"  validate:"required,gte=0"`
	Category_id int    `json:"category_id" validate:"required,gte=0"`
}
type UpdateInventoryRequest struct {
	Name         string `json:"name" validate:"required,min=6"`
	Price        int    `json:"price" validate:"required,gte=0"`
	Stock        int    `json:"stock"  validate:"required,gte=0"`
	Category_id  int    `json:"category_id" validate:"required,gte=0"`
	Inventory_id int    `json:"inventory_id" validate:"required,gte=0"`
}
type DeleteInventoryRequest struct {
	Inventory_id int `json:"inventory_id" validate:"required,gte=0"`
}
