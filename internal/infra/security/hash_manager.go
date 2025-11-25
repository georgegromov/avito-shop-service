package security

import "errors"

type HashManager interface {
	Hash(password string) (string, error)
	CompareHash(password, hash string) error
}

var (
	ErrHashPassword  = errors.New("an error occurred while hashing password")
	ErrEmptyPassword = errors.New("password is empty")
)
