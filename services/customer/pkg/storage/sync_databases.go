package storage

import (
	"diploma/services/customer/pkg/models"
)

func SyncDatabases(msg models.SyncDatabasesMessage) error {
	_, err := db.Exec(`delete from categories`)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(`delete from products`)
	if err != nil {
		return err
	}

	for _, category := range msg.Categories {
		_, err := db.Exec(
			`insert into categories (id, name)
			select * from (select cast($1 as integer), $2) as tmp
			where not exists (
				select 1 from categories where id = cast($1 as integer)
			)`, category.Id, category.Name,
		)

		if err != nil {
			return err
		}
	}

	for _, prodcut := range msg.Products {
		_, err := db.Exec(
			`insert into products (id, category_id, name, amount, cost)
			select * from (select cast($1 as integer), cast($2 as integer), $3, cast($4 as integer), cast($5 as numeric(10, 2))) as tmp
			where not exists (
				select 1 from products where id = cast($1 as integer)
			)`, prodcut.Id, prodcut.CategoryID, prodcut.Name, prodcut.Amount, prodcut.Cost,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
