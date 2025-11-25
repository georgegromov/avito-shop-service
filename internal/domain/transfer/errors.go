package transfer

import "errors"

var (
	ErrHistoryEmpty         = errors.New("transfer history is empty")
	ErrCannotSendToYourself = errors.New("cannot send coins to yourself")
	ErrInvalidAmount        = errors.New("invalid amount")
	ErrInvalidUser          = errors.New("invalid user")
)
