package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}
