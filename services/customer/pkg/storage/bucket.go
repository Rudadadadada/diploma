package storage

import (
	"diploma/services/customer/pkg/models"
	"log"
)

func InsertIntoBucket(customerId int, items map[int]int) error {
	_, err := db.Query(
		`insert into bucket (customer_id, preparing)
		select * from (select cast($1 as integer), true) as tmp
		where not exists (
			select 1 from bucket where customer_id = cast($1 as integer) and preparing = true
		) 
		limit 1`, customerId,
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
	if err != nil {
		return nil, -1, err
	}

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

	if err = rows.Close(); err != nil {
		return nil, -1, err
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

func GetAllProductCost(bucketId int) (float64, error) {
	rows, err := db.Query(
		`select 
			products.cost * bucket_items.amount as total_cost
		from bucket_items as bucket_items
		left join products as products on bucket_items.product_id = products.id 
		where bucket_id = $1`, bucketId,
	)

	if err != nil {
		return -1, err
	}

	var allProductCost float64
	for rows.Next() {
		var tmp float64
		err = rows.Scan(&tmp)
		if err != nil {
			return -1, err
		}

		allProductCost += tmp
	}

	if err = rows.Close(); err != nil {
		return -1, err
	}

	return allProductCost, nil
}

func GetChangesAndUpdate(currentOrderItems []models.BucketItem, orderId int, newTotalCost int) (*bool, error) {
	var previousOrderItems []models.BucketItem

	for _, product := range currentOrderItems {
		rows, err := db.Query(`select id, amount from bucket_items where id = $1`, product.Id)

		if err != nil {
			log.Print(err)
			return nil, err
		}

		for rows.Next() {
			var tmp models.BucketItem
			err = rows.Scan(&tmp.Id, &tmp.Amount)
			if err != nil {
				return nil, err
			}

			previousOrderItems = append(previousOrderItems, tmp)
		}

		if err = rows.Close(); err != nil {
			return nil, err
		}
	}

	changed := false
	for i := 0; i < len(currentOrderItems); i++ {
		if currentOrderItems[i].Amount == previousOrderItems[i].Amount {
			continue
		}

		changed = true
		_, err := db.Query(`
			update bucket_items 
			set amount = $1 
			where id = $2`, currentOrderItems[i].Amount, currentOrderItems[i].Id,
		)

		if err != nil {
			return nil, err
		}
	}

	log.Print(newTotalCost)
	_, err := db.Query(`update orders set total_cost = $1 where id = $2`, newTotalCost, orderId)
	if err != nil {
		return nil, err
	}

	return &changed, nil
}

func CheckOrderIsEmpty(orderId int) (*bool, error) {
	rows, err := db.Query(`select total_cost from orders where id = $1`, orderId)
	if err != nil {
		return nil, err
	}

	var total_cost float32
	for rows.Next() {
		err = rows.Scan(&total_cost)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	isEmpty := false
	if total_cost == 0 {
		isEmpty = true
	}

	return &isEmpty, nil
}