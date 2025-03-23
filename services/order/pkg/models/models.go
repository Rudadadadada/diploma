package models

import "time"

type BucketItem struct {
	Id        uint    `json:"id"`
	ProductId uint    `json:"product_id"`
	Name      string  `json:"name"`
	Amount    uint    `json:"amount"`
	TotalCost float32 `json:"total_cost"`
}

type OrderMessage struct {
	OrderId     int          `json:"order_id"`
	CustomerId  int          `json:"customer_id"`
	TotalCost   int          `json:"total_cost"`
	Status      string       `json:"status"`
	DeliveredAt *time.Time   `json:"delivered_at,omitempty"`
	OrderItems  []BucketItem `json:"bucket_items"`
}
