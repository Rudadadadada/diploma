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

type SyncDatabasesMessage struct {
	Categories []Category `json:"categories"`
	Products   []Product  `json:"products"`
}

type BucketItem struct {
	Id        uint    `json:"id"`
	ProductId uint    `json:"product_id"`
	Name      string  `json:"name"`
	Amount    uint    `json:"amount"`
	TotalCost float32 `json:"total_cost"`
}

type Order struct {
	Id          uint       `json:"id"`
	CustomerId  uint       `json:"customer_id"`
	BucketId    uint       `json:"bucket_id"`
	TotalCost   float32    `json:"total_cost"`
	CreatedAt   time.Time  `json:"created_at"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	Status      string     `json:"status"`
}

type Courier struct {
	Id uint `json:"id"`
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