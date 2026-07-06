package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Ping(ctx context.Context, pool *pgxpool.Pool) error {
	return pool.Ping(ctx)
}

/*
pgxpool - pgx conn pool manager
rather than 1 req, maintains pool of reusable connections.

connString - postgres URL "postgres://user:password@host:port/dbname"
*/
func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
