/*
REFERENCES - enforces that you cant insert price row for a ticker not in watchlist.
NUMERIC - no floating point rounding

INDEX - make look up fast, no scanning whole table
*/

CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    ticker TEXT NOT NULL REFERENCES watchlist(ticker),
    price NUMERIC NOT NULL,
    change NUMERIC NOT NULL,
    percent_change NUMERIC NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_prices_ticker_fetched_at ON prices (ticker, fetched_at DESC);