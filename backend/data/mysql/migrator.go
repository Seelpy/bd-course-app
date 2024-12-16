package mysql

import (
	"embed"
	"log"
	"server/pkg/infrastructure/mysql"

	"github.com/pressly/goose/v3"
)

const (
	dbDriverName  = "mysql"
	migrationsDir = "migrations"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func InitMigrations() {
	log.Println("Initializing DB connection for migrations")
	db, err := mysql.InitDBConnection()
	if err != nil {
		log.Fatalf("Failed to initialize DB connection: %v", err)
	}

	log.Println("Setting base FS for goose")
	goose.SetBaseFS(embedMigrations)

	log.Println("Setting dialect for goose")
	if err := goose.SetDialect(dbDriverName); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	log.Println("Running migrations")
	if err := goose.Up(db.DB, migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}
