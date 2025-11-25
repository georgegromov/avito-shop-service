package server

import (
	infrahttp "avito-shop-service/internal/infra/http"
	"avito-shop-service/pkg/config"
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type HttpServer struct {
	log        *slog.Logger
	config     *config.HttpServerConfig
	httpServer *http.Server
}

func New(log *slog.Logger, cfg *config.HttpServerConfig, handler http.Handler) *HttpServer {

	srv := &http.Server{
		Addr:           cfg.Host + ":" + cfg.Port,
		Handler:        handler,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		IdleTimeout:    cfg.IdleTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	return &HttpServer{
		log:        log,
		config:     cfg,
		httpServer: srv,
	}
}

func (h *HttpServer) MustStart() {
	const op = "server.HttpServer.MustStart"
	if err := h.Start(); err != nil {
		panic(fmt.Sprintf("%s: an error occurred: %v", op, err))
	}
}

func (h *HttpServer) Start() error {
	const op = "server.HttpServer.Start"

	log := h.log.With(slog.String("op", op))

	log.Info("starting http server...", slog.String("addr", h.config.Host+":"+h.config.Port))

	if err := h.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (h *HttpServer) Stop(ctx context.Context) error {
	const op = "server.HttpServer.Stop"

	log := h.log.With(slog.String("op", op))

	log.Info("stopping http server...")

	return h.httpServer.Shutdown(ctx)
}

var _ infrahttp.HttpServer = (*HttpServer)(nil)
