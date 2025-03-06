package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "secret"
	dbname   = "authentificationgo"
)

func ConnectDB() (*sql.DB, error) {
	pgConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	err = CreateTableIfNotExists(conn)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the PostgreSQL database")
	return conn, nil
}

func CreateTableIfNotExists(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		age INT,
		mobile_number VARCHAR(20),
		email VARCHAR(100) UNIQUE,
		password VARCHAR(255)
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return fmt.Errorf("failed to create users table: %w", err)
	}

	log.Println("Table 'users' is ready.")
	return nil
}
