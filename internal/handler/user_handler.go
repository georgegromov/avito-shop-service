package handler

import (
	"avito-shop-service/internal/domain/user"
	infrahttp "avito-shop-service/internal/infra/http"
	"avito-shop-service/internal/infra/security"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	log         *slog.Logger
	validator   infrahttp.Validator
	userService user.UserService
	jwtManager  security.JwtManager
}

func NewUserHandler(
	log *slog.Logger,
	validator infrahttp.Validator,
	userService user.UserService,
	jwtManager security.JwtManager,
) user.UserHandler {
	return &userHandler{
		log:         log,
		validator:   validator,
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// SignUpRoute implements user.UserHandler.
func (h *userHandler) SignUpRoute(gc *gin.Context) {
	const op = "handler.userHandler.SignUpRoute"

	ctx := gc.Request.Context()

	log := h.log.With(slog.String("op", op))

	var reqBody user.UserSignUpRequestDTO
	if err := gc.ShouldBindJSON(&reqBody); err != nil {
		log.Error("an error occurred while binding request body", slog.Any("err", err))
		infrahttp.SendErrorResponse(
			gc,
			http.StatusBadRequest,
			infrahttp.ErrBadRequest.Error(),
			gin.H{"error": infrahttp.ErrBadRequest.Error()},
		)
		return
	}

	if err := h.validator.Validate(ctx, reqBody); err != nil {
		log.Error("an error occurred while validating request body", slog.Any("err", err))
		infrahttp.SendErrorResponse(
			gc,
			http.StatusBadRequest,
			infrahttp.ErrBadRequest.Error(),
			gin.H{"error": infrahttp.ErrBadRequest.Error()},
		)
		return
	}

	usr, err := h.userService.SignUp(ctx, reqBody.Username, reqBody.Password)
	if err != nil {
		log.Error("an error occurred while signing up user", slog.Any("err", err))
		if err == user.ErrUserAlreadyExists {
			infrahttp.SendErrorResponse(
				gc,
				http.StatusConflict,
				user.ErrUserAlreadyExists.Error(),
				gin.H{"error": user.ErrUserAlreadyExists.Error()},
			)
			return
		}
		infrahttp.SendErrorResponse(
			gc,
			http.StatusInternalServerError,
			infrahttp.ErrInternalServer.Error(),
			gin.H{"error": infrahttp.ErrInternalServer.Error()},
		)
		return
	}

	accessToken, err := h.jwtManager.Generate(usr.ID)
	if err != nil {
		log.Error(security.ErrGenerateToken.Error(), slog.Any("err", err))
		infrahttp.SendErrorResponse(
			gc,
			http.StatusInternalServerError,
			infrahttp.ErrInternalServer.Error(),
			gin.H{"error": infrahttp.ErrInternalServer.Error()},
		)
		return
	}

	respBody := &user.UserSignUpResponseDTO{
		User:        *usr,
		AccessToken: accessToken,
	}

	http.SetCookie(gc.Writer, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		Domain:   "",
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 дней
		MaxAge:   3600 * 24 * 7,                      // В секундах (7 дней)
		HttpOnly: true,                               // Защита от доступа из JS
		Secure:   true,                               // Только HTTPS
		SameSite: http.SameSiteNoneMode,              // Запрещает кросс-доменные запросы с куки
	})

	infrahttp.SendSuccessResponse(
		gc,
		http.StatusCreated,
		"user successfully signed up",
		respBody,
	)
}

// SignInRoute implements user.UserHandler.
func (h *userHandler) SignInRoute(gc *gin.Context) {
	const op = "handler.userHandler.SignInRoute"

	ctx := gc.Request.Context()

	log := h.log.With(slog.String("op", op))

	var reqBody user.UserSignInRequestDTO
	if err := gc.ShouldBindJSON(&reqBody); err != nil {
		log.Error("failed to bind request body", slog.Any("err", err))
		infrahttp.SendErrorResponse(
			gc,
			http.StatusBadRequest,
			infrahttp.ErrBadRequest.Error(),
			gin.H{"error": infrahttp.ErrBadRequest.Error()},
		)
		return
	}

	if err := h.validator.Validate(ctx, reqBody); err != nil {
		log.Error("validation failed", slog.Any("err", err))
		infrahttp.SendErrorResponse(
			gc,
			http.StatusBadRequest,
			infrahttp.ErrBadRequest.Error(),
			gin.H{"error": infrahttp.ErrBadRequest.Error()},
		)
		return
	}

	usr, err := h.userService.SignIn(ctx, reqBody.Username, reqBody.Password)
	if err != nil {
		log.Error("sign in failed", slog.Any("err", err))
		status := http.StatusInternalServerError
		errMsg := infrahttp.ErrInternalServer.Error()

		if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrInvalidPassword) {
			status = http.StatusUnauthorized
			errMsg = user.ErrInvalidCredentials.Error()
		}

		infrahttp.SendErrorResponse(
			gc,
			status,
			errMsg,
			gin.H{"error": errMsg},
		)
		return
	}

	accessToken, err := h.jwtManager.Generate(usr.ID)
	if err != nil {
		log.Error(security.ErrGenerateToken.Error(), slog.Any("err", err))
		infrahttp.SendErrorResponse(
			gc,
			http.StatusInternalServerError,
			infrahttp.ErrInternalServer.Error(),
			gin.H{"error": infrahttp.ErrInternalServer.Error()},
		)
		return
	}

	http.SetCookie(gc.Writer, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		Domain:   "",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		MaxAge:   3600 * 24 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	respBody := &user.UserSignInResponseDTO{
		User:        *usr,
		AccessToken: accessToken,
	}

	infrahttp.SendSuccessResponse(
		gc,
		http.StatusOK,
		"user successfully signed in",
		respBody,
	)
}
