package main

import (
	"log"
	"server/api"
	"server/data/mysql"
	"server/pkg/infrastructure/transport"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	e := echo.New()

	mysql.InitMigrations()

	public := transport.NewPublicAPI()

	api.RegisterHandlersWithBaseURL(e, public, "")

	e.File("/api/v1/openapi.yaml", "./api/api.yaml")

	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URLs = []string{"/api/v1/openapi.yaml"}
		c.InstanceName = "custom"
		c.DocExpansion = "list"
		c.DeepLinking = true
	}))

	if err := e.Start(":8082"); err != nil {
		log.Fatal(err)
	}
}
