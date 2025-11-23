package transfer

import "errors"

var (
	ErrInvalidAmount = errors.New("invalid amount")
	ErrInvalidUser   = errors.New("invalid user")
)
