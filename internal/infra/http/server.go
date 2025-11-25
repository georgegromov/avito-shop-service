package infrahttp

import (
	"context"
	"errors"
)

type HttpServer interface {
	Start() error
	Stop(ctx context.Context) error
}

var (
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
)
