package storage

import (
	"diploma/services/admin/pkg/models"
)

func ViewAllProducts() ([]models.Product, error) {
	rows, err := db.Query(`select * from products`)

	if err != nil {
		return nil, err
	}

	var products []models.Product
	for rows.Next() {
		var tmp models.Product
		
		err = rows.Scan(&tmp.Id, &tmp.Name, &tmp.CategoryID)
		if err != nil {
			return nil, err
		}
		products = append(products, tmp)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return products, nil
}