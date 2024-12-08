package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"log"
	"server/api"
	"server/data/mysql"
	"server/pkg/domain/service"
	inframysql "server/pkg/infrastructure/mysql"
	"server/pkg/infrastructure/mysql/repo"
	"server/pkg/infrastructure/transport"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	e := echo.New()

	mysql.InitMigrations()

	db, err := inframysql.InitDBConnection()
	if err != nil {
		panic(err)
	}
	dependencyContainer := NewDependencyContainer(db)
	public := transport.NewPublicAPI(
		dependencyContainer.UserService(),
		dependencyContainer.BookService(),
	)

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

type DependencyContainer struct {
	userService service.UserService
	bookService service.BookService
}

func NewDependencyContainer(connection *sqlx.DB) *DependencyContainer {
	userRepository := repo.NewUserRepository(connection)
	userService := service.NewUserService(userRepository)

	bookRepository := repo.NewBookRepository(connection)
	bookService := service.NewBookService(bookRepository)

	return &DependencyContainer{
		userService: userService,
		bookService: bookService,
	}
}

func (container *DependencyContainer) UserService() service.UserService {
	return container.userService
}

func (container *DependencyContainer) BookService() service.BookService {
	return container.bookService
}
