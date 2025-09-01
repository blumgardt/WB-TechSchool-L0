package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToDb() (*sqlx.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "orders_user")
	password := os.Getenv("DB_PASSWORD")
	dbname := getEnv("DB_NAME", "orders_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	var dsn string
	if password == "" {
		dsn = fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=%s", user, host, port, dbname, sslmode)
	} else {
		dsn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
	}

	log.Printf("Connecting to DB: %s@%s:%s/%s (sslmode=%s)\n", user, host, port, dbname, sslmode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅ Successfully connected to PostgreSQL")
	return db, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
