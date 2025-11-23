package postgres

import (
	"avito-shop-service/pkg/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

func MustConnect(log *slog.Logger, cfg config.PostgresConfig) *sqlx.DB {
	db, err := Connect(log, cfg)
	if err != nil {
		panic(err)
	}
	return db
}

func Connect(log *slog.Logger, cfg config.PostgresConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Таймаут для пинга базы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Info("connected to PostgreSQL", "host", cfg.Host, "db", cfg.DBName)
	return db, nil
}
