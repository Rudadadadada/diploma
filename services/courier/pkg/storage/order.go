package storage

import (
	"diploma/services/courier/pkg/models"
	"diploma/services/courier/pkg/utils"
	"time"
)

func InsertOrder(order models.OrderMessage) error {
	rows, err := db.Query(`
		insert into orders (id, customer_id, courier_id, total_cost, created_at, took, status) 
		values ($1, $2, -1, $3, $4, false, $5) returning id`, order.OrderId, order.CustomerId, order.TotalCost, order.CreatedAt, order.Status,
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
			created_at,
			status
		from orders where took = false and status != 'delivered' and status != 'order declined' order by id`,
	)

	if err != nil {
		return nil, err
	}

	var orders []models.Order
	for rows.Next() {
		var tmp models.Order
		err = rows.Scan(&tmp.Id, &tmp.TotalCost, &tmp.CreatedAt, &tmp.Status)
		if err != nil {
			return nil, err
		}

		tmp.CreatedAtStr = utils.TruncateTime(tmp.CreatedAt)

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

func CheckOrderTaken(orderId int) (int, *bool, error) {
	rows, err := db.Query(`select id, took from orders where id = $1`, orderId)
	if err != nil {
		return -1, nil, err
	}

	var id int
	var took bool
	for rows.Next() {
		err = rows.Scan(&id, &took)
		if err != nil {
			return -1, nil, err
		}
	}

	if err = rows.Close(); err != nil {
		return -1, nil, err
	}

	return id, &took, nil
}

func TakeOrder(orderId int, courierId int) error {
	_, err := db.Query(`update orders
		set courier_id = $1, took = true, delivery_started = $3
		where id = $2`, courierId, orderId, time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func GetOrderId(courierId int) (int, error) {
	rows, err := db.Query(`
		select
			id
		from orders where courier_id = $1`, courierId,
	)

	if err != nil {
		return -1, err
	}

	var orderId int
	for rows.Next() {
		err = rows.Scan(&orderId)
		if err != nil {
			return -1, err
		}
	}

	if err = rows.Close(); err != nil {
		return -1, err
	}

	return orderId, nil
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

func GetOrderCost(orderId int) (float32, error) {
	rows, err := db.Query(`select total_cost from orders where id = $1`, orderId)
	if err != nil {
		return -1, err
	}

	var totalCost float32
	for rows.Next() {
		err = rows.Scan(&totalCost)
		if err != nil {
			return -1, err
		}
	}

	if err = rows.Close(); err != nil {
		return -1, err
	}

	return totalCost, nil
}

func UpdateOrderStatus(orderId int, status string) error {
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

func GetOrderStatus(orderId int) (string, error) {
	rows, err := db.Query(`select status from orders where id = $1`, orderId)
	if err != nil {
		return "", err
	}

	var status string
	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			return "", err
		}
	}

	if err = rows.Close(); err != nil {
		return "", err
	}

	return status, nil
}

func DeclineOrder(orderId int) error {
	_, err := db.Query(`update orders set status = 'order declined', courier_id = -1 where id = $1`, orderId)

	if err != nil {
		return err
	}

	return nil
}

func GetChangesAndUpdate(currentOrderItems []models.BucketItem, orderId int, newTotalCost int) (*bool, error) {
	var previousOrderItems []models.BucketItem

	for _, product := range currentOrderItems {
		rows, err := db.Query(`select id, amount from order_items where id = $1`, product.Id)

		if err != nil {
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
			update order_items 
			set amount = $1 
			where id = $2`, currentOrderItems[i].Amount, currentOrderItems[i].Id,
		)

		if err != nil {
			return nil, err
		}
	}

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

	var totalCost float32
	for rows.Next() {
		err = rows.Scan(&totalCost)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	isEmpty := false
	if totalCost == 0 {
		isEmpty = true
	}

	return &isEmpty, nil
}