package repository

import (
	"avito-shop-service/internal/domain/transfer"
	"avito-shop-service/internal/domain/user"
	"avito-shop-service/internal/domain/wallet"
	"avito-shop-service/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type transferRepo struct {
	log *slog.Logger
	db  *sqlx.DB
}

func NewTransferRepo(log *slog.Logger, db *sqlx.DB) transfer.TransferRepo {
	return &transferRepo{log: log, db: db}
}

// SendCoins implements transfer.TransferRepo.
func (r *transferRepo) SendCoins(ctx context.Context, fromID user.UserID, toID user.UserID, amount wallet.Amount) error {
	const op = "repository.transferRepo.SendCoins"
	log := r.log.With(slog.String("op", op))

	if fromID == toID {
		log.Warn("attempt to send coins to yourself", slog.Any("user_id", fromID))
		return transfer.ErrCannotSendToYourself
	}

	if amount == 0 {
		return transfer.ErrInvalidAmount
	}

	return utils.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var fromBalance wallet.Amount
		queryFrom := `SELECT balance FROM wallets WHERE user_id = $1 FOR UPDATE`
		if err := tx.QueryRowContext(ctx, queryFrom, fromID).Scan(&fromBalance); err != nil {
			log.Error("an error occurred while getting sender wallet", slog.Any("err", err))
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("sender wallet not found: %w", wallet.ErrWalletNotFound)
			}
			return err
		}

		if fromBalance < wallet.Amount(amount) {
			return wallet.ErrInsufficientFunds
		}

		var toBalance wallet.Amount
		queryTo := `SELECT balance FROM wallets WHERE user_id = $1 FOR UPDATE`
		if err := tx.QueryRowContext(ctx, queryTo, toID).Scan(&toBalance); err != nil {
			log.Error("an error occurred while getting receiver wallet", slog.Any("err", err))
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("receiver wallet not found: %w", wallet.ErrWalletNotFound)
			}
			return err
		}

		updateFrom := `UPDATE wallets SET balance = balance - $1, updated_at = now() WHERE user_id = $2`
		if _, err := tx.ExecContext(ctx, updateFrom, amount, fromID); err != nil {
			log.Error("an error occurred while updating sender wallet", slog.Any("err", err))
			return err
		}

		updateTo := `UPDATE wallets SET balance = balance + $1, updated_at = now() WHERE user_id = $2`
		if _, err := tx.ExecContext(ctx, updateTo, amount, toID); err != nil {
			log.Error("an error occurred while updating receiver wallet", slog.Any("err", err))
			return err
		}

		insertTransfer := `INSERT INTO transfers (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)`
		if _, err := tx.ExecContext(ctx, insertTransfer, fromID, toID, amount); err != nil {
			log.Error("an error occurred while inserting transfer record", slog.Any("err", err))
			return err
		}

		log.Info("transfer completed successfully", slog.Any("from_user_id", fromID), slog.Any("to_user_id", toID), slog.Any("amount", amount))
		return nil
	})
}

// GetHistory implements transfer.TransferRepo.
func (r *transferRepo) GetHistory(ctx context.Context, userID user.UserID) ([]*transfer.Transfer, error) {
	const op = "repository.transferRepo.GetHistory"
	log := r.log.With(slog.String("op", op))

	transfers := []*transfer.Transfer{}

	q := `SELECT * FROM transfers WHERE from_user_id = $1 OR to_user_id = $1 ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &transfers, q)
	if err != nil {
		log.Error("an error occurred while getting transfer history", slog.Any("err", err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, transfer.ErrHistoryEmpty
		}
		return nil, err
	}

	return transfers, nil
}
