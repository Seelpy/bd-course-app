package transport

import (
	"errors"
	"fmt"
	"net/http"
	"server/api"
	domainmodel "server/pkg/domain/model"
	"server/pkg/domain/service"
	"server/pkg/infrastructure/model"
	"server/pkg/infrastructure/mysql/provider"
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
	bookChapterService service.BookChapterService,
	bookChapterTranslationService service.BookChapterTranslationService,
	readingSession service.ReadingSessionService,
	verifyBookRequestService service.VerifyBookRequestService,
	userQueryService query.UserQueryService,
	verifyBookRequestProvider provider.VerifyBookRequestProvider,
) api.ServerInterface {
	return &public{
		userService:                   userService,
		bookService:                   bookService,
		bookChapterService:            bookChapterService,
		bookChapterTranslationService: bookChapterTranslationService,
		readingSession:                readingSession,
		verifyBookRequestService:      verifyBookRequestService,

		userQueryService: userQueryService,

		verifyBookRequestProvider: verifyBookRequestProvider,
	}
}

type PublicAPI interface {
	api.ServerInterface
}

type public struct {
	userService                   service.UserService
	bookService                   service.BookService
	bookChapterService            service.BookChapterService
	bookChapterTranslationService service.BookChapterTranslationService
	readingSession                service.ReadingSessionService
	verifyBookRequestService      service.VerifyBookRequestService

	userQueryService query.UserQueryService

	verifyBookRequestProvider provider.VerifyBookRequestProvider
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

func (p public) EditUser(ctx echo.Context) error {
	var input api.EditUserRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.userService.EditUser(service.EditUserInput{
		ID:       domainmodel.UserID(input.Id),
		AboutMe:  input.AboutMe,
		Login:    input.Login,
		Password: input.Password,
	})
	if errors.Is(err, domainmodel.ErrUserNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to edit user: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("User edited successfully"),
	})
}

func (p public) DeleteUser(ctx echo.Context) error {
	var input api.DeleteUserRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.userService.DeleteUser(domainmodel.UserID(input.Id))
	if errors.Is(err, domainmodel.ErrUserNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	if errors.Is(err, domainmodel.ErrNotDeleteAdmin) {
		return echo.NewHTTPError(http.StatusForbidden, "Not allowed")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete user: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("User deleted successfully"),
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
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create book: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Book created successfully"),
	})
}

func (p public) EditBook(ctx echo.Context) error {
	var input api.EditBookRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.bookService.EditBook(service.EditBookInput{
		ID:          domainmodel.BookID(input.Id),
		Title:       input.Title,
		Description: input.Description,
	})
	if errors.Is(err, domainmodel.ErrBookNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Book not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to edit book: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Book edited successfully"),
	})
}

func (p public) DeleteBook(ctx echo.Context) error {
	var input api.DeleteBookRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.bookService.DeleteBook(domainmodel.BookID(input.Id))
	if errors.Is(err, domainmodel.ErrBookNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Book not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete book: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Book deleted successfully"),
	})
}

func (p public) CreateBookChapter(ctx echo.Context) error {
	var input api.CreateBookChapterRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.bookChapterService.CreateBookChapter(service.CreateBookChapterInput{
		BookID: domainmodel.BookID(input.BookId),
		Title:  input.Title,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create book chapter: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Book chapter created successfully"),
	})
}

func (p public) EditBookChapter(ctx echo.Context) error {
	var input api.EditBookChapterRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.bookChapterService.EditBookChapter(service.EditBookChapterInput{
		BookChapterID: domainmodel.BookChapterID(input.BookChapterId),
		Title:         input.Title,
	})
	if errors.Is(err, domainmodel.ErrBookChapterNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Book chapter not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create book chapter: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Book chapter edited successfully"),
	})
}

func (p public) DeleteBookChapter(ctx echo.Context) error {
	var input api.DeleteBookChapterRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.bookChapterService.DeleteBookChapter(domainmodel.BookChapterID(input.BookChapterId))
	if errors.Is(err, domainmodel.ErrBookChapterNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Book chapter not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete book chapter: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Book chapter deleted successfully"),
	})
}

func (p public) StoreBookChapterTranslation(ctx echo.Context) error {
	var input api.StoreBookChapterTranslationRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.bookChapterTranslationService.StoreBookChapterTranslation(service.StoreBookChapterTranslationInput{
		BookChapterID: domainmodel.BookChapterID(input.BookChapterId),
		TranslatorID:  domainmodel.UserID(input.TranslatorId),
		Text:          input.Text,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store book chapter translation: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Book chapter translation stored successfully"),
	})
}

func (p public) StoreReadingSession(ctx echo.Context) error {
	var input api.StoreReadingSessionRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.readingSession.StoreReadingSession(service.StoreReadingSessionInput{
		BookID:        domainmodel.BookID(input.BookId),
		BookChapterID: domainmodel.BookChapterID(input.BookChapterId),
		UserID:        domainmodel.UserID(input.UserId),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store reading session: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Reading session stored successfully"),
	})
}

func (p public) CreateVerifyBookRequest(ctx echo.Context) error {
	var input api.CreateVerifyBookRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.verifyBookRequestService.CreateVerifyBookRequest(service.CreateVerifyBookRequestInput{
		TranslatorID: domainmodel.UserID(input.TranslatorId),
		BookID:       domainmodel.BookID(input.BookId),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create verify book request: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Verify book request created successfully"),
	})
}

func (p public) DeleteVerifyBookRequest(ctx echo.Context) error {
	var input api.DeleteVerifyBookRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.verifyBookRequestService.DeleteVerifyBookRequest(domainmodel.VerifyBookRequestID(input.VerifyBookRequestId))
	if errors.Is(err, domainmodel.ErrVerifyBookRequestNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Verify book request not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete verify book request: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Verify book request deleted successfully"),
	})
}

func (p public) AcceptVerifyBookRequest(ctx echo.Context) error {
	var input api.AcceptBookChapterRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.verifyBookRequestService.AcceptVerifyBookRequest(service.AcceptVerifyBookRequestInput{
		VerifyBookRequestID: domainmodel.VerifyBookRequestID(input.VerifyBookRequestId),
		Accept:              input.Accept,
	})
	if errors.Is(err, domainmodel.ErrVerifyBookRequestNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Verify book request not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to accept verify book request: %s", err))
	}

	bookID, err := p.verifyBookRequestProvider.FindBookIDByVerifyBookRequestID(domainmodel.VerifyBookRequestID(input.VerifyBookRequestId))
	if errors.Is(err, domainmodel.ErrVerifyBookRequestNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Verify book request not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to find book to verify: %s", err))
	}

	err = p.bookService.PublishBook(service.PublishBookInput{
		ID:          bookID,
		IsPublished: input.Accept,
	})
	if errors.Is(err, domainmodel.ErrBookNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Book not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to find book to verify: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Verify book request verified successfully"),
	})
}

func ptr(s string) *string {
	return &s
}
