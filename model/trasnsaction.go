package model

type Transaction struct {
	Model
	UserId      int
	InventoryId int
	Quantity    int
}
