package storage

import (
	"diploma/services/admin/pkg/models"
)

func SyncDatabases() (*models.SyncDatabasesMessage, error) {
	categories, err := ViewAllCategories()
	if err != nil {
		return nil, err
	}
	
	products, _, err := ViewAllProducts()
	if err != nil {
		return nil, err
	}
	
	err = SetCategory(products)
	if err != nil {
		return nil, err
	}

	sync := models.SyncDatabasesMessage{
		Categories: categories,
		Products: products,
	}

	return &sync, nil
}