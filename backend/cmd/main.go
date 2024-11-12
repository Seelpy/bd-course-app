package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"server/api"
	"server/data/mysql"
	"server/pkg/infrastructure/transport"
)

func main() {
	e := echo.New()

	mysql.InitMigrations()

	public := transport.NewPublicAPI()

	api.RegisterHandlersWithBaseURL(e, public, "/")

	if err := e.Start(":8082"); err != nil {
		log.Fatal(err)
	}
}
