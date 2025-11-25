package transfer

import (
	"avito-shop-service/internal/domain/user"
	"avito-shop-service/internal/domain/wallet"
	"context"

	"github.com/gin-gonic/gin"
)

type TransferRepo interface {
	SendCoins(ctx context.Context, fromID, toID user.UserID, amount wallet.Amount) error
	GetHistory(ctx context.Context, userID user.UserID) ([]*Transfer, error)
}

type TransferService interface {
	SendCoins(ctx context.Context, fromID, toID user.UserID, amount wallet.Amount) error
	GetHistory(ctx context.Context, userID user.UserID) ([]*Transfer, error)
}

type TransferHandler interface {
	SendCoinsRoute(*gin.Context)
	GetHistoryRoute(*gin.Context)
}
