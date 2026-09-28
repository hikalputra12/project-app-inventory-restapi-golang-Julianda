package dto

type CreateTransactionRequest struct {
	InventoryId int `json:"inventory_id" validate:"required,gte=1"`
	Quantity    int `json:"quantity" validate:"required,gt=0"`
}

type TransactionListResponse struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	InventoryID int    `json:"inventory_id"`
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	TotalPrice  int    `json:"total_price"`
	CreatedAt   string `json:"created_at"`
}

type UpdateTransactionRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

type TransactionByIdResponse struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	InventoryID int    `json:"inventory_id"`
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	TotalPrice  int    `json:"total_price"`
	CreatedAt   string `json:"created_at"`
}
