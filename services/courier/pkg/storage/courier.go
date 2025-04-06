package storage

import (
	"diploma/services/courier/pkg/models"
	"time"
)

func AddCourier(courierId int) error {
	_, err := db.Query(
		`insert into couriers (id, active, in_progress, rating, order_delivered)
		select * from (select cast($1 as integer), false, false, 5, 0) as tmp
		where not exists (
			select 1 from couriers where id = cast($1 as integer)
		)`, courierId,
	)

	if err != nil {
		return err
	}

	return nil
}

func SetActive(courierId int) error {
	_, err := db.Query(
		`update couriers
		set active = not active
		where id = $1;`, courierId,
	)

	if err != nil {
		return err
	}

	return nil
}

func SetInProgress(courierId int) error {
	_, err := db.Query(
		`update couriers
		set in_progress = not in_progress
		where id = $1;`, courierId,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetState(courierId int) (*models.Courier, error) {
	rows, err := db.Query(`select * from couriers where id = $1`, courierId)

	if err != nil {
		return nil, err
	}

	var tmp models.Courier
	for rows.Next() {
		err = rows.Scan(&tmp.Id, &tmp.Active, &tmp.InProgress, &tmp.Rating, &tmp.OrderDelivered)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return &tmp, nil
}

func FinishDelivery(courierId int, orderId int, deliveredAt time.Time) error {
	_, err := db.Query(`
		update orders 
		set delivered_at = $1, courier_id = -1
		where id = $2`, deliveredAt, orderId,
	)

	if err != nil {
		return err
	}

	var rating float32
	var orderDelivered int

	rows, err := db.Query(`select rating, order_delivered from couriers where id = $1`, courierId)
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&rating, &orderDelivered)
		if err != nil {
			return err
		}
	}

	if err = rows.Close(); err != nil {
		return err
	}

	var DeliveryStarted time.Time
	var DeliveredAt time.Time
	rows, err = db.Query(`select delivery_started, delivered_at from orders where id = $1`, orderId)
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&DeliveryStarted, &DeliveredAt)
		if err != nil {
			return err
		}
	}

	if err = rows.Close(); err != nil {
		return err
	}

	diff := DeliveredAt.Sub(DeliveryStarted)
	minutes := float32(diff.Minutes())

	var curOrderRating float32
	if minutes < 15 {
		curOrderRating = 5
	} else if 15 <= minutes && minutes < 45 {
		curOrderRating = 4
	} else {
		curOrderRating = 3
	}

	newOrderDeliveredCount := float32(orderDelivered) + 1
	newRating := (rating*float32(orderDelivered) + curOrderRating) / newOrderDeliveredCount

	_, err = db.Query(`
		update couriers
		set in_progress = false, rating = $1, order_delivered = $2
		where id = $3`, newRating, newOrderDeliveredCount, courierId,
	)

	if err != nil {
		return err
	}

	return nil
}

func DisableInProgress(courierId int) error {
	_, err := db.Query(`update couriers set in_progress = false where id = $1`, courierId)
	
	if err != nil {
		return err
	}

	return nil
}

func DeclinedByCourier(courierId int) error {
	var rating float32
	var orderDelivered int

	rows, err := db.Query(`select rating, order_delivered from couriers where id = $1`, courierId)
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&rating, &orderDelivered)
		if err != nil {
			return err
		}
	}

	curOrderRating := 2
	newOrderDeliveredCount := float32(orderDelivered) + 1
	newRating := (rating*float32(orderDelivered) + float32(curOrderRating)) / newOrderDeliveredCount

	_, err = db.Query(`
		update couriers
		set in_progress = false, rating = $1, order_delivered = $2
		where id = $3`, newRating, newOrderDeliveredCount, courierId,
	)

	if err != nil {
		return err
	}

	return nil
}