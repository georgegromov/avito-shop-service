package wallet

import (
	"avito-shop-service/internal/domain/user"
	"context"

	"github.com/jmoiron/sqlx"
)

type WalletRepo interface {
	GetByUserID(ctx context.Context, userID user.UserID) (*Wallet, error)
	UpdateBalanceTx(ctx context.Context, tx sqlx.Tx, wallet *Wallet) error
}

type WalletService interface {
	Increase(ctx context.Context, userID user.UserID, amount Amount) error
	Decrease(ctx context.Context, userID user.UserID, amount Amount) error
}
