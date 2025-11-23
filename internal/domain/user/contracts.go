package user

import "context"

type UserRepo interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, userID UserID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type UserService interface {
	SignUp(ctx context.Context, username, password string) (*User, error)
	SignIn(ctx context.Context, username, password string) (*User, error)
}
