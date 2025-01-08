package transport

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"math"
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
	bookQueryService query.BookQueryService,
	bookChapterQueryService query.BookChapterQueryService,
	bookChapterTranslationQueryService query.BookChapterTranslationQueryService,
	verifyBookRequestProvider provider.VerifyBookRequestProvider,
	bookRatingService service.BookRatingService,
) api.ServerInterface {
	return &public{
		userService:                   userService,
		bookService:                   bookService,
		bookChapterService:            bookChapterService,
		bookChapterTranslationService: bookChapterTranslationService,
		readingSession:                readingSession,
		verifyBookRequestService:      verifyBookRequestService,
		bookRatingService:             bookRatingService,

		userQueryService:                   userQueryService,
		bookQueryService:                   bookQueryService,
		bookChapterQueryService:            bookChapterQueryService,
		bookChapterTranslationQueryService: bookChapterTranslationQueryService,

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
	bookRatingService             service.BookRatingService

	userQueryService                   query.UserQueryService
	bookQueryService                   query.BookQueryService
	bookChapterQueryService            query.BookChapterQueryService
	bookChapterTranslationQueryService query.BookChapterTranslationQueryService

	verifyBookRequestProvider provider.VerifyBookRequestProvider
}

func (p public) UpdateBookRating(ctx echo.Context, id string) error {
	var req api.UpdateBookRatingRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	userID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	var bookID uuid.UUID
	err = bookID.Parse(id)
	if err != nil {
		return err
	}

	return p.bookRatingService.StoreRating(service.StoreBookRatingInput{BookID: bookID, UserID: userID, Value: req.Value})
}

func (p public) DeleteBookRating(ctx echo.Context, id string) error {
	userID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	var bookID uuid.UUID
	err = bookID.Parse(id)
	if err != nil {
		return err
	}

	return p.bookRatingService.DeleteRating(bookID, userID)
}

func (p public) GetBookRating(ctx echo.Context, id string) error {
	var bookID uuid.UUID
	err := bookID.Parse(id)
	if err != nil {
		return err
	}

	stat, err := p.bookRatingService.GetStatistics(bookID)

	return ctx.JSON(http.StatusCreated, api.GetBookRatingResponse{
		Average: ptr(stat.Average),
		Count:   ptr(stat.Count),
	})
}

func (p public) LoginUser(ctx echo.Context) error {
	var userReq api.LoginUserRequest
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
		Login:  user.Login,
		Role:   user.Role,
		UserID: user.ID.String(),
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

func (p public) ListBook(ctx echo.Context, page int, size int) error {
	bookOutputs, err := p.bookQueryService.List(page, size)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete book: %s", err))
	}

	booksRespData := make([]api.Book, len(bookOutputs))
	for i, b := range bookOutputs {
		cover, ok := b.Cover.Get()

		booksRespData[i] = api.Book{
			BookId:      openapi_types.UUID(b.BookID),
			Cover:       ptr(cover),
			Title:       b.Title,
			Description: b.Description,
		}

		if !ok {
			booksRespData[i].Cover = nil
		}
	}

	countBook, err := p.bookQueryService.CountBook(true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list book: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.ListBookResponse{
		Books:      booksRespData,
		CountPages: ptr(int(math.Ceil(float64(countBook) / float64(size)))),
	})
}

func (p public) GetBook(ctx echo.Context, id string) error {
	var bookID uuid.UUID
	err := bookID.Parse(id)
	if err != nil {
		return err
	}

	book, err := p.bookQueryService.FindByID(domainmodel.BookID(bookID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list book: %s", err))
	}

	cover, ok := book.Cover.Get()

	bookRespData := api.Book{
		BookId:      openapi_types.UUID(book.BookID),
		Cover:       ptr(cover),
		Title:       book.Title,
		Description: book.Description,
	}

	if !ok {
		bookRespData.Cover = nil
	}

	return ctx.JSON(http.StatusCreated, api.GetBookResponse{
		Book: bookRespData,
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

func (p public) ListBookChapter(ctx echo.Context) error {
	var input api.ListBookChapterRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	bookChaptersOutput, err := p.bookChapterQueryService.ListByBookID(domainmodel.BookID(input.BookId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list book chapter: %s", err))
	}

	bookChaptersRespData := make([]api.BookChapter, len(bookChaptersOutput))
	for i, b := range bookChaptersOutput {
		bookChaptersRespData[i] = api.BookChapter{
			BookChapterId: openapi_types.UUID(b.BookChapterID),
			Index:         b.Index,
			Title:         b.Title,
		}
	}

	return ctx.JSON(http.StatusCreated, api.ListBookChapterResponse{
		BookChapters: ptr(bookChaptersRespData),
	})
}

func (p public) StoreBookChapterTranslation(ctx echo.Context) error {
	var input api.StoreBookChapterTranslationRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	translatorID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	err = p.bookChapterTranslationService.StoreBookChapterTranslation(service.StoreBookChapterTranslationInput{
		BookChapterID: domainmodel.BookChapterID(input.BookChapterId),
		TranslatorID:  translatorID,
		Text:          input.Text,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store book chapter translation: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.SuccessResponse{
		Message: ptr("Book chapter translation stored successfully"),
	})
}

func (p public) GetBookChapterTranslation(ctx echo.Context) error {
	var input api.GetBookChapterTranslationRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	bookChapterTranslation, err := p.bookChapterTranslationQueryService.GetByBookChapterIDAndTranslatorID(
		domainmodel.BookChapterID(input.BookChapterId),
		domainmodel.UserID(input.TranslatorId),
	)
	if errors.Is(err, domainmodel.ErrBookChapterTranslationNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Book chapter translation not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get book chapter: %s", err))
	}

	return ctx.JSON(http.StatusCreated, api.GetBookChapterTranslationResponse{
		Text: bookChapterTranslation.Text,
	})
}

func (p public) ListTranslatorsByBookChapterId(ctx echo.Context) error {
	var input api.ListTranslatorsByBookChapterIdRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	translatorsID, err := p.bookChapterTranslationQueryService.ListTranslatorsByBookChapterId(
		domainmodel.BookChapterID(input.BookChapterId),
	)
	if errors.Is(err, domainmodel.ErrBookChapterTranslationNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Book chapter translation not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list translators book chapter: %s", err))
	}

	translatorsRespID := make([]openapi_types.UUID, len(translatorsID))
	for i, t := range translatorsID {
		translatorsRespID[i] = openapi_types.UUID(t)
	}

	return ctx.JSON(http.StatusCreated, api.ListTranslatorsByBookChapterIdResponse{
		TranslatorsId: translatorsRespID,
	})
}

func (p public) StoreReadingSession(ctx echo.Context) error {
	var input api.StoreReadingSessionRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	userID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	err = p.readingSession.StoreReadingSession(service.StoreReadingSessionInput{
		BookID:        domainmodel.BookID(input.BookId),
		BookChapterID: domainmodel.BookChapterID(input.BookChapterId),
		UserID:        userID,
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

	translatorID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	err = p.verifyBookRequestService.CreateVerifyBookRequest(service.CreateVerifyBookRequestInput{
		TranslatorID: translatorID,
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

func ptr[T any](s T) *T {
	return &s
}

func extractUserIDFromContext(ctx echo.Context) (domainmodel.UserID, error) {
	tokenString := ctx.Request().Header.Get("Authorization")
	if tokenString == "" {
		return domainmodel.UserID{}, echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return domainmodel.UserID{}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	var userID uuid.UUID
	err = userID.Parse(claims.UserID)
	return domainmodel.UserID(userID), err
}
