package mysql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	dbDriverName = "mysql"
)

func InitDBConnection() (*sqlx.DB, error) {
	dbx, err := sqlx.Open(dbDriverName, "user:userpassword@tcp(127.0.0.1:3306)/mydatabase?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
	if err != nil {
		fmt.Printf("open db err: %v\n", err)
		return nil, err
	}

	if err := dbx.Ping(); err != nil {
		fmt.Printf("open db err: %v\n", err)
		return nil, err
	}

	return dbx, nil
}
