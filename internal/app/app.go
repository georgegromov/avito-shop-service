package app

import (
	"avito-shop-service/internal/handler"
	"avito-shop-service/internal/infra/db/postgres"
	"avito-shop-service/internal/repository"
	"avito-shop-service/internal/service"
	"avito-shop-service/pkg/config"
	"avito-shop-service/pkg/hash"
	"avito-shop-service/pkg/jwt"
	"avito-shop-service/pkg/logger"
	"avito-shop-service/pkg/server"
	"avito-shop-service/pkg/validator"
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type App struct {
	Cfg *config.Config
	log *slog.Logger

	db *sqlx.DB

	HttpServer *server.HttpServer
}

func New() *App {
	cfg := config.MustLoad()
	logger := logger.New(cfg.Env)

	db := postgres.MustConnect(logger, cfg.Postgres)

	// if ENV == 'development' when seed items

	hashManager := hash.New(int(cfg.Hash.Cost))
	jwtManager := jwt.New(cfg.Jwt.Secret, cfg.Jwt.TokenTTL)

	v := validator.New()

	// user module
	userRepo := repository.NewUserRepo(logger, db)
	userService := service.NewUserService(logger, userRepo, hashManager, jwtManager)
	userHandler := handler.NewUserHandler(logger, v, userService, jwtManager)

	// purchase module
	purchaseRepo := repository.NewPurchaseRepo(logger, db)
	purchaseService := service.NewPurchaseService(logger, purchaseRepo)
	purchaseHandler := handler.NewPurchaseHandler(logger, v, purchaseService)

	// transfer module
	transferRepo := repository.NewTransferRepo(logger, db)
	transferService := service.NewTransferService(logger, transferRepo)
	transferHandler := handler.NewTransferHandler(logger, v, transferService)

	h := handler.NewHandler(logger, userHandler, purchaseHandler, transferHandler)
	engine := h.NewEngine()

	srv := server.New(logger, &cfg.HttpServer, engine)

	return &App{
		Cfg:        cfg,
		log:        logger,
		db:         db,
		HttpServer: srv,
	}
}

func (a *App) Close(ctx context.Context) {
	const op = "app.App.Close"
	log := a.log.With(slog.String("op", op))

	if err := a.HttpServer.Stop(ctx); err != nil {
		log.Error("an error occurred while stopping http server", slog.Any("err", err))
	}

	if err := a.db.Close(); err != nil {
		log.Error("an error occurred while closing database", slog.Any("err", err))
	}

	log.Info("application gracefully stopped")
}
