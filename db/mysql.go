package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("MYSQL_DSN environment variable is not set")
	}
	return sql.Open("mysql", dsn)
}
