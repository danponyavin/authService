package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgreSQLStorage() *PostgresStorage {
	connStr := "host=localhost port=5440 user=postgres password=mypass dbname=test sslmode=disable"

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return &PostgresStorage{db: db}
}
