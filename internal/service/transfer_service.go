package service

import (
	"avito-shop-service/internal/domain/transfer"
	"avito-shop-service/internal/domain/user"
	"avito-shop-service/internal/domain/wallet"
	"context"
	"log/slog"
)

type transferService struct {
	log  *slog.Logger
	repo transfer.TransferRepo
}

func NewTransferService(
	log *slog.Logger,
	repo transfer.TransferRepo,
) transfer.TransferService {
	return &transferService{
		log:  log,
		repo: repo,
	}
}

// SendCoins implements transfer.TransferService.
func (t *transferService) SendCoins(ctx context.Context, fromID user.UserID, toID user.UserID, amount wallet.Amount) error {
	return t.repo.SendCoins(ctx, fromID, toID, amount)
}

// GetHistory implements transfer.TransferService.
func (t *transferService) GetHistory(ctx context.Context, userID user.UserID) ([]*transfer.Transfer, error) {
	return t.repo.GetHistory(ctx, userID)
}
