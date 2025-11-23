package wallet

import "errors"

var (
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("invalid amount")
)
