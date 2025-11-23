package security

import (
	"avito-shop-service/internal/domain/user"
	"errors"
)

type JwtManager interface {
	Generate(userID user.UserID) (string, error)
	Validate(token string) (string, error)
}

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)
