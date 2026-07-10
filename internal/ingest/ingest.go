package ingest

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kenoi1/parasol/internal/finnhub"
	"github.com/kenoi1/parasol/internal/store"
)

/*
Retrieve stocks data from finnhub of the stocks in watchlist
*/
func RunOnce(ctx context.Context, pool *pgxpool.Pool, fh *finnhub.Client) error {
	items, err := store.GetWatchlist(ctx, pool)
	if err != nil {
		return err
	}

	for _, item := range items {
		quote, err := fh.GetQuote(item.Ticker)
		if err != nil {
			log.Printf("failed to fetch quote for %s: %v", item.Ticker, err)
			continue
		}

		err = store.SavePrice(ctx, pool, item.Ticker, quote.CurrentPrice, quote.Change, quote.PercentChange)
		if err != nil {
			log.Printf("failed to save price for %s: %v", item.Ticker, err)
			continue
		}

		log.Printf("saved price for %s: %.2f", item.Ticker, quote.CurrentPrice)

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
