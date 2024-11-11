package transport

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"server/api"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewPublicAPI() api.ServerInterface {
	return &public{}
}

type PublicAPI struct {
	api.ServerInterface
}

type public struct {
}

func (p public) ListUsers(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, User{ID: "1", Name: "Igor", Email: "maks@mail.ru"})
}

func (p public) GetUser(ctx echo.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
