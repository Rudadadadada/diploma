package storage

import "diploma/services/courier/pkg/models"

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
		err = rows.Scan(&tmp.Id, &tmp.Active, &tmp.InProgress, &tmp.Rating, &tmp.Order_delivered)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return &tmp, nil
}