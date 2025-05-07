package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/rs/zerolog/log"
)

func ConnectDBCandyShop() *sqlx.DB {
	err := godotenv.Load()
	if err != nil {
		panic("Can't read .env file")
	}

	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	sslMode := os.Getenv("DATABASE_SSLMODE")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" || sslMode == "" {
		log.Panic().Err(err).Int("status", 500).Str("function", "db connection").Msg("failed to connect to database candy shop")
		panic("One or more environment variables are missing")
	}

	// Create the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslMode)

	// Open a connection to the database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Panic().Err(err).Int("status", 500).Str("function", "db connection").Msg("failed to connect to database candy shop")
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Panic().Err(err).Int("status", 500).Str("function", "db connection").Msg("failed to connect to database candy shop")
		panic(fmt.Sprintf("Failed to ping the database: %v", err))
	}

	fmt.Println("Successfully connected to the database candy shop!")
	return db
}
