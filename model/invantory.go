package model

type Inventory struct {
	Model
	Name                  string `json:"name"`
	Price                 int    `json:"price"`
	Stock                 int    `json:"stock"`
	Category              string `json:"category,omitempty"`
	Rack                  string `json:"rack,omitempty"`
	Warehouse             string `json:"warehouse,omitempty"`
	Category_inventory_id int    `json:"category_inventory_id"`
}
