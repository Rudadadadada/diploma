package models

import "time"

type BucketItem struct {
	Id        uint    `json:"id"`
	ProductId uint    `json:"product_id"`
	Name      string  `json:"name"`
	Amount    uint    `json:"amount"`
	TotalCost float32 `json:"total_cost"`
}

type Courier struct {
	Id              int `json:"id"`
	Active          bool `json:"active"`
	In_progress     bool `json:"in_progress"`
	Rating          int  `json:"rating"`
	Order_delivered int  `json:"order_delivered"`
}

type OrderMessage struct {
	OrderId           int          `json:"order_id"`
	CustomerId        int          `json:"customer_id"`
	TotalCost         float32      `json:"total_cost"`
	Status            string       `json:"status"`
	CreatedAt         time.Time    `json:"created_at"`
	DeliveryStartedAt time.Time    `json:"delivery_started_at"`
	DeliveredAt       time.Time    `json:"delivered_at,omitempty"`
	Courier           Courier      `json:"courier"`
	OrderItems        []BucketItem `json:"bucket_items"`
}

