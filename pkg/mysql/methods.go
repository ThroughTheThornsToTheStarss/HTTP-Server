package mysql

import (
	"fmt"
	"os"

	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func NewGormFromEnv() (*gorm.DB, error) {
	cfg := LoadConfigFromEnv()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
