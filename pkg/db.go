package pkg

import (
	"car-backend/pkg/config"
	"database/sql"
	"log"
)

// DB is a database connection variable.
var DB *sql.DB

// InitDB establishes a connection to the database
func InitDB() error {
	config := config.Load()

	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		log.Println("Failed to open DB Connection")
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	DB = db
	log.Println("Successfully connected to the database!")
	return nil
}
