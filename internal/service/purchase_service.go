package service

import (
	"avito-shop-service/internal/domain/item"
	"avito-shop-service/internal/domain/purchase"
	"avito-shop-service/internal/domain/user"
	"context"
	"log/slog"
)

type purchaseService struct {
	log  *slog.Logger
	repo purchase.PurchaseRepo
}

func NewPurchaseService(
	log *slog.Logger,
	repo purchase.PurchaseRepo,
) purchase.PurchaseService {
	return &purchaseService{
		log:  log,
		repo: repo,
	}
}

// BuyItem implements purchase.PurchaseService.
func (p *purchaseService) BuyItem(ctx context.Context, userID user.UserID, itemID item.ItemID, quantity uint) error {
	return p.repo.BuyItem(ctx, userID, itemID, quantity)
}

// GetHistory implements purchase.PurchaseService.
func (p *purchaseService) GetHistory(ctx context.Context, userID user.UserID) ([]*purchase.Purchase, error) {
	return p.repo.GetHistory(ctx, userID)
}
