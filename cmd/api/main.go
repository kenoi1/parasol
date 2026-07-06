package main

import (
	"context"
	"fmt"

	"encoding/json" // json for apis
	"log"           // logging
	"net/http"      // http
	"os"            // for env variables (port)
	"time"          //  timeouts for slow clients

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

func main() {
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

	// multiplexer: map url paths to functions
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)

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
