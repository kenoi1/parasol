/*
make watchlist table if name doesnt exist.

id - column name
serial - auto increment int (as rows inserted)
primary key - marks column as unique id for each row

ticker - column name (ex: XEQT)
text - variable string
NOT NULL - can not be missing
UNIQUE - 2 rows cant have same ticker

added_at - records when row was created
TIMESTAMPTZ - timestamp with timezone
NOT NULL - cant missing
DEFAULT now() - if no val provide, default to current timestamp
*/
CREATE TABLE IF NOT EXISTS watchlist ( 
    id SERIAL PRIMARY KEY,
    ticker TEXT NOT NULL UNIQUE,
    added_at TIMESTAMPTZ NOT NULL DEFAULT now()
);