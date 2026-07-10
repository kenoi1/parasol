package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
populate ingested finnhub data in database
*/
func SavePrice(ctx context.Context, pool *pgxpool.Pool, ticker string, price, change, percentChange float64) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO prices (ticker, price, change, percent_change) VALUES ($1, $2, $3, $4)`,
		ticker, price, change, percentChange,
	)
	return err
}
