package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	return sql.Open("mysql", "root:viraj3rathod@tcp(localhost:3306)/test?parseTime=true")
}
