package dto

type ReportResponse struct {
	TotalTransactions int `json:"total_transactions"`
	TotalItemsSold    int `json:"total_item_sold"`
	TotalRevenue      int `json:"total_revenue"`
}
