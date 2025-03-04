package storage

import (
	"database/sql"
	"diploma/db_configs"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB


func New(db_config *db_configs.DBConfig) error {
	var err error
	const op = "storage.postgres.New"

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
		db_config.Host, db_config.Port, db_config.User, db_config.Password, db_config.DBname)
	fmt.Println(psqlconn)

	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func Ping() error {
	return db.Ping()
}

func Close() error {
	return db.Close()
}