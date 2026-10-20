package database

import (
	"database/sql"
	"fmt"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}