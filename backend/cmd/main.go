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
	"server/pkg/infrastructure/mysql/provider"
	"server/pkg/infrastructure/mysql/query"
	"server/pkg/infrastructure/mysql/repo"
	"server/pkg/infrastructure/transport"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	log.Println("Starting server...")
	e := echo.New()

	log.Println("Initiating migrations")
	mysql.InitMigrations()

	log.Println("Initiating DB connection")
	db, err := inframysql.InitDBConnection()
	if err != nil {
		panic(err)
	}

	log.Println("Creating dependency container")
	dependencyContainer := NewDependencyContainer(db)

	log.Println("Loading API")
	public := transport.NewPublicAPI(
		dependencyContainer.UserService(),
		dependencyContainer.BookService(),
		dependencyContainer.BookChapterService(),
		dependencyContainer.BookChapterTranslationService(),
		dependencyContainer.ReadingSessionService(),
		dependencyContainer.VerifyBookRequestService(),
		dependencyContainer.UserQueryService(),
		dependencyContainer.BookQueryService(),
		dependencyContainer.BookChapterQueryService(),
		dependencyContainer.BookChapterTranslationQueryService(),
		dependencyContainer.VerifyBookRequestQueryService(),
		dependencyContainer.VerifyBookRequestProvider(),
		dependencyContainer.BookRatingService(),
	)

	log.Println("Creating endpoints")
	api.RegisterHandlersWithBaseURL(e, public, "")

	e.POST("/api/v1/verify-book-request/accept", public.AcceptVerifyBookRequest, MiddlewareRole(0))

	e.File("/api/v1/openapi.yaml", "./api/api.yaml")

	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URLs = []string{"/api/v1/openapi.yaml"}
		c.InstanceName = "custom"
		c.DocExpansion = "list"
		c.DeepLinking = true
	}))

	log.Println("Starting listening...")
	if err := e.Start(":8082"); err != nil {
		log.Fatal(err)
	}
}

type DependencyContainer struct {
	userService                   service.UserService
	bookService                   service.BookService
	bookChapterService            service.BookChapterService
	bookChapterTranslationService service.BookChapterTranslationService
	readingSessionService         service.ReadingSessionService
	verifyBookRequestService      service.VerifyBookRequestService
	bookRatingService             service.BookRatingService

	userQueryService                   query.UserQueryService
	bookQueryService                   query.BookQueryService
	bookChapterQueryService            query.BookChapterQueryService
	bookChapterTranslationQueryService query.BookChapterTranslationQueryService
	verifyBookRequestQueryService      query.VerifyBookRequestQueryService

	verifyBookRequestProvider provider.VerifyBookRequestProvider
}

func NewDependencyContainer(connection *sqlx.DB) *DependencyContainer {
	userRepository := repo.NewUserRepository(connection)
	userService := service.NewUserService(userRepository)

	bookRepository := repo.NewBookRepository(connection)
	bookService := service.NewBookService(bookRepository)

	bookChapterRepository := repo.NewBookChapterRepository(connection)
	bookChapterService := service.NewBookChapterService(bookChapterRepository)

	bookChapterTranslationRepository := repo.NewBookChapterTranslationRepository(connection)
	bookChapterTranslationService := service.NewBookChapterTranslationService(bookChapterTranslationRepository)

	readingSessionRepository := repo.NewReadingSessionRepository(connection)
	readingSessionService := service.NewReadingSessionService(readingSessionRepository)

	bookRatingRepository := repo.NewBookRatingRepository(connection)
	bookRatingService := service.NewBookRatingService(bookRatingRepository)

	verifyBookRequestRepository := repo.NewVerifyBookRequestRepository(connection)
	verifyBookRequestService := service.NewVerifyBookRequestService(verifyBookRequestRepository)

	userQueryService := query.NewUserQueryService(connection)

	bookQueryService := query.NewBookQueryService(connection)

	bookChapterQueryService := query.NewBookChapterQueryService(connection)

	bookChapterTranslationQueryService := query.NewBookChapterTranslationQueryService(connection)

	verifyBookRequestQueryService := query.NewVerifyBookRequestQueryService(connection)

	verifyBookRequestProvider := provider.NewVerifyBookRequestProvider(connection)

	return &DependencyContainer{
		userService:                   userService,
		bookService:                   bookService,
		bookChapterService:            bookChapterService,
		bookChapterTranslationService: bookChapterTranslationService,
		readingSessionService:         readingSessionService,
		verifyBookRequestService:      verifyBookRequestService,
		bookRatingService:             bookRatingService,

		userQueryService:                   userQueryService,
		bookQueryService:                   bookQueryService,
		bookChapterQueryService:            bookChapterQueryService,
		bookChapterTranslationQueryService: bookChapterTranslationQueryService,
		verifyBookRequestQueryService:      verifyBookRequestQueryService,

		verifyBookRequestProvider: verifyBookRequestProvider,
	}
}

func (container *DependencyContainer) UserService() service.UserService {
	return container.userService
}

func (container *DependencyContainer) BookService() service.BookService {
	return container.bookService
}

func (container *DependencyContainer) BookRatingService() service.BookRatingService {
	return container.bookRatingService
}

func (container *DependencyContainer) BookChapterService() service.BookChapterService {
	return container.bookChapterService
}

func (container *DependencyContainer) BookChapterTranslationService() service.BookChapterTranslationService {
	return container.bookChapterTranslationService
}

func (container *DependencyContainer) ReadingSessionService() service.ReadingSessionService {
	return container.readingSessionService
}

func (container *DependencyContainer) VerifyBookRequestService() service.VerifyBookRequestService {
	return container.verifyBookRequestService
}

func (container *DependencyContainer) UserQueryService() query.UserQueryService {
	return container.userQueryService
}

func (container *DependencyContainer) BookQueryService() query.BookQueryService {
	return container.bookQueryService
}

func (container *DependencyContainer) BookChapterQueryService() query.BookChapterQueryService {
	return container.bookChapterQueryService
}

func (container *DependencyContainer) BookChapterTranslationQueryService() query.BookChapterTranslationQueryService {
	return container.bookChapterTranslationQueryService
}

func (container *DependencyContainer) VerifyBookRequestQueryService() query.VerifyBookRequestQueryService {
	return container.verifyBookRequestQueryService
}

func (container *DependencyContainer) VerifyBookRequestProvider() provider.VerifyBookRequestProvider {
	return container.verifyBookRequestProvider
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
