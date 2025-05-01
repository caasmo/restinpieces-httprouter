package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/caasmo/restinpieces"
	"github.com/caasmo/restinpieces-httprouter"
	r "github.com/caasmo/restinpieces/router"
)

// Simple logging middleware example
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request received", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
		// We could log response status here if needed, but requires a ResponseWriter wrapper
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func echoHandler(w http.ResponseWriter, req *http.Request) {
	message := req.URL.Query().Get("msg")
	if message == "" {
		message = "Echo!"
	}
	fmt.Fprintln(w, message)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// In a real app, you'd read the request body
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Resource created (simulated)")
}

func main() {
	dbPath := flag.String("dbpath", "", "Path to the SQLite database file (required)") // Changed flag name
	ageKeyPath := flag.String("age-key", "", "Path to the age identity (private key) file (required)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -dbpath <database-path> -age-key <identity-file-path>\n\n", os.Args[0]) // Changed usage string
		fmt.Fprintf(os.Stderr, "Start the restinpieces example application server using httprouter.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *dbPath == "" || *ageKeyPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	dbPool, err := restinpieces.NewZombiezenPool(*dbPath)
	if err != nil {
		slog.Error("failed to create database pool", "error", err)
		os.Exit(1)
	}

	defer func() {
		slog.Info("Closing database pool...")
		if err := dbPool.Close(); err != nil {
			slog.Error("Error closing database pool", "error", err)
		}
	}()

	app, srv, err := restinpieces.New(
		restinpieces.WithDbZombiezen(dbPool),
		restinpieces.WithAgeKeyPath(*ageKeyPath),
		httprouter.WithRouterHttprouter(),
		restinpieces.WithCacheRistretto(),
		restinpieces.WithTextLogger(nil),
	)
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Register example routes
	app.Router().Register(map[string]*r.Chain{
		"GET /hello": r.NewChain(http.HandlerFunc(helloHandler)).WithMiddleware(loggingMiddleware), // Added middleware
		"/echo":      r.NewChain(http.HandlerFunc(echoHandler)),                                     // Defaults to GET
		"POST /items": r.NewChain(http.HandlerFunc(postHandler)),
	})

	srv.Run()

	slog.Info("Server shut down gracefully.")
}
