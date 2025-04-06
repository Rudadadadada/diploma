package storage

import (
	"diploma/services/distribution/pkg/models"
)

func AddCourier(courier models.Courier) error {
	_, err := db.Query(
		`insert into couriers (id, active, in_progress, rating, order_delivered)
		select * from (select cast($1 as integer), cast($2 as bool), cast($3 as bool), cast($4 as integer), cast($5 as integer)) as tmp
		where not exists (
			select 1 from couriers where id = cast($1 as integer)
		)`, courier.Id, courier.Active, courier.InProgress, courier.Rating, courier.OrderDelivered,
	)

	if err != nil {
		return err
	}

	return nil
}

func SetState(courierId int, courier models.Courier) error {
	_, err := db.Query(
		`update couriers
		set active = $2, in_progress = $3
		where id = $1`, courierId, courier.Active, courier.InProgress,
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

func GetActiveCouriers() ([]models.Courier, error) {
	rows, err := db.Query(`select * from couriers where active = true and in_progress = false`)

	if err != nil {
		return nil, err
	}

	var couriers []models.Courier
	for rows.Next() {
		var tmp models.Courier
		err = rows.Scan(&tmp.Id, &tmp.Active, &tmp.InProgress, &tmp.Rating, &tmp.OrderDelivered)
		if err != nil {
			return nil, err
		}

		couriers = append(couriers, tmp)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return couriers, nil
}