package main

import (
	"context"
	"fmt"

	"encoding/json" // json for apis
	"log"           // logging
	"net/http"      // http
	"os"            // for env variables (port)
	"time"          //  timeouts for slow clients

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/kenoi1/parasol/internal/finnhub"
	"github.com/kenoi1/parasol/internal/store"
)

func Hello(name string) string {
	// return gretting with name in a msg
	message := fmt.Sprintf("hi, %v. welcome to parasol!", name)
	return message
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("ok"))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

/*
wraps an outer functions that takes pool, and returns actual handler to inject dependencies
into handles w/o global vars. cal watchlist handler and it hands back http handler func.

r.Contect() - every incoming req carries context which is tied to req lifecycle.
pass context down into db calls.
*/
func watchlistHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := store.GetWatchlist(r.Context(), pool)
		if err != nil {
			http.Error(w, "failed to fetch watchlist", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	fmt.Println("Hello, World!") // greetings!
	fmt.Println(Hello("derek"))
	// fmt.Println(quote.Opt())

	port := os.Getenv("PORT")
	if port == "" { // fallback to default port
		port = "8080"
	}

	// == init postgres ==
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://parasol:parasol@localhost:5432/parasol"
	}

	pool, err := store.NewPostgresPool(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to create postgres pool: %v", err)
	}
	defer pool.Close()

	if err := store.Ping(ctx, pool); err != nil {
		log.Fatalf("failed to ping postgres: %v", err)
	}
	log.Println("connected to postgres successfully")

	// get finnhub data
	fh := finnhub.NewClient(os.Getenv("FINNHUB_API_KEY"))
	quote, err := fh.GetQuote("TQQQ")
	if err != nil {
		log.Fatalf("failed to get quote: %v", err)
	}
	log.Printf("TQQQ quote: %+v", quote)

	// get stored tickers
	tickers := []string{"TQQQ", "SMH", "XEQT"}
	for _, t := range tickers {
		if err := store.AddTicker(ctx, pool, t); err != nil {
			log.Fatalf("failed to add ticker %s: %v", t, err)
		}
	}

	items, err := store.GetWatchlist(ctx, pool)
	if err != nil {
		log.Fatalf("failed to get watchlist %v", err)
	}
	log.Printf("current watchlist: %+v", items)

	// multiplexer: map url paths to functions
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/watchlist", watchlistHandler(pool))

	srv := &http.Server{ // server start logic
		Addr:         ":" + port,       //  addr to listen on
		Handler:      mux,              // router, maps requests to handler funcs
		ReadTimeout:  5 * time.Second,  // max time for reading inc reqs (header+body)
		WriteTimeout: 10 * time.Second, // max time writing response back to client
	}

	// make barebones server
	log.Printf("starting server on :%s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
