package service

import (
	"avito-shop-service/internal/domain/user"
	"avito-shop-service/internal/infra/security"
	"context"
	"errors"
	"log/slog"
)

type userService struct {
	log         *slog.Logger
	repo        user.UserRepo
	hashManager security.HashManager
	jwtManager  security.JwtManager
}

func NewUserService(
	log *slog.Logger,
	repo user.UserRepo,
	hashManager security.HashManager,
	jwtManager security.JwtManager,
) user.UserService {
	return &userService{
		log:         log,
		repo:        repo,
		hashManager: hashManager,
		jwtManager:  jwtManager,
	}
}

// SignUp implements user.UserService.
func (s *userService) SignUp(ctx context.Context, username string, password string) (*user.User, error) {
	const op = "service.userService.SignUp"

	log := s.log.With(slog.String("op", op))

	log.Info("attempt to sign up a user with the username", slog.String("username", username))

	existingUser, err := s.repo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		log.Warn(user.ErrUserAlreadyExists.Error(), slog.String("username", username))
		return nil, user.ErrUserAlreadyExists
	}

	passwordHash, err := s.hashManager.Hash(password)
	if err != nil {
		log.Error("failed to hash password", slog.Any("err", err))
		return nil, err
	}

	usr := &user.User{
		Username:     username,
		PasswordHash: passwordHash,
	}

	userID, err := s.repo.Create(ctx, usr)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			log.Warn(user.ErrUserAlreadyExists.Error(), slog.String("username", username))
			return nil, user.ErrUserAlreadyExists
		}
		return nil, err
	}

	usr.ID = userID

	log.Info("the user has been successfully signed up", slog.Any("user_id", usr.ID))
	return usr, nil
}

// SignIn implements user.UserService.
func (s *userService) SignIn(ctx context.Context, username string, password string) (*user.User, error) {
	const op = "service.userService.SignIn"

	log := s.log.With(slog.String("op", op))

	usr, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		log.Error("an error occurred while getting user by username", slog.Any("err", err))
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	err = s.hashManager.CompareHash(password, usr.PasswordHash)
	if err != nil {
		log.Error("an error occurred while comparing password hash", slog.Any("err", err))
		return nil, user.ErrInvalidPassword
	}

	return usr, nil
}
