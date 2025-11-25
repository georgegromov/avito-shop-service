package purchase

import "errors"

var (
	ErrInvalidQuantity = errors.New("invalid quantity")
	ErrHistoryEmpty    = errors.New("purchase history is empty")
)
