package main

import (
	"car-backend/pkg/config"
	"car-backend/pkg/handlers"
	"car-backend/pkg/repository"
	"context"
	"log"
	"net/http"

	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func initDB(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	// At the start, configure log to use JSON format
	log.SetFlags(0)          // Remove timestamp prefix
	log.SetOutput(os.Stdout) // Ensure logs go to stdout

	// Add environment variable check with better logging
	requiredEnvVars := []string{
		"INSTANCE_CONNECTION_NAME",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"CLERK_SECRET_KEY",
	}

	log.Printf("{\"severity\":\"INFO\",\"message\":\"Checking required environment variables\"}")

	for _, envVar := range requiredEnvVars {
		value := os.Getenv(envVar)
		if value == "" {
			log.Printf("{\"severity\":\"ERROR\",\"message\":\"Missing required environment variable: %s\"}", envVar)
			// Print all environment variables for debugging (excluding sensitive ones)
			for _, env := range os.Environ() {
				if !strings.Contains(env, "PASSWORD") && !strings.Contains(env, "SECRET") {
					log.Printf("{\"severity\":\"DEBUG\",\"message\":\"Environment: %s\"}", env)
				}
			}
			os.Exit(1) // Change from return to os.Exit(1)
		}
		if envVar != "DB_PASSWORD" && envVar != "CLERK_SECRET_KEY" {
			log.Printf("{\"severity\":\"INFO\",\"message\":\"Found environment variable %s: %s\"}", envVar, value)
		} else {
			log.Printf("{\"severity\":\"INFO\",\"message\":\"Found environment variable %s: [REDACTED]\"}", envVar)
		}
	}

	log.Printf("{\"severity\":\"INFO\",\"message\":\"Starting server\"}")
	// Load configuration
	cfg := config.Load()
	log.Printf("{\"severity\":\"INFO\",\"message\":\"Starting server on port %s\"}", cfg.Port)

	// Initialize DB connection with Cloud SQL
	var dbURI string
	if os.Getenv("INSTANCE_CONNECTION_NAME") != "" {
		// Cloud SQL connection format should be:
		dbURI = fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("INSTANCE_CONNECTION_NAME"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
	} else {
		// Local development connection
		dbURI = cfg.DatabaseURL
	}

	log.Printf("{\"severity\":\"INFO\",\"message\":\"Attempting database connection with instance: %s\"}",
		os.Getenv("INSTANCE_CONNECTION_NAME"))

	db, err := sql.Open("cloudsqlpostgres", dbURI)
	if err != nil {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"Database connection failed: %v\"}", err)
		os.Exit(1) // Change from return to os.Exit(1)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"Failed to ping database: %v\"}", err)
		os.Exit(1) // Change from return to os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo)

	// Setup router
	r := mux.NewRouter()

	// Add welcome message for root path
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Welcome to Car Backend API!",
		})
	}).Methods("GET")

	// Public routes
	r.HandleFunc("/webhook/clerk", userHandler.HandleWebhook).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(clerkhttp.RequireHeaderAuthorization())

	protected.HandleFunc("/profile", userHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/profile", userHandler.UpdateProfile).Methods("PUT")

	// Add health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
		})
	}).Methods("GET")

	// Create server with proper port handling
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 instead of using config
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Add more detailed startup logging
	log.Printf("{\"severity\":\"INFO\",\"message\":\"Server starting on port %s\"}", port)

	// Use a channel to handle shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("{\"severity\":\"ERROR\",\"message\":\"HTTP server Shutdown: %v\"}", err)
		}
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"HTTP server ListenAndServe: %v\"}", err)
		os.Exit(1)
	}

	<-idleConnsClosed
}
