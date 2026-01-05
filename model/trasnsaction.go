package model

type Transaction struct {
	Model
	Name        string
	UserId      int
	InventoryId int
	Quantity    int
	Price       int
}
