package item

import "errors"

var (
	ErrInvalidName  = errors.New("invalid name")
	ErrInvalidPrice = errors.New("invalid price")
	ErrItemNotFound = errors.New("item not found")
)
