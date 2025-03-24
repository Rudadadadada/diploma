package storage

import (
	"diploma/services/courier/pkg/models"
)

func InsertOrder(order models.OrderMessage) error {
	rows, err := db.Query(`
		insert into orders (id, customer_id, courier_id, total_cost, created_at, took) 
		values ($1, $2, -1, $3, $4, false) returning id`, order.OrderId, order.CustomerId, order.TotalCost, order.CreatedAt,
	)

	if err != nil {
		return err
	}

	var orderId int
	for rows.Next() {
		err = rows.Scan(&orderId)
		if err != nil {
			return err
		}
	}

	if err = rows.Close(); err != nil {
		return err
	}

	for _, product := range order.OrderItems {
		_, err := db.Query(
			`insert into order_items (id, order_id, product_id, name, amount, total_cost) values ($1, $2, $3, $4, $5, $6)`, 
			product.Id, order.OrderId, product.ProductId, product.Name, product.Amount, product.TotalCost,
		)

		if err != nil {
			return err
		}

	}

	return nil
}

func ViewOrders() ([]models.Order, error) {
	rows, err := db.Query(
		`select
			id,
			total_cost,
			created_at
		from orders where took = false`,
	)

	if err != nil {
		return nil, err
	}

	var orders []models.Order
	for rows.Next() {
		var tmp models.Order
		err = rows.Scan(&tmp.Id, &tmp.TotalCost, &tmp.CreatedAt)
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

func ViewOrderItem(orderId int) ([]models.OrderItem, error) {
	rows, err := db.Query(`select id, name, amount from order_items where order_id = $1`, orderId)

	if err != nil {
		return nil, err
	}

	var orderItems []models.OrderItem
	for rows.Next() {
		var tmp models.OrderItem
		err = rows.Scan(&tmp.Id, &tmp.Name, &tmp.Amount)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, tmp)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return orderItems, nil

}

func CheckOrderTaken(orderId int) (*bool, error) {
	rows, err := db.Query(`select took from orders where id = $1`, orderId)
	if err != nil {
		return nil, err
	}

	var took bool
	for rows.Next() {
		err = rows.Scan(&took)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return &took, nil
}

func TakeOrder(orderId int, courierId int) error {
	_, err := db.Query(`update orders
		set courier_id = $1, took = true
		where id = $2`, courierId, orderId,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetFullOrderInfo(orderId int) (*models.OrderMessage, error) {
	rows, err := db.Query(`
		select
			id,
			customer_id,
			total_cost,
			created_at
		from orders where id = $1`, orderId,
	)

	if err != nil {
		return nil, err
	}

	var order models.OrderMessage
	for rows.Next() {
		err = rows.Scan(&order.OrderId, &order.CustomerId, &order.TotalCost, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}
	
	rows, err = db.Query(`
		select
			id,
			product_id,
			name,
			amount,
			total_cost 
		from order_items where order_id = $1`, orderId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tmp models.BucketItem
		err = rows.Scan(&tmp.Id, &tmp.ProductId, &tmp.Name, &tmp.Amount, &tmp.TotalCost)
		if err != nil {
			return nil, err
		}

		order.OrderItems = append(order.OrderItems, tmp)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return &order, nil
}