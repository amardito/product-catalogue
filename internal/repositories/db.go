package repositories

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDBConnection() (*sql.DB, error) {
	host := "localhost"
	port := 5432
	dbname := "product_catalog"
	user := "postgres"
	password := "postgres"

	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, user, password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Set connection pool settings (optional)
	db.SetMaxOpenConns(10) // Maximum number of open connections
	db.SetMaxIdleConns(5)  // Maximum number of idle connections

	return db, nil
}
