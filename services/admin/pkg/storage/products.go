package storage

import (
	"diploma/services/admin/pkg/models"
)

func ViewAllProducts() ([]models.Product, []string, error) {
	rows, err := db.Query(
		`select 
			products.id,
			products.name,
			products.amount,
			products.cost,
			categories.name 
		from products as products 
		left join 
		categories as categories on products.category_id = categories.id`)

	if err != nil {
		return nil, nil, err
	}

	var products []models.Product
	categoriesNames := []string{}

	for rows.Next() {
		var tmp models.Product
		var categoryName string

		err = rows.Scan(&tmp.Id, &tmp.Name, &tmp.Amount, &tmp.Cost, &categoryName)
		if err != nil {
			return nil, nil, err
		}
		products = append(products, tmp)
		categoriesNames = append(categoriesNames, categoryName)
	}

	if err = rows.Close(); err != nil {
		return nil, nil, err
	}

	return products, categoriesNames, nil
}

func CreateProduct(newProduct models.Product) error {
	_, err := db.Query(`insert into products(category_id, name, amount, cost) values ($1, $2, $3, $4)`, 
		newProduct.CategoryID, newProduct.Name, newProduct.Amount, newProduct.Cost)

	if err != nil {
		return err
	}

	return nil
}

func RemoveProduct(removeProduct models.Product) error {
	_, err := db.Query(`delete from products where id = ($1)`, removeProduct.Id)

	if err != nil {
		return err
	}

	return nil
}