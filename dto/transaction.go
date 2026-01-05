package dto

type CreateTransactionRequest struct {
	UserId      int `json:"user_id" validate:"required,gte=0"`
	InventoryId int `json:"inventory_id"  validate:"required,gte=0"`
	Quantity    int `json:"quantity" validate:"required,gte=0"`
	Price       int `json:"price" validate:"required,gte=0"`
}
