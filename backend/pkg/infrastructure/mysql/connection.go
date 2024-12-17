package mysql

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	dbDriverName = "mysql"
	dsn          = "user:userpassword@tcp(db_db:3306)/mydatabase?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true"
)

func InitDBConnection() (*sqlx.DB, error) {
	dbx, err := sqlx.Open(dbDriverName, dsn)
	if err != nil {
		fmt.Printf("open db err: %v\n", err)
		return nil, err
	}

	waitForDB(dbx)

	return dbx, nil
}

func waitForDB(db *sqlx.DB) {
	log.Println("Pinging DB...")

	for {
		if err := db.Ping(); err == nil {
			log.Println("DB is reachable")
			break
		} else {
			// log.Printf("DB ping failed: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}
