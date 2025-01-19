package transport

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
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
	imageService service.ImageService,
	bookRatingService service.BookRatingService,
	userBookFavouritesService service.UserBookFavouritesService,
	authorService service.AuthorService,
	bookAuthorService service.BookAuthorService,
	genreService service.GenreService,
	bookGenreService service.BookGenreService,

	userQueryService query.UserQueryService,
	bookQueryService query.BookQueryService,
	bookChapterQueryService query.BookChapterQueryService,
	bookChapterTranslationQueryService query.BookChapterTranslationQueryService,
	verifyBookRequestQueryService query.VerifyBookRequestQueryService,
	readingSessionQueryService query.ReadingSessionQueryService,
	imageQueryService query.ImageQueryService,
	userBookFavouritesQueryService query.UserBookFavouritesQueryService,
	authorQueryService query.AuthorQueryService,
	genreQueryService query.GenreQueryService,

	verifyBookRequestProvider provider.VerifyBookRequestProvider,
) api.ServerInterface {
	return &public{
		userService:                   userService,
		bookService:                   bookService,
		bookChapterService:            bookChapterService,
		bookChapterTranslationService: bookChapterTranslationService,
		readingSession:                readingSession,
		verifyBookRequestService:      verifyBookRequestService,
		bookRatingService:             bookRatingService,
		imageService:                  imageService,
		userBookFavouritesService:     userBookFavouritesService,
		authorService:                 authorService,
		bookAuthorService:             bookAuthorService,
		genreService:                  genreService,
		bookGenreService:              bookGenreService,

		userQueryService:                   userQueryService,
		bookQueryService:                   bookQueryService,
		bookChapterQueryService:            bookChapterQueryService,
		bookChapterTranslationQueryService: bookChapterTranslationQueryService,
		verifyBookRequestQueryService:      verifyBookRequestQueryService,
		readingSessionQueryService:         readingSessionQueryService,
		imageQueryService:                  imageQueryService,
		userBookFavouritesQueryService:     userBookFavouritesQueryService,
		authorQueryService:                 authorQueryService,
		genreQueryService:                  genreQueryService,

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
	imageService                  service.ImageService
	userBookFavouritesService     service.UserBookFavouritesService
	authorService                 service.AuthorService
	bookAuthorService             service.BookAuthorService
	genreService                  service.GenreService
	bookGenreService              service.BookGenreService

	userQueryService                   query.UserQueryService
	bookQueryService                   query.BookQueryService
	bookChapterQueryService            query.BookChapterQueryService
	bookChapterTranslationQueryService query.BookChapterTranslationQueryService
	verifyBookRequestQueryService      query.VerifyBookRequestQueryService
	readingSessionQueryService         query.ReadingSessionQueryService
	imageQueryService                  query.ImageQueryService
	userBookFavouritesQueryService     query.UserBookFavouritesQueryService
	authorQueryService                 query.AuthorQueryService
	genreQueryService                  query.GenreQueryService

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

	return ctx.JSON(http.StatusOK, api.GetBookRatingResponse{
		Average: ptr(float32(stat.Average)),
		Count:   ptr(stat.Count),
	})
}

func (p public) LoginUser(ctx echo.Context) error {
	var userReq api.LoginUserRequest
	if err := ctx.Bind(&userReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	user, err := p.userQueryService.FindByLogin(userReq.Login)
	if errors.Is(err, domainmodel.ErrUserNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	if user == (model.User{}) || user.Password != userReq.Password {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	if err != nil {
		return err
	}

	accessToken, accessExpirationTime, err := createToken(user, 5*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create access token")
	}

	refreshToken, refreshExpirationTime, err := createToken(user, 30*24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create refresh token")
	}

	setCookie(ctx, "access_token", accessToken, accessExpirationTime)

	setCookie(ctx, "refresh_token", refreshToken, refreshExpirationTime)

	return ctx.NoContent(http.StatusOK)
}

func (p public) RefreshToken(ctx echo.Context) error {
	refreshTokenCookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token not found")
	}

	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(refreshTokenCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil || !token.Valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token expired")
	}

	user, err := p.userQueryService.FindByLogin(claims.Login)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	accessToken, accessExpirationTime, err := createToken(user, 5*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create access token")
	}

	setCookie(ctx, "access_token", accessToken, accessExpirationTime)

	msg := "Token refreshed successfully"
	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr(msg),
	})
}

func (p public) LogoutUser(ctx echo.Context) error {
	deleteCookie(ctx, "access_token")
	deleteCookie(ctx, "refresh_token")

	return ctx.NoContent(http.StatusOK)
}

func (p public) GetLoginUser(ctx echo.Context) error {
	userID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	user, err := p.userQueryService.FindByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get login user: %s", err))
	}

	return ctx.JSON(http.StatusOK, convertUserModelToAPI(user))
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
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

	user, err := p.userQueryService.FindByID(domainmodel.UserID(input.Id))
	if errors.Is(err, domainmodel.ErrUserNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to edit user: %s", err))
	}

	isAdmin, err := checkIsAdmin(ctx)
	if err != nil {
		return err
	}

	if user.Password != input.ConfirmPassword && !isAdmin {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Failed to edit user password: %s", err))
	}

	if user.Login != input.Login && !isAdmin {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Failed to edit user login: %s", err))
	}

	err = p.userService.EditUser(service.EditUserInput{
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("User deleted successfully"),
	})
}

func (p public) GetUser(ctx echo.Context, id string) error {
	user, err := p.userQueryService.FindByID(domainmodel.UserID(uuid.FromStringOrNil(id)))
	if errors.Is(err, domainmodel.ErrUserNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to edit user: %s", err))
	}

	return ctx.JSON(http.StatusOK, convertUserModelToAPI(user))
}

func (p public) ListUser(ctx echo.Context) error {
	usersOutput, err := p.userQueryService.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list user: %s", err))
	}

	users := make([]api.User, len(usersOutput))
	for i, user := range usersOutput {
		users[i] = convertUserModelToAPI(user)
	}

	return ctx.JSON(http.StatusOK, users)
}

func (p public) CreateBook(ctx echo.Context) error {
	var input api.CreateBookRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	bookID, err := p.bookService.CreateBook(service.CreateBookInput{
		Title:       input.Title,
		Description: input.Description,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create book: %s", err))
	}

	userID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	err = p.verifyBookRequestService.CreateVerifyBookRequest(service.CreateVerifyBookRequestInput{
		TranslatorID: userID,
		BookID:       bookID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create verify book request: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Book deleted successfully"),
	})
}

func (p public) SearchBook(ctx echo.Context, queryParams api.SearchBookParams) error {
	spec := convertListBookParamsToListSpec(queryParams)

	if spec.Page < 0 || spec.Size < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Page and size needed positive and not zero")
	}

	bookOutputs, err := p.bookQueryService.List(spec)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list book: %s", err))
	}

	booksRespData := make([]api.Book, len(bookOutputs))
	for i, b := range bookOutputs {
		authors, err2 := p.authorQueryService.ListByBookID(b.BookID)
		if err2 != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list author: %s", err2))
		}
		genres, err2 := p.genreQueryService.ListByBookID(b.BookID)
		if err2 != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list genre: %s", err2))
		}

		booksRespData[i] = convertBookOutputModelToAPI(b, authors, genres)
	}

	countBook, err := p.bookQueryService.CountBook(true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list book: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.ListBookResponse{
		Books:      booksRespData,
		CountPages: ptr(int(math.Ceil(float64(countBook) / float64(spec.Size)))),
	})
}

func (p public) GetBook(ctx echo.Context, id string) error {
	var bookID uuid.UUID
	err := bookID.Parse(id)
	if err != nil {
		return err
	}

	book, err := p.bookQueryService.FindByID(domainmodel.BookID(bookID))
	if errors.Is(err, domainmodel.ErrBookNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Book not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get book: %s", err))
	}

	authors, err2 := p.authorQueryService.ListByBookID(book.BookID)
	if err2 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list author: %s", err2))
	}

	genres, err2 := p.genreQueryService.ListByBookID(book.BookID)
	if err2 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list genre: %s", err2))
	}

	bookRespData := convertBookOutputModelToAPI(book, authors, genres)

	return ctx.JSON(http.StatusOK, api.GetBookResponse{
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
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

	return ctx.JSON(http.StatusOK, api.ListBookChapterResponse{
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

	if input.TranslaterId != nil {
		translatorID = domainmodel.UserID(*input.TranslaterId)
	}

	err = p.bookChapterTranslationService.StoreBookChapterTranslation(service.StoreBookChapterTranslationInput{
		BookChapterID: domainmodel.BookChapterID(input.BookChapterId),
		TranslatorID:  translatorID,
		Text:          input.Text,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store book chapter translation: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
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

	return ctx.JSON(http.StatusOK, api.GetBookChapterTranslationResponse{
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

	return ctx.JSON(http.StatusOK, api.ListTranslatorsByBookChapterIdResponse{
		TranslatorsId: translatorsRespID,
	})
}

func (p public) GetLastReadingSession(ctx echo.Context) error {
	var input api.GetLastReadingSessionRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	userID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	lastReadingSession, err := p.readingSessionQueryService.GetLastReadingSession(domainmodel.BookID(input.BookId), userID)
	if errors.Is(err, domainmodel.ErrBookChapterTranslationNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Reading session not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get last reading session: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.GetLastReadingSessionResponse{
		BookChapterId: openapi_types.UUID(lastReadingSession.BookChapterId),
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Reading session stored successfully"),
	})
}

func (p public) ListVerifyBookRequest(ctx echo.Context) error {
	verifyBooksRequests, err := p.verifyBookRequestQueryService.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list verify book request session: %s", err))
	}

	bookIDs := make([]domainmodel.BookID, 0)
	for _, verifyBooksRequest := range verifyBooksRequests {
		bookIDs = append(bookIDs, verifyBooksRequest.BookID)
	}

	books, err := p.bookQueryService.ListByIDs(bookIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list verify book request session: %s", err))
	}
	bookIDToBookMap := make(map[domainmodel.BookID]query.BookOutput)
	for _, book := range books {
		bookIDToBookMap[book.BookID] = book
	}

	verifyBookRespRequests := make([]api.VerifyBookRequest, len(verifyBooksRequests))
	for i, v := range verifyBooksRequests {
		isVerified, ok := v.IsVerified.Get()

		sendDateMilli := int(v.SendDate.UnixNano() / int64(time.Millisecond))

		authors, err2 := p.authorQueryService.ListByBookID(v.BookID)
		if err2 != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list author: %s", err2))
		}

		genres, err2 := p.genreQueryService.ListByBookID(v.BookID)
		if err2 != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list genre: %s", err2))
		}

		verifyBookRespRequests[i] = api.VerifyBookRequest{
			VerifyBookRequestId: openapi_types.UUID(v.VerifyBookRequestID),
			TranslatorId:        openapi_types.UUID(v.TranslatorID),
			IsVerified:          ptr(isVerified),
			SendDateMilli:       sendDateMilli,
			Book:                convertBookOutputModelToAPI(bookIDToBookMap[v.BookID], authors, genres),
		}

		if !ok {
			verifyBookRespRequests[i].IsVerified = nil
		}
	}

	return ctx.JSON(http.StatusOK, api.ListVerifyBookRequestResponse{
		VerifyBookRequests: verifyBookRespRequests,
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
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

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Verify book request verified successfully"),
	})
}

func (p public) GetImage(ctx echo.Context) error {
	var input api.GetImageRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	imageData, err := p.imageQueryService.FindByID(domainmodel.ImageID(input.ImageId))
	if errors.Is(err, domainmodel.ErrImageNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Image not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get image: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.GetImageResponse{
		ImageData: imageData,
	})
}

func (p public) DeleteImage(ctx echo.Context) error {
	var input api.DeleteImageRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.imageService.DeleteImage(domainmodel.ImageID(input.ImageId))
	if errors.Is(err, domainmodel.ErrImageNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Reading image not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get last reading session: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Image deleted successfully"),
	})
}

func (p public) StoreImageBook(ctx echo.Context) error {
	var input api.StoreBookImageRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	imageID, err := p.imageService.StoreImage(input.ImageData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store image book: %s", err))
	}

	err = p.bookService.EditBookImage(service.EditBookImageInput{
		ID:      domainmodel.BookID(input.BookId),
		ImageID: imageID,
	})
	if errors.Is(err, domainmodel.ErrBookNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "Book not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store image book: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Store image book successfully"),
	})
}

func (p public) StoreImageUser(ctx echo.Context) error {
	var input api.StoreUserImageRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	userID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	imageID, err := p.imageService.StoreImage(input.ImageData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store image user: %s", err))
	}

	err = p.userService.EditImageUser(service.EditUserImageInput{
		ID:      userID,
		ImageID: imageID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store image user: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Store image user successfully"),
	})
}

func (p public) StoreImageAuthor(ctx echo.Context) error {
	var input api.StoreAuthorImageRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	imageID, err := p.imageService.StoreImage(input.ImageData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store image author: %s", err))
	}

	err = p.authorService.EditAuthorAvatar(service.EditAuthorAvatarInput{
		ID:      domainmodel.AuthorID(input.AuthorId),
		ImageID: imageID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store image author: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Store author user successfully"),
	})
}

func (p public) ListUserBookFavouritesByBook(ctx echo.Context) error {
	var input api.ListUserBookFavouritesByBookRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	userID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	modelType, err := p.userBookFavouritesQueryService.ListUserBookFavouritesByBook(userID, domainmodel.BookID(input.BookId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list user book favourites: %s", err))
	}

	apiType, err := convertUserBookFavouritesTypeModelToAPI(modelType)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, api.ListUserBookFavouritesByBookResponse{
		Type: apiType,
	})
}

func (p public) StoreUserBookFavourites(ctx echo.Context) error {
	var input api.StoreUserBookFavouritesRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	userID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	modelType, err := convertUserBookFavouritesTypeAPIToModel(input.Type)
	if err != nil {
		return err
	}

	err = p.userBookFavouritesService.StoreUserBookFavourites(service.StoreUserBookFavouritesInput{
		UserID: userID,
		BookID: domainmodel.BookID(input.BookId),
		Type:   modelType,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store user book favourites: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Store user book favourites successfully"),
	})
}

func (p public) DeleteUserBookFavourites(ctx echo.Context) error {
	var input api.DeleteUserBookFavouritesRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	userID, err := extractUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	err = p.userBookFavouritesService.DeleteUserBookFavourites(service.DeleteUserBookFavouritesInput{
		UserID: userID,
		BookID: domainmodel.BookID(input.BookId),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete user book favourites: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Delete user book favourites successfully"),
	})
}

func (p public) ListBookByUserBookFavourites(ctx echo.Context) error {
	var input api.ListBookByUserBookFavouritesRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	modelTypes := make([]domainmodel.UserBookFavouritesType, len(input.Types))
	for i, t := range input.Types {
		modelType, err2 := convertUserBookFavouritesTypeAPIToModel(t)
		if err2 != nil {
			return err2
		}

		modelTypes[i] = modelType
	}

	outputs, err := p.userBookFavouritesQueryService.ListBookByUserBookFavourites(domainmodel.UserID(input.UserId), modelTypes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list user book favourites: %s", err))
	}

	userBookFavouritesBooks := make([]api.UserBookFavouritesBooks, len(outputs))
	for i, output := range outputs {
		apiType, err2 := convertUserBookFavouritesTypeModelToAPI(output.Type)
		if err2 != nil {
			return err2
		}

		books := make([]api.Book, len(output.Books))
		for j, book := range output.Books {
			authors, err3 := p.authorQueryService.ListByBookID(book.BookID)
			if err3 != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list author: %s", err3))
			}

			genres, err3 := p.genreQueryService.ListByBookID(book.BookID)
			if err3 != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list genre: %s", err3))
			}

			books[j] = convertBookOutputModelToAPI(book, authors, genres)
		}

		userBookFavouritesBooks[i] = api.UserBookFavouritesBooks{
			Type:  apiType,
			Books: books,
		}
	}

	return ctx.JSON(http.StatusOK, api.ListBookByUserBookFavouritesResponse{
		UserBookFavouritesBooks: userBookFavouritesBooks,
	})
}

func (p public) ListAuthors(ctx echo.Context) error {
	outputs, err := p.authorQueryService.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list authors: %s", err))
	}

	authorsRespData := make([]api.Author, len(outputs))
	for i, a := range outputs {
		authorsRespData[i] = convertAuthorOutputModelToAPI(a)
	}

	return ctx.JSON(http.StatusOK, api.ListAuthorResponse{
		Authors: authorsRespData,
	})
}

func (p public) CreateAuthor(ctx echo.Context) error {
	var input api.CreateAuthorRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	middleName := maybe.Nothing[string]()
	if input.MiddleName != nil {
		middleName = maybe.Just(*input.MiddleName)
	}

	nickName := maybe.Nothing[string]()
	if input.MiddleName != nil {
		nickName = maybe.Just(*input.NickName)
	}

	err := p.authorService.CreateAuthor(service.CreateAuthorInput{
		FirstName:  input.FirstName,
		SecondName: input.SecondName,
		MiddleName: middleName,
		Nickname:   nickName,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create author: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Author created successfully"),
	})
}

func (p public) EditAuthor(ctx echo.Context) error {
	var input api.EditAuthorRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	middleName := maybe.Nothing[string]()
	if input.MiddleName != nil {
		middleName = maybe.Just(*input.MiddleName)
	}

	nickName := maybe.Nothing[string]()
	if input.MiddleName != nil {
		nickName = maybe.Just(*input.NickName)
	}

	err := p.authorService.EditAuthor(service.EditAuthorInput{
		ID:         domainmodel.AuthorID(input.Id),
		FirstName:  input.FirstName,
		SecondName: input.SecondName,
		MiddleName: middleName,
		Nickname:   nickName,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to edit author: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Author edited successfully"),
	})
}

func (p public) DeleteAuthor(ctx echo.Context) error {
	var input api.DeleteAuthorRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.authorService.DeleteAuthor(domainmodel.AuthorID(input.Id))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete author: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Author deleted successfully"),
	})
}

func (p public) GetAuthor(ctx echo.Context, id string) error {
	uid, err := uuid.FromString(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get author: %s", err))
	}

	output, err := p.authorQueryService.FindByID(domainmodel.AuthorID(uid))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get author: %s", err))
	}

	return ctx.JSON(http.StatusOK, convertAuthorOutputModelToAPI(output))
}

func (p public) StoreBookAuthor(ctx echo.Context) error {
	var input api.StoreBookAuthorRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.bookAuthorService.StoreBookAuthor(domainmodel.BookID(input.BookId), domainmodel.GenreID(input.AuthorId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store book author: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Book author stored successfully"),
	})
}

func (p public) DeleteBookAuthor(ctx echo.Context) error {
	var input api.DeleteBookAuthorRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.bookAuthorService.DeleteBookAuthor(domainmodel.BookID(input.BookId), domainmodel.GenreID(input.AuthorId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete book author: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Book author deleted successfully"),
	})
}

func (p public) ListGenres(ctx echo.Context) error {
	outputs, err := p.genreQueryService.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to list genres: %s", err))
	}

	genresRespData := make([]api.Genre, len(outputs))
	for i, g := range outputs {
		genresRespData[i] = convertGenreOutputModelToAPI(g)
	}

	return ctx.JSON(http.StatusOK, api.ListGenreResponse{
		Genres: genresRespData,
	})
}

func (p public) CreateGenre(ctx echo.Context) error {
	var input api.CreateGenreRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.genreService.CreateGenre(service.CreateGenreInput{
		Name: input.Name,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create genre: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Genre created successfully"),
	})
}

func (p public) EditGenre(ctx echo.Context) error {
	var input api.EditGenreRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.genreService.EditGenre(service.EditGenreInput{
		ID:   domainmodel.GenreID(input.Id),
		Name: input.Name,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to edit genre: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Genre edited successfully"),
	})
}

func (p public) DeleteGenre(ctx echo.Context) error {
	var input api.DeleteGenreRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.genreService.DeleteGenre(domainmodel.GenreID(input.Id))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete genre: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Genre deleted successfully"),
	})
}

func (p public) StoreBookGenre(ctx echo.Context) error {
	var input api.StoreBookGenreRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.bookGenreService.StoreBookGenre(domainmodel.BookID(input.BookId), domainmodel.GenreID(input.GenreId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to store book genre: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Book genre stored successfully"),
	})
}

func (p public) DeleteBookGenre(ctx echo.Context) error {
	var input api.DeleteBookGenreRequest
	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.BadRequestResponse{
			Message: ptr(fmt.Sprintf("Invalid request: %s", err)),
		})
	}

	err := p.bookGenreService.DeleteBookGenre(domainmodel.BookID(input.BookId), domainmodel.GenreID(input.GenreId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to delete book genre: %s", err))
	}

	return ctx.JSON(http.StatusOK, api.SuccessResponse{
		Message: ptr("Book genre deleted successfully"),
	})
}

func ptr[T any](s T) *T {
	return &s
}

func extractUserIDFromContext(ctx echo.Context) (domainmodel.UserID, error) {
	claims, err := parseClaims(ctx)
	if err != nil {
		return domainmodel.UserID{}, err
	}

	var userID uuid.UUID
	err = userID.Parse(claims.UserID)
	return domainmodel.UserID(userID), err
}

func checkIsAdmin(ctx echo.Context) (bool, error) {
	claims, err := parseClaims(ctx)
	if err != nil {
		return false, err
	}

	return domainmodel.UserRole(claims.Role) == domainmodel.Admin, nil
}

func parseClaims(ctx echo.Context) (model.Claims, error) {
	tokenString, err := ctx.Cookie("access_token")
	if err != nil {
		return model.Claims{}, err
	}
	if tokenString.Value == "" {
		return model.Claims{}, echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenString.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return model.Claims{}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	return *claims, nil
}

func convertUserBookFavouritesTypeAPIToModel(apiType api.UserBookFavouritesType) (domainmodel.UserBookFavouritesType, error) {
	switch apiType {
	case api.READING:
		return domainmodel.READING, nil
	case api.PLANNED:
		return domainmodel.PLANNED, nil
	case api.DEFERRED:
		return domainmodel.DEFERRED, nil
	case api.READ:
		return domainmodel.READ, nil
	case api.DROPPED:
		return domainmodel.DROPPED, nil
	case api.FAVORITE:
		return domainmodel.FAVORITE, nil
	default:
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Unknown UserBookFavouritesType "+apiType)
	}
}

func convertUserBookFavouritesTypeModelToAPI(modelType domainmodel.UserBookFavouritesType) (api.UserBookFavouritesType, error) {
	switch modelType {
	case domainmodel.READING:
		return api.READING, nil
	case domainmodel.PLANNED:
		return api.PLANNED, nil
	case domainmodel.DEFERRED:
		return api.DEFERRED, nil
	case domainmodel.READ:
		return api.READ, nil
	case domainmodel.DROPPED:
		return api.DROPPED, nil
	case domainmodel.FAVORITE:
		return api.FAVORITE, nil
	default:
		return "", echo.NewHTTPError(http.StatusInternalServerError, "Unknown UserBookFavouritesType "+fmt.Sprint(modelType))
	}
}

func convertBookOutputModelToAPI(
	bookOutput query.BookOutput,
	authors []query.AuthorOutput,
	genres []query.GenreOutput,
) api.Book {
	authorsAPI := make([]api.Author, len(authors))
	for i, author := range authors {
		authorsAPI[i] = convertAuthorOutputModelToAPI(author)
	}

	genresAPI := make([]api.Genre, len(genres))
	for i, genre := range genres {
		genresAPI[i] = convertGenreOutputModelToAPI(genre)
	}

	cover, ok := bookOutput.Cover.Get()

	bookAPI := api.Book{
		BookId:                 openapi_types.UUID(bookOutput.BookID),
		Cover:                  ptr(cover),
		Title:                  bookOutput.Title,
		Description:            bookOutput.Description,
		Authors:                authorsAPI,
		Genres:                 genresAPI,
		IsLoggedUserTranslator: true,
		Rating:                 float32(bookOutput.AverageRating),
	}

	if !ok {
		bookAPI.Cover = nil
	}

	return bookAPI
}

func convertAuthorOutputModelToAPI(output query.AuthorOutput) api.Author {
	var middleName *string
	if middleNameValue, ok := output.MiddleName.Get(); ok {
		middleName = &middleNameValue
	}

	var nickname *string
	if nicknameValue, ok := output.Nickname.Get(); ok {
		nickname = &nicknameValue
	}

	return api.Author{
		Id:         openapi_types.UUID(output.AuthorID),
		FirstName:  output.FirstName,
		SecondName: output.SecondName,
		MiddleName: middleName,
		Nickname:   nickname,
	}
}

func convertGenreOutputModelToAPI(output query.GenreOutput) api.Genre {
	return api.Genre{
		Id:   openapi_types.UUID(output.GenreID),
		Name: output.Name,
	}
}

func convertListBookParamsToListSpec(params api.SearchBookParams) query.ListSpec {
	spec := query.ListSpec{
		Page: params.Page,
		Size: params.Size,
	}

	if params.BookTitle != nil {
		spec.BookTitle = maybe.Just(*params.BookTitle)
	} else {
		spec.BookTitle = maybe.Nothing[string]()
	}

	if params.AuthorIds != nil {
		authorIDs := make([]domainmodel.AuthorID, len(*params.AuthorIds))
		for i, authorID := range *params.AuthorIds {
			authorIDs[i] = domainmodel.AuthorID(uuid.FromStringOrNil(authorID))
		}

		spec.AuthorIDs = maybe.Just(authorIDs)
	} else {
		spec.AuthorIDs = maybe.Nothing[[]domainmodel.AuthorID]()
	}

	if params.GenreIds != nil {
		genreIDs := make([]domainmodel.GenreID, len(*params.GenreIds))
		for i, genreID := range *params.GenreIds {
			genreIDs[i] = domainmodel.GenreID(uuid.FromStringOrNil(genreID))
		}

		spec.GenreIDs = maybe.Just(genreIDs)
	} else {
		spec.GenreIDs = maybe.Nothing[[]domainmodel.GenreID]()
	}

	if params.MinRating != nil {
		spec.MinRating = maybe.Just(float64(*params.MinRating))
	} else {
		spec.MinRating = maybe.Nothing[float64]()
	}

	if params.MaxRating != nil {
		spec.MaxRating = maybe.Just(float64(*params.MaxRating))
	} else {
		spec.MaxRating = maybe.Nothing[float64]()
	}

	if params.MinChaptersCount != nil {
		spec.MinChaptersCount = maybe.Just(*params.MinChaptersCount)
	} else {
		spec.MinChaptersCount = maybe.Nothing[int]()
	}

	if params.MaxChaptersCount != nil {
		spec.MaxChaptersCount = maybe.Just(*params.MaxChaptersCount)
	} else {
		spec.MaxChaptersCount = maybe.Nothing[int]()
	}

	if params.MinRatingCount != nil {
		spec.MinRatingCount = maybe.Just(*params.MinRatingCount)
	} else {
		spec.MinRatingCount = maybe.Nothing[int]()
	}

	if params.MaxRatingCount != nil {
		spec.MaxRatingCount = maybe.Just(*params.MaxRatingCount)
	} else {
		spec.MaxRatingCount = maybe.Nothing[int]()
	}

	if params.SortBy != nil {
		spec.SortBy = maybe.Just(string(*params.SortBy))
	} else {
		spec.SortBy = maybe.Nothing[string]()
	}

	if params.SortType != nil {
		spec.SortType = maybe.Just(string(*params.SortType))
	} else {
		spec.SortType = maybe.Nothing[string]()
	}

	return spec
}

func convertUserModelToAPI(user model.User) api.User {
	avatar, ok := user.Avatar.Get()

	userAPI := api.User{
		Id:      openapi_types.UUID(user.ID),
		Avatar:  ptr(avatar),
		Login:   user.Login,
		Role:    user.Role,
		AboutMe: user.AboutMe,
	}

	if !ok {
		userAPI.Avatar = nil
	}

	return userAPI
}

func createToken(user model.User, expirationTimeDur time.Duration) (string, time.Time, error) {
	expirationTime := time.Now().Add(expirationTimeDur)
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
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func setCookie(ctx echo.Context, name, value string, expirationTime time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = expirationTime
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode

	ctx.SetCookie(cookie)
}

func deleteCookie(ctx echo.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	ctx.SetCookie(cookie)
}
