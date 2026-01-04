package model

type Inventory struct {
	Model
	Name      string
	Price     int
	Stock     int
	Category  string
	Rack      string
	Warehouse string
}
