package main

import (
	"log"
	"server/api"

	"github.com/labstack/echo/v4"

	"server/pkg/infrastructure/transport"
)

func main() {
	e := echo.New()

	public := transport.NewPublicAPI()

	api.RegisterHandlersWithBaseURL(e, public, "/")

	if err := e.Start(":8082"); err != nil {
		log.Fatal(err)
	}
}
