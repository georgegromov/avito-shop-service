package hash

import (
	"avito-shop-service/internal/infra/security"

	"golang.org/x/crypto/bcrypt"
)

type HashManager struct {
	cost int
}

func New(cost int) *HashManager {
	// default bcrypt cost 10–14.
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}

	return &HashManager{cost: cost}
}

func (m *HashManager) Hash(password string) (string, error) {
	if password == "" {
		return "", security.ErrEmptyPassword
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), m.cost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (m *HashManager) CompareHash(password, hash string) error {
	if password == "" {
		return security.ErrEmptyPassword
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

var _ security.HashManager = (*HashManager)(nil)
