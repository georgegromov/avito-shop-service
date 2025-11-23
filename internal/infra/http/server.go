package http

import "context"

type HttpServer interface {
	Start() error
	Stop(ctx context.Context) error
}
