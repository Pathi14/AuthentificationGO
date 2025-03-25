package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func ConnectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	pgConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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

func ConnectTestDB() (*sql.DB, error) {
	host := os.Getenv("TEST_DB_HOST")
	port := os.Getenv("TEST_DB_PORT")
	user := os.Getenv("TEST_DB_USER")
	password := os.Getenv("TEST_DB_PASSWORD")
	dbname := os.Getenv("TEST_DB_NAME")

	pgConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion à la DB de test : %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("échec de connexion à la DB de test : %w", err)
	}

	err = CreateTableIfNotExists(db)
	if err != nil {
		return nil, err
	}

	log.Println("✅ Connexion réussie à la DB de test")
	return db, nil
}
