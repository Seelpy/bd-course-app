package mysql

import (
	"embed"
	"github.com/pressly/goose/v3"
	"server/pkg/infrastructure/mysql"
)

const (
	dbDriverName  = "mysql"
	migrationsDir = "migrations"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func InitMigrations() {
	db, err := mysql.InitDBConnection()
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dbDriverName); err != nil {
		panic(err)
	}

	if err := goose.Up(db.DB, migrationsDir); err != nil {
		panic(err)
	}
}
