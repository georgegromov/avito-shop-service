package repository

import (
	"avito-shop-service/internal/domain/item"
	"avito-shop-service/internal/domain/purchase"
	"avito-shop-service/internal/domain/user"
	"avito-shop-service/internal/domain/wallet"
	"avito-shop-service/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type purchaseRepo struct {
	log *slog.Logger
	db  *sqlx.DB
}

func NewPurchaseRepo(log *slog.Logger, db *sqlx.DB) purchase.PurchaseRepo {
	return &purchaseRepo{log: log, db: db}
}

// BuyItem implements purchase.PurchaseRepo.
func (p *purchaseRepo) BuyItem(ctx context.Context, userID user.UserID, itemID item.ItemID, quantity uint) error {
	const op = "repository.purchaseRepo.BuyItem"
	log := p.log.With(slog.String("op", op))

	return utils.WithTx(ctx, p.db, func(tx *sqlx.Tx) error {
		var price item.Price
		itemQuery := `SELECT price FROM items WHERE id = $1`
		if err := tx.QueryRowContext(ctx, itemQuery, itemID).Scan(&price); err != nil {
			log.Error("an error occurred while getting item price", slog.Any("err", err))
			if errors.Is(err, sql.ErrNoRows) {
				return item.ErrItemNotFound
			}
			return err
		}

		totalPrice := price * item.Price(quantity)

		var balance uint64
		walletQuery := `SELECT balance FROM wallets WHERE user_id = $1 FOR UPDATE`
		if err := tx.QueryRowContext(ctx, walletQuery, userID).Scan(&balance); err != nil {
			log.Error("an error occurred while getting wallet balance", slog.Any("err", err))
			if errors.Is(err, sql.ErrNoRows) {
				return wallet.ErrWalletNotFound
			}
			return err
		}

		if item.Price(balance) < totalPrice {
			log.Error("an error occurred while checking wallet balance", slog.Any("err", wallet.ErrInsufficientFunds.Error()))
			return wallet.ErrInsufficientFunds
		}

		updateWallet := `UPDATE wallets SET balance = balance - $1, updated_at = now() WHERE user_id = $2`
		if _, err := tx.ExecContext(ctx, updateWallet, totalPrice, userID); err != nil {
			log.Error("an error occurred while updating wallet balance", slog.Any("err", err))
			return err
		}

		insertPurchase := `INSERT INTO purchases (user_id, item_id, quantity, total_price) VALUES ($1, $2, $3, $4)`
		if _, err := tx.ExecContext(ctx, insertPurchase, userID, itemID, quantity, totalPrice); err != nil {
			log.Error("an error occurred while inserting purchase", slog.Any("err", err))
			return err
		}

		log.Info("purchase completed successfully", slog.Any("user_id", userID), slog.Any("item_id", itemID))
		return nil
	})
}

// GetHistory implements purchase.PurchaseRepo.
func (r *purchaseRepo) GetHistory(ctx context.Context, userID user.UserID) ([]*purchase.Purchase, error) {
	const op = "repository.purchaseRepo.GetHistory"
	log := r.log.With(slog.String("op", op))

	purchases := []*purchase.Purchase{}

	q := `SELECT * FROM purchases WHERE user_id = $1 ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &purchases, q)
	if err != nil {
		log.Error("an error occurred while getting purchase history", slog.Any("err", err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, purchase.ErrHistoryEmpty
		}
		return nil, err
	}

	return purchases, nil
}
