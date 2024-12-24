package main

import (
	"car-backend/pkg/handlers"
	"car-backend/pkg/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gorilla/mux"
)

func setupLogging() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

func checkEnvironment() {
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
			os.Exit(1)
		}
		if envVar != "DB_PASSWORD" && envVar != "CLERK_SECRET_KEY" {
			log.Printf("{\"severity\":\"INFO\",\"message\":\"Found environment variable %s: %s\"}", envVar, value)
		} else {
			log.Printf("{\"severity\":\"INFO\",\"message\":\"Found environment variable %s: [REDACTED]\"}", envVar)
		}
	}
}

func setupDatabase() *sql.DB {
	dbURI := fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("INSTANCE_CONNECTION_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("cloudsqlpostgres", dbURI)
	if err != nil {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"Database connection failed: %v\"}", err)
		os.Exit(1)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}

func setupRouter(userHandler *handlers.UserHandler) *mux.Router {
	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	// Welcome endpoint
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

	return r
}

func main() {
	setupLogging()
	checkEnvironment()

	db := setupDatabase()
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)
	router := setupRouter(userHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown setup
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("{\"severity\":\"ERROR\",\"message\":\"HTTP server Shutdown: %v\"}", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("{\"severity\":\"INFO\",\"message\":\"Server starting on port %s\"}", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"HTTP server ListenAndServe: %v\"}", err)
		os.Exit(1)
	}

	<-idleConnsClosed
}
