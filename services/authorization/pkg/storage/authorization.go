package storage

import (
	"diploma/services/authorization/pkg/models"

)

func CustomerAuthorization(user models.Customer) (*models.Customer, error) {
	rows, err := db.Query(`select * from customers where email = $1`, user.Email)
	if err != nil {
		return nil, err
	}

	var tmp models.Customer
	for rows.Next() {
		err = rows.Scan(&tmp.Id, &tmp.Name, &tmp.Surname, &tmp.Email, &tmp.HashPassword)
		if err != nil {
			return nil, err
		}

	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return &tmp, nil
}

func CourierAuthorization(user models.Courier) (*models.Courier, error) {
	rows, err := db.Query(`select * from couriers where email = $1`, user.Email)

	if err != nil {
		return nil, err
	}

	var tmp models.Courier
	for rows.Next() {
		err = rows.Scan(&tmp.Id, &tmp.Name, &tmp.Surname, &tmp.Email, &tmp.HashPassword)
		if err != nil {
			return nil, err
		}

	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return &tmp, nil
}

func AdminAuthorization(admin models.Admin) (string, error) {
	rows, err := db.Query(`select * from admins where admin = $1`, admin.Admin)

	if err != nil {
		return "", err
	}

	var tmp models.Admin
	for rows.Next() {
		err = rows.Scan(&tmp.Id, &tmp.Admin, &tmp.HashPassword)
		if err != nil {
			return "", err
		}

	}

	if err = rows.Close(); err != nil {
		return "", err
	}

	return tmp.HashPassword, nil
}