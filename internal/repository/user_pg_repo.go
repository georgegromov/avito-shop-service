package repository

import (
	"avito-shop-service/internal/domain/user"
	"avito-shop-service/internal/domain/wallet"
	"avito-shop-service/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type userRepo struct {
	log *slog.Logger
	db  *sqlx.DB
}

func NewUserRepo(log *slog.Logger, db *sqlx.DB) user.UserRepo {
	return &userRepo{log: log, db: db}
}

// Create implements user.UserRepo.
func (r *userRepo) Create(ctx context.Context, u *user.User) (user.UserID, error) {
	const op = "repository.userRepo.Create"

	log := r.log.With(slog.String("op", op))

	log.Info("attempt to create a new user with the username...", slog.String("username", u.Username))

	var userID user.UserID

	err := utils.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		userQuery := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`
		if err := tx.QueryRowContext(ctx, userQuery, u.Username, u.PasswordHash).Scan(&userID); err != nil {
			log.Error("an error occurred while creating new user", slog.String("username", u.Username), slog.Any("err", err))
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				return user.ErrUserAlreadyExists
			}
			return err
		}

		walletQuery := `INSERT INTO wallets (user_id, balance) VALUES ($1, $2)`
		if _, err := tx.ExecContext(ctx, walletQuery, userID, wallet.InitWalletBalance); err != nil {
			log.Error("an error occurred while creating wallet for new user", slog.String("username", u.Username), slog.Any("err", err))
			return err
		}

		log.Info("new user and wallet created successfully", slog.String("username", u.Username), slog.Any("user_id", userID))
		return nil
	})

	if err != nil {
		return user.UserID(uuid.Nil), err
	}

	return userID, err
}

// GetByUsername implements user.UserRepo.
func (r *userRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	const op = "repository.userRepo.GetByUsername"

	log := r.log.With(slog.String("op", op))

	log.Info("getting user by username", slog.String("username", username))

	q := `SELECT * FROM users WHERE username = $1`
	u := &user.User{}

	if err := r.db.GetContext(ctx, u, q, username); err != nil {
		log.Error("an error occurred while getting user by username", slog.String("username", u.Username), slog.Any("err", err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}
