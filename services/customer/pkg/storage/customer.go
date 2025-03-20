package storage

import (
	"diploma/services/customer/pkg/models"
)

func ViewAllCategories() ([]models.Category, error) {
	rows, err := db.Query(`select * from categories`)

	if err != nil {
		return nil, err
	}

	var categories []models.Category
	for rows.Next() {
		var tmp models.Category
		
		err = rows.Scan(&tmp.Id, &tmp.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, tmp)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return categories, nil
}

func ViewProductsByCategoryId(categoryID int) ([]models.Product, error) {
	rows, err := db.Query(`select * from products where category_id = $1`, categoryID)

	if err != nil {
		return nil, err
	}

	var products []models.Product
	for rows.Next() {
		var tmp models.Product
		
		err = rows.Scan(&tmp.Id, &tmp.Name, &tmp.Amount, &tmp.Cost, &tmp.CategoryID)
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