package user

import (
	"time"

	"github.com/google/uuid"
)

type UserID uuid.UUID

type User struct {
	ID           UserID
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type UserSignUpRequestDTO struct {
	Username string `json:"username" binding:"required,min=4" validate:"required,min=4"`
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`
}
type UserSignInRequestDTO struct {
	Username string `json:"username" binding:"required,min=4" validate:"required,min=4"`
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`
}

type UserSignUpResponseDTO struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

type UserSignInResponseDTO struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

func New(username, passwordHash string) (*User, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}
	if passwordHash == "" {
		return nil, ErrInvalidPassword
	}

	return &User{
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}, nil
}
