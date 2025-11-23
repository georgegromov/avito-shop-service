package item

import "github.com/google/uuid"

type ItemID uuid.UUID
type Price uint64

type Item struct {
	ID    ItemID
	Name  string
	Price Price
}

func New(name string, price Price) (*Item, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	if price <= 0 {
		return nil, ErrInvalidPrice
	}

	return &Item{
		Name:  name,
		Price: price,
	}, nil
}
