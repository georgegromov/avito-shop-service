package repository

import (
	"avito-shop-service/internal/domain/user"
	"avito-shop-service/internal/domain/wallet"
	"context"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type walletRepo struct {
	log *slog.Logger
	db  *sqlx.DB
}

func NewWalletRepo(log *slog.Logger, db *sqlx.DB) wallet.WalletRepo {
	return &walletRepo{log: log, db: db}
}

// GetByUserID implements wallet.WalletRepo.
func (r *walletRepo) GetByUserID(ctx context.Context, userID user.UserID) (*wallet.Wallet, error) {
	const op = "repository.walletRepo.GetByUserID"

	log := r.log.With(slog.String("op", op))

	query := `SELECT * FROM wallets WHERE user_id = $1`

	var w wallet.Wallet
	if err := r.db.GetContext(ctx, &w, query, userID); err != nil {
		log.Error("an error occurred while fetching wallet by user ID", slog.Any("user_id", userID), slog.Any("err", err))
		return nil, err
	}

	return &w, nil
}

// UpdateBalance implements wallet.WalletRepo.
func (r *walletRepo) UpdateBalanceTx(ctx context.Context, tx sqlx.Tx, w *wallet.Wallet) error {
	const op = "repository.walletRepo.UpdateBalance"

	log := r.log.With(slog.String("op", op))

	log.Info("updating wallet balance...", slog.Any("user_id", w.UserID))

	w.UpdatedAt = time.Now().UTC()

	query := `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE user_id = $2`

	res, err := tx.ExecContext(ctx, query, w.Balance, w.UserID)
	if err != nil {
		log.Error("an error occurred while updating wallet balance", slog.Any("user_id", w.UserID), slog.Any("err", err))
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return wallet.ErrWalletNotFound
	}

	log.Info("wallet balance updated successfully", slog.Any("user_id", w.UserID))
	return nil
}
