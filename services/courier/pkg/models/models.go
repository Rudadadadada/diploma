package models

import "time"

type BucketItem struct {
	Id        uint    `json:"id"`
	ProductId uint    `json:"product_id"`
	Name      string  `json:"name"`
	Amount    uint    `json:"amount"`
	TotalCost float32 `json:"total_cost"`
}

type Order struct {
	Id           uint        `json:"id"`
	TotalCost    float64     `json:"total_cost"`
	CreatedAt    time.Time   `json:"created_at"`
	CreatedAtStr string      `json:"created_at_str"`
	OrderItems   []OrderItem `json:"order_items"`
	Status       string      `json:"status"`
}

type OrderItem struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Amount uint   `json:"amount"`
}

type Courier struct {
	Id              int     `json:"id"`
	Active          bool    `json:"active"`
	InProgress      bool    `json:"in_progress"`
	Rating          float32 `json:"rating"`
	OrderDelivered int     `json:"order_delivered"`
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

type OrderStatus struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}
