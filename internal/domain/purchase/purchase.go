package purchase

import (
	"avito-shop-service/internal/domain/item"
	"avito-shop-service/internal/domain/user"
	"time"

	"github.com/google/uuid"
)

type PurchaseID uuid.UUID

type Purchase struct {
	ID         PurchaseID
	UserID     user.UserID
	ItemID     item.ItemID
	Quantity   uint
	TotalPrice item.Price
	CreatedAt  time.Time
}

func New(userID user.UserID, itemID item.ItemID, quantity uint, price item.Price) (*Purchase, error) {
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	return &Purchase{
		UserID:     userID,
		ItemID:     itemID,
		Quantity:   quantity,
		TotalPrice: price,
		CreatedAt:  time.Now(),
	}, nil
}
