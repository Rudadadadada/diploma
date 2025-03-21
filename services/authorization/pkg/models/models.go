package models

type Customer struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	HashPassword string `json:"hash_password"`
}

type Courier struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	HashPassword string `json:"hash_password"`
}

type Admin struct {
	Id           uint   `json:"id"`
	Admin        string `json:"admin"`
	HashPassword string `json:"hash_password"`
}
