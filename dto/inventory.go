package dto

type InventoryListResponse struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Price               int    `json:"price"`
	Stock               int    `json:"stock"`
	Category            string `json:"category"`
	Rack                string `json:"rack"`
	Warehouse           string `json:"warehouse"`
	CategoryInventoryID int    `json:"category_inventory_id"`
}

type CreateInventoryRequest struct {
	Name        string `json:"name" validate:"required,min=3"`
	Price       int    `json:"price" validate:"required,gte=0"`
	Stock       int    `json:"stock" validate:"required,gte=0"`
	Category_id int    `json:"category_id" validate:"required,gte=1"`
}

type UpdateInventoryRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3"`
	Price       *int   `json:"price" validate:"omitempty,gte=0"`
	Stock       *int   `json:"stock" validate:"omitempty,gte=0"`
	Category_id *int   `json:"category_id" validate:"omitempty,gte=1"`
}

type InventoryByIdResponse struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Price               int    `json:"price"`
	Stock               int    `json:"stock"`
	CategoryInventoryID int    `json:"category_inventory_id"`
	Category            string `json:"category"`
	Rack                string `json:"rack"`
	Warehouse           string `json:"warehouse"`
}
