package user

import (
	"context"

	"github.com/gin-gonic/gin"
)

type UserRepo interface {
	Create(ctx context.Context, u *User) (UserID, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	// GetByID(ctx context.Context, userID UserID) (*User, error)
}

type UserService interface {
	SignUp(ctx context.Context, username, password string) (*User, error)
	SignIn(ctx context.Context, username, password string) (*User, error)
}

type UserHandler interface {
	SignUpRoute(*gin.Context)
	SignInRoute(*gin.Context)
}
