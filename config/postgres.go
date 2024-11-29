package config

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ghulammuzz/backend-parkerin/pkg/log"
	_ "github.com/lib/pq"
)

var (
	dbInstance *sql.DB
	once       sync.Once
	initErr    error
)

func InitPostgres() (*sql.DB, error) {
	once.Do(func() {
		host := os.Getenv("PG_HOST")
		port := os.Getenv("PG_PORT")
		user := os.Getenv("PG_USER")
		password := os.Getenv("PG_PASS")
		dbname := os.Getenv("PG_NAME")

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			initErr = fmt.Errorf("failed to open database connection: %w", err)
			return
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err := db.Ping(); err != nil {
			initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		log.Info("Database connected successfully")
		dbInstance = db
	})

	return dbInstance, initErr
}
