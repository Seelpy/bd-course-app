package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"server/api"
	"server/data/mysql"
	"server/pkg/domain/service"
	"server/pkg/infrastructure/model"
	inframysql "server/pkg/infrastructure/mysql"
	"server/pkg/infrastructure/mysql/query"
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
		dependencyContainer.UserQueryService(),
	)

	api.RegisterHandlersWithBaseURL(e, public, "")

	e.POST("/login", public.Login)
	e.POST("/api/v1/book/create", public.CreateBook, MiddlewareRole(0))

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

	userQueryService query.UserQueryService
}

func NewDependencyContainer(connection *sqlx.DB) *DependencyContainer {
	userRepository := repo.NewUserRepository(connection)
	userService := service.NewUserService(userRepository)

	bookRepository := repo.NewBookRepository(connection)
	bookService := service.NewBookService(bookRepository)

	userQueryService := query.NewUserQueryService(connection)

	return &DependencyContainer{
		userService: userService,
		bookService: bookService,

		userQueryService: userQueryService,
	}
}

func (container *DependencyContainer) UserService() service.UserService {
	return container.userService
}

func (container *DependencyContainer) BookService() service.BookService {
	return container.bookService
}

func (container *DependencyContainer) UserQueryService() query.UserQueryService {
	return container.userQueryService
}

func MiddlewareRole(requiredRole int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
			}

			claims := &model.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return transport.JwtKey, nil
			})

			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			if claims.Role != requiredRole {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
			}

			c.Set("login", claims.Login)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}
