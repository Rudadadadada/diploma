package storage

import (
	"diploma/services/admin/pkg/models"
)

func ViewAllCategories() ([]models.Category, error) {
	rows, err := db.Query(`select * from categories order by id`)

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

func CreateCategory(newCategory models.Category) error {
	_, err := db.Query(`insert into categories(name) values ($1)`, newCategory.Name)

	if err != nil {
		return err
	}

	return nil
}

func RemoveCategory(removeCategory models.Category) error {
	_, err := db.Query(`delete from categories where id = ($1)`, removeCategory.Id)

	if err != nil {
		return err
	}

	return nil
}