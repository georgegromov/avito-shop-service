package wallet

import (
	"avito-shop-service/internal/domain/user"
	"time"
)

type Amount uint64

type Wallet struct {
	UserID    user.UserID
	Balance   Amount
	UpdatedAt time.Time
}

func (w *Wallet) CanSpend(amount Amount) bool {
	return amount > 0 && w.Balance >= amount
}

func (w *Wallet) Spend(amount Amount) error {
	if !w.CanSpend(amount) {
		return ErrInsufficientFunds
	}
	w.Balance -= amount
	w.UpdatedAt = time.Now()
	return nil
}

func (w *Wallet) Gain(amount Amount) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	w.Balance += amount
	w.UpdatedAt = time.Now()
	return nil
}
