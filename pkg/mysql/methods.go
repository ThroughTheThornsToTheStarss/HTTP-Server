package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func LoadConfigFromEnv() Config {
	return Config{
		Host:     getenv("db_host", "fullstack-mysql"),
		Port:     getenv("db_port", "3306"),
		User:     getenv("db_user", "app_user"),
		Password: getenv("db_password", "app_password"),
		Name:     getenv("db_name", "amoCRM_http_server"),
	}
}

func ConnectFromEnv() (*DB, error) {
	cfg := LoadConfigFromEnv()
	return connectWithRetry(cfg)
}

func connectWithRetry(cfg Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	const (
		maxAttempts  = 10
		sleepBetween = time.Second
	)

	var (
		db  *sql.DB
		err error
	)

	for i := 1; i <= maxAttempts; i++ {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			fmt.Printf("mysql open attempt %d failed: %v\n", i, err)
		} else if err = db.Ping(); err != nil {
			fmt.Printf("mysql ping attempt %d failed: %v\n", i, err)
		} else {
			fmt.Printf("Connected to MySQL at %s:%s, db=%s\n", cfg.Host, cfg.Port, cfg.Name)
			db.SetMaxOpenConns(10)
			db.SetMaxIdleConns(5)
			db.SetConnMaxLifetime(time.Hour)
			return &DB{DB: db}, nil
		}

		time.Sleep(sleepBetween)
	}

	return nil, fmt.Errorf("mysql connect failed after %d attempts: %w", maxAttempts, err)
}

func (db *DB) Close() error {
	if db == nil || db.DB == nil {
		return nil
	}
	return db.DB.Close()
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
