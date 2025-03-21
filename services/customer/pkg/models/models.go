package models

import "time"

type Product struct {
	Id         uint    `json:"id"`
	Name       string  `json:"name"`
	Amount     uint    `json:"amount"`
	Cost       float32 `json:"cost"`
	CategoryID uint    `json:"category_id"`
}

type Category struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type BucketItem struct {
	Id        uint    `json:"id"`
	Name      string  `json:"name"`
	Amount    uint    `json:"amount"`
	TotalCost float32 `json:"total_cost"`
}

type Order struct {
	Id          uint      `json:"id"`
	CustomerId  uint      `json:"customer_id"`
	BucketId    uint      `json:"bucket_id"`
	TotalCost   float64   `json:"total_cost"`
	CreatedAt   time.Time `json:"created_at"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	Status      string    `json:"status"`
}
