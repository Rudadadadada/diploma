package models

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
