package dto

type InventoryListResponse struct {
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Stock     int    `json:"stock"`
	Category  string `json:"category"`
	Rack      string `json:"rack"`
	Warehouse string `json:"warehouse"`
}
