package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sqlx.DB
}

func Connect(dbUsername, dbName, password string) (*Database, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", "postgres", "saurav", "agri")
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &Database{
		DB: conn,
	}, nil
}
