package storage

import (
	"diploma/services/customer/pkg/models"
	"time"
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
		
		err = rows.Scan(&tmp.Id, &tmp.CategoryID, &tmp.Name, &tmp.Amount, &tmp.Cost)
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

func MakeOrder(bucketId int, customerId int, allProductsCost float64) error {
	_, err := db.Query(`insert into orders (bucket_id, customer_id, total_cost, status) values ($1, $2, $3, 'created')`, 
		bucketId, customerId, allProductsCost)

	if err != nil {
		return err
	}

	return nil
}

func SelectOrderIdAndCreatedAt(bucketId int, customerId int) (int, time.Time, error) {
	rows, err := db.Query(`select id, created_at from orders where bucket_id = $1 and customer_id = $2`, bucketId, customerId)

	if err != nil {
		return -1, time.Time{}, err
	}

	var id int
	var created_at time.Time
	for rows.Next() {
		err = rows.Scan(&id, &created_at)
		if err != nil {
			return -1, time.Time{}, err
		}
	}
	
	if err = rows.Close(); err != nil {
		return -1, time.Time{}, err
	}

	return id, created_at, nil
}

func ViewOrders(customerId int) ([]models.Order, error) {
	rows, err := db.Query(`select * from orders where customer_id = $1 order by id`, customerId)
	if err != nil {
		return nil, err
	}

	var orders []models.Order
	for rows.Next() {
		var tmp models.Order
		
		err = rows.Scan(&tmp.Id, &tmp.CustomerId, &tmp.BucketId, &tmp.TotalCost, &tmp.CreatedAt, &tmp.DeliveredAt, &tmp.Status)
		if err != nil {
			return nil, err
		}

		orders = append(orders, tmp)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return orders, nil
}

func ViewOrderItems(bucketId int) ([]models.BucketItem, error) {
	rows, err := db.Query(
		`select 
			bucket_items.id,
			products.id,
			products.name,
			bucket_items.amount,
			products.cost * bucket_items.amount as total_cost
		from bucket_items as bucket_items 
		left join 
		products as products on bucket_items.product_id = products.id
		where bucket_id = $1`, bucketId,
	)
	if err != nil {
		return nil, err
	}

	var bucketItems []models.BucketItem
	for rows.Next() {
		var tmp models.BucketItem
		
		err = rows.Scan(&tmp.Id, &tmp.ProductId, &tmp.Name, &tmp.Amount, &tmp.TotalCost)
		if err != nil {
			return nil, err
		}

		bucketItems = append(bucketItems, tmp)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return bucketItems, nil
}

func UpdateStatus(orderId int, status string) error {
	_, err := db.Query(
		`update orders
		set status = $2
		where id = $1`, orderId, status,
	)
	
	if err != nil {
		return err
	}

	return nil
}

func GetOrderStatuses(customerId int) ([]models.OrderStatus, error) {
	rows, err := db.Query(`select id, status from orders where customer_id = $1`, customerId)
	if err != nil {
		return nil, err
	}

	var statuses []models.OrderStatus
	for rows.Next() {
		var tmp models.OrderStatus
		err = rows.Scan(&tmp.Id, &tmp.Status)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, tmp)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}


	return statuses, nil
}