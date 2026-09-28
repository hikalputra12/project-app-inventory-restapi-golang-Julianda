package model

type Transaction struct {
	Model
	Name        string `json:"name,omitempty"`
	UserId      int    `json:"user_id"`
	InventoryId int    `json:"inventory_id"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
}
