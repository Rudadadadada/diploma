package storage

import (
	"diploma/services/customer/pkg/models"
)

func InertIntoBucket(customerId int, items map[int]int) error {
	_, err := db.Query(
		`insert into bucket (customer_id, preparing)
		select * from (select cast($1 as integer), true) as tmp
		where not exists (
			select 1 from bucket where customer_id = cast($1 as integer) and preparing = true
		) 
		limit 1
		returning id`, customerId,
	)

	if err != nil {
		return err
	}

	rows, err := db.Query(`select id from bucket where customer_id = $1 and preparing = true`, customerId)

	if err != nil {
		return err
	}

	var bucketId int
	for rows.Next() {
		err = rows.Scan(&bucketId)
		if err != nil {
			return err
		}
	}

	if err = rows.Close(); err != nil {
		return err
	}

	for productId, amount := range items {
		_, err := db.Query(
			`insert into bucket_items (bucket_id, product_id, amount) values ($1, $2, $3)
			on conflict (bucket_id, product_id)
			do update set amount = bucket_items.amount + $3`,
			bucketId, productId, amount)
		if err != nil {
			return err
		}
	}

	return nil
}

func ViewBucket(customerId int) ([]models.BucketItem, int, error) {
	rows, err := db.Query(
		`select 
			bucket_items.product_id,
			bucket_items.bucket_id,
			products.name,
			bucket_items.amount,
			products.cost * bucket_items.amount as total_cost
		from bucket_items as bucket_items 
		left join 
		bucket as bucket on bucket.id = bucket_items.bucket_id
		left join
		products as products on bucket_items.product_id = products.id
		where customer_id = $1 and preparing = true`, customerId,
	)

	bucketItems := []models.BucketItem{}
	var bucketId int
	for rows.Next() {
		var tmp models.BucketItem
		err = rows.Scan(&tmp.Id, &bucketId, &tmp.Name, &tmp.Amount, &tmp.TotalCost)
		if err != nil {
			return nil, -1, err
		}
		bucketItems = append(bucketItems, tmp)
	}
	return bucketItems, bucketId, nil
}

func RemoveItemFromBucket(bucketId int, productId int) error {
	_, err := db.Query(`delete from bucket_items where bucket_id = cast($1 as integer) and product_id = cast($2 as integer)`, bucketId, productId)

	if err != nil {
		return err
	}

	return nil
}

func UpdateBucketItems(bucketId int, items map[int]int) error {
	for productId, amount := range items {
		_, err := db.Query(
			`update bucket_items
			set amount = $1
			where bucket_id = $2 and product_id = $3;`, amount, bucketId, productId,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateBucketStatus(bucketId int, customerId int) error {
	_, err := db.Query(`update bucket set preparing = false where id = $1 and customer_id = $2`, bucketId, customerId)

	if err != nil {
		return err
	}

	return nil
}