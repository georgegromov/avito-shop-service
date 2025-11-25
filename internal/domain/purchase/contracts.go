package purchase

import (
	"avito-shop-service/internal/domain/item"
	"avito-shop-service/internal/domain/user"
	"context"

	"github.com/gin-gonic/gin"
)

type PurchaseRepo interface {
	BuyItem(ctx context.Context, userID user.UserID, itemID item.ItemID, quantity uint) error
	GetHistory(ctx context.Context, userID user.UserID) ([]*Purchase, error)
}

type PurchaseService interface {
	BuyItem(ctx context.Context, userID user.UserID, itemID item.ItemID, quantity uint) error
	GetHistory(ctx context.Context, userID user.UserID) ([]*Purchase, error)
}

type PurchaseHandler interface {
	BuyItemRoute(*gin.Context)
	GetHistoryRoute(*gin.Context)
}
