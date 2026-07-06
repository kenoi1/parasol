package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WatchlistItem struct {
	ID     int
	Ticker string
}

// add an individual value pair (id, ticker)
func AddTicker(ctx context.Context, pool *pgxpool.Pool, ticker string) error {
	// $1 - placeholder for ticker arg.
	// ON CONFLICT - ticker is unique. skip if added
	_, err := pool.Exec(ctx,
		`INSERT INTO watchlist (ticker) VALUES ($1) ON CONFLICT (ticker) DO NOTHING`,
		ticker,
	)
	return err
}

func GetWatchlist(ctx context.Context, pool *pgxpool.Pool) ([]WatchlistItem, error) {
	// return rows, which is iterated over with rows.Next()
	// rows.Scan() copies cols into Go vars, in SELECTed order
	rows, err := pool.Query(ctx, `SELECT id, ticker FROM watchlist ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []WatchlistItem
	for rows.Next() { // store as go vars
		var item WatchlistItem
		if err := rows.Scan(&item.ID, &item.Ticker); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
