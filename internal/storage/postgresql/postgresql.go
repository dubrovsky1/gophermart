package postgresql

import (
	"context"
	"embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"time"
)

//go:embed migrations
var migrations embed.FS

type Storage struct {
	pool *pgxpool.Pool
}

func (s *Storage) Close() error {
	s.pool.Close()
	return nil
}

func New(connectionString string, connTimeout time.Duration, version int64) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}

	if err = migrate(pool, version); err != nil {
		return nil, err
	}

	return &Storage{pool: pool}, nil
}

func migrate(pool *pgxpool.Pool, version int64) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("postgres migrate set dialect postgres: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.UpTo(db, "migrations", version); err != nil {
		return fmt.Errorf("postgres migrate up: %w", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("postgres migrate close db: %w", err)
	}
	return nil
}
