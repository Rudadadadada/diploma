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
		categories as categories on products.category_id = categories.id order by products.id`)

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

func GetActualState(orderItems []models.BucketItem) ([]models.Product, error) {
	var actualProductsState []models.Product

	for _, product := range orderItems {
		rows, err := db.Query(`select id, amount, cost from products where id = $1`, product.ProductId)
		
		if err != nil {
			return nil, err
		}

		var tmp models.Product
		for rows.Next() {
			err = rows.Scan(&tmp.Id, &tmp.Amount, &tmp.Cost)
			if err != nil {
				return nil, err
			}
		}
	
		if err = rows.Close(); err != nil {
			return nil, err
		}

		actualProductsState = append(actualProductsState, tmp)
	}

	return actualProductsState, nil
}

func UpadteProducts(toUpdate []models.Product) error {
	for _, product := range toUpdate {
		_, err := db.Query(`
			update products
			set amount = $1
			where id = $2`, product.Amount, product.Id,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func SetCategory(products []models.Product) error {
	for i, product := range products {
		rows, err := db.Query(`select category_id from products where id = $1`, product.Id)
		if err != nil {
			return err 
		}

		for rows.Next() {
			err = rows.Scan(&products[i].CategoryID)
			if err != nil {
				return err
			}
		}
	
		if err = rows.Close(); err != nil {
			return  err
		}
	}

	return nil
}