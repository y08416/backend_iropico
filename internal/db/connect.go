package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQLドライバ
)

var DB *sql.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in environment")
	}

	// DB接続
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// 実際に接続確認
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL!")
}
