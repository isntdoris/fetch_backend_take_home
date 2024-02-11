package model

type ReceiptItem struct {
	Description string `json:"shortDescription"`
	Price       string `json:"price"`
}

type ReceiptRequest struct {
	Retailer     string        `json:"retailer"`
	PurchaseDate string        `json:"purchaseDate"`
	PurchaseTime string        `json:"purchaseTime"`
	Items        []ReceiptItem `json:"items"`
	Total        string        `json:"total"`
}
