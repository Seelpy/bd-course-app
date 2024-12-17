package transport

import (
	"fmt"
	"net/http"
	"server/api"
	"server/pkg/domain/service"
	"server/pkg/infrastructure/model"
	"server/pkg/infrastructure/mysql/query"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

var JwtKey = []byte("your_secret_key")

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewPublicAPI(
	userService service.UserService,
	bookService service.BookService,
	userQueryService query.UserQueryService,
) api.ServerInterface {
	return &public{
		userService:      userService,
		bookService:      bookService,
		userQueryService: userQueryService,
	}
}

type PublicAPI interface {
	api.ServerInterface
}

type public struct {
	userService service.UserService
	bookService service.BookService

	userQueryService query.UserQueryService
}

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (p public) LoginUser(ctx echo.Context) error {
	var userReq LoginUserRequest
	if err := ctx.Bind(&userReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	user, err := p.userQueryService.FindByLogin(userReq.Login)
	if err != nil {
		return err
	}
	if user == (model.User{}) || user.Password != userReq.Password {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &model.Claims{
		Login: user.Login,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create token")
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = tokenString
	cookie.Expires = expirationTime
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode

	ctx.SetCookie(cookie)
	return ctx.NoContent(http.StatusOK)
}

func (p public) ListUsers(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, User{ID: "1", Name: "Igor", Email: "maks@mail.ru"})
}

func (p public) CreateUser(ctx echo.Context) error {
	var input api.CreateUserRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.userService.CreateUser(service.CreateUserInput{
		AboutMe:  input.AboutMe,
		Login:    input.Login,
		Password: input.Password,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("User created successfully"),
	})
}

func (p public) GetUser(ctx echo.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (p public) CreateBook(ctx echo.Context) error {
	var input api.CreateBookRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.bookService.CreateBook(service.CreateBookInput{
		Title:       input.Title,
		Description: input.Description,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Book created successfully"),
	})
}

func ptr(s string) *string {
	return &s
}
