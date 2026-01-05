package dto

type CreateTransactionRequest struct {
	UserId      int `json:"user_id" validate:"required,gte=0"`
	InventoryId int `json:"inventory_id"  validate:"required,gte=0"`
	Quantity    int `json:"quantity" validate:"required,gte=0"`
	Price       int `json:"price" validate:"required,gte=0"`
}

type TransactionListResponse struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type UpdateTransactionRequest struct {
	Quantity    int `json:"quantity" validate:"required,gte=0"`
	SalesItemID int `json:"sales_item_ID" validate:"required"`
}
type DeleteTransactionRequest struct {
	SalesItemID int `json:"sales_item_ID" validate:"required,gte=0"`
}
