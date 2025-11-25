package utils

import (
	"avito-shop-service/internal/domain/user"
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserFromCtx(gc *gin.Context) (*user.User, error) {
	uctx, ok := gc.Get("user")
	if !ok {
		return nil, errors.New("unauthorized: user not found in context")
	}

	user, ok := uctx.(*user.User)
	if !ok {
		return nil, errors.New("unauthorized: invalid user type")
	}

	return user, nil
}
