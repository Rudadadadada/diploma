package storage

import (
	"diploma/services/distribution/pkg/models"
)

func AddCourier(courierId int) error {
	_, err := db.Query(
		`insert into active_couriers (id, active)
		select * from (select cast($1 as integer), true) as tmp
		where not exists (
			select 1 from active_couriers where id = cast($1 as integer)
		)`, courierId,
	)

	if err != nil {
		return err
	}

	return nil
}

func SetState(courierId int, state bool) error {
	_, err := db.Query(
		`update active_couriers
		set active = $2
		where id = $1`, courierId, state,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetActiveCouriers() ([]models.CourierState, error) {
	rows, err := db.Query(`select * from active_couriers where active = true`)

	if err != nil {
		return nil, err
	}

	var courierStates []models.CourierState
	for rows.Next() {
		var tmp models.CourierState
		err = rows.Scan(&tmp.CourierId, &tmp.State)
		if err != nil {
			return nil, err
		}

		courierStates = append(courierStates, tmp)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return courierStates, nil
}