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

func New(username, passwordHash string) (*User, error) {
    if username == "" {
        return nil, ErrInvalidUsername
    }
    if passwordHash == "" {
        return nil, ErrInvalidPassword
    }

    return &User{
        Username:  username,
        PasswordHash:  passwordHash,
        CreatedAt: time.Now(),
    }, nil
}