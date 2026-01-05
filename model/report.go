package model

type Report struct {
	TotalTransactions int `json:"total_transactions"`
	TotalItemsSold    int `json:"total_items_sold"`
	TotalRevenue      int `json:"total_revenue"`
}
