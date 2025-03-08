package storage

import (
	"database/sql"
	"diploma/postgres_configs"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB
var PostgresCfg postgres_configs.DBConfig


func New() error {
	var err error
	const op = "storage.postgres.New"

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
		PostgresCfg.Host, PostgresCfg.Port, PostgresCfg.User, PostgresCfg.Password, PostgresCfg.DBname)

	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}