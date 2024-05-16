package initializers

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)
const createTableQuery = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    age INTEGER NOT NULL,
    phone_no INTEGER NOT NULL,
	secret_code TEXT NOT NULL,
    role_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`;


func ConnectToDB() (*pgx.Conn, error) {
	var err error

	connString := os.Getenv("DB_URI")
	db, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	//defer db.Close(context.Background())

	// Create the users table if it doesn't exist
	_, err = db.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	return db, nil

}
