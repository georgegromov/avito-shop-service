package jwt

import (
	"avito-shop-service/internal/domain/user"
	"avito-shop-service/internal/infra/security"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	secret   []byte
	tokenTTL time.Duration
}

func New(secret string, tokenTTL time.Duration) *JwtManager {
	return &JwtManager{
		secret:   []byte(secret),
		tokenTTL: tokenTTL,
	}
}

func (j *JwtManager) Generate(userID user.UserID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JwtManager) Validate(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, security.ErrInvalidSigningMethod
		}
		return j.secret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
	}

	return "", security.ErrInvalidToken
}

var _ security.JwtManager = (*JwtManager)(nil)
