package storage

func AddCourier(courierId int) error {
	_, err := db.Query(
		`insert into activity (id, active)
		select * from (select cast($1 as integer), false) as tmp
		where not exists (
			select 1 from activity where id = cast($1 as integer)
		)`, courierId,
	)

	if err != nil {
		return err
	}

	return nil
}

func SetActive(courierId int) error {
	_, err := db.Query(
		`update activity
		set active = not active
		where id = $1;`, courierId,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetState(courierId int) (*bool, error) {
	rows, err := db.Query(`select active from activity where id = $1`, courierId)

	if err != nil {
		return nil, err
	}

	var active bool
	for rows.Next() {
		err = rows.Scan(&active)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return &active, nil
}
