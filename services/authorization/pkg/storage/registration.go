package storage

import (
	"diploma/services/authorization/pkg/models"
)

func CustomerRegistration(user models.Customer) error {
	_, err := db.Query(`insert into customers(name, surname, email, hash_password) values ($1, $2, $3, $4)`,
		user.Name, user.Surname, user.Email, user.HashPassword)

	if err != nil {
		return err
	}

	return nil
}

func CourierRegistration(user models.Courier) error {
	_, err := db.Query(`insert into couriers(name, surname, email, hash_password) values ($1, $2, $3, $4)`,
		user.Name, user.Surname, user.Email, user.HashPassword)

	if err != nil {
		return err
	}

	return nil
}