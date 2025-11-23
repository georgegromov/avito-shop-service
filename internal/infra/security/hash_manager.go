package security

import "errors"

type HashManager interface {
	Hash(password string) (string, error)
	CompareHash(password, hash string) error
}

var (
	ErrEmptyPassword = errors.New("password is empty")
)
