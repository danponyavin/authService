package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgreSQLStorage() *PostgresStorage {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5440"
	}
	connStr := fmt.Sprintf("host=%s port=%s user=postgres password=mypass dbname=test sslmode=disable", host, port)
	fmt.Println(connStr)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return &PostgresStorage{db: db}
}
