package transfer

import (
	"avito-shop-service/internal/domain/user"
	"time"

	"github.com/google/uuid"
)

type TransferID uuid.UUID

type Transfer struct {
	ID         TransferID
	FromUserID user.UserID
	ToUserID   user.UserID
	Amount     uint64
	CreatedAt  time.Time
}

func New(fromUserID, toUserID user.UserID, amount uint64) (*Transfer, error) {
	if uuid.UUID(fromUserID) == uuid.Nil || uuid.UUID(toUserID) == uuid.Nil {
		return nil, ErrInvalidUser
	}

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	return &Transfer{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
		CreatedAt:  time.Now(),
	}, nil
}
