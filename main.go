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
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var debugMode bool

func debugLog(format string, v ...interface{}) {
	if debugMode {
		log.Printf("{\"severity\":\"DEBUG\",\"message\":\""+format+"\"}", v...)
	}
}

func setupLogging() {
	// Configure logging to write to stdout and include timestamp
	log.SetOutput(os.Stdout)
	log.SetFlags(0) // Remove timestamp prefix as we're using structured logging
}

func checkEnvironment() {
	var requiredEnvVars []string

	if os.Getenv("ENV") == "local" {
		requiredEnvVars = []string{
			"DB_USER",
			"DB_PASSWORD",
			"DB_NAME",
			"CLERK_SECRET_KEY",
		}
	} else {
		requiredEnvVars = []string{
			"INSTANCE_CONNECTION_NAME",
			"DB_USER",
			"DB_PASSWORD",
			"DB_NAME",
			"CLERK_SECRET_KEY",
		}
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
	var db *sql.DB
	var err error

	if os.Getenv("ENV") == "local" {
		// Local development connecting to Cloud SQL
		dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
		db, err = sql.Open("postgres", dbURI)
	} else {
		// Cloud SQL connection using unix socket
		dbURI := fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s",
			os.Getenv("INSTANCE_CONNECTION_NAME"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)

		debugLog("Attempting to connect to Cloud SQL with connection name: %s", dbURI)
		db, err = sql.Open("postgres", dbURI)
	}

	if err != nil {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"Database connection failed: %v\"}", err)
		os.Exit(1)
	} else {
		log.Printf("{\"severity\":\"INFO\",\"message\":\"Database connection successful\"}")
	}

	// Test the connection
	if err := db.PingContext(context.Background()); err != nil {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"Database ping failed: %v\"}", err)
		os.Exit(1)
	} else {
		log.Printf("{\"severity\":\"INFO\",\"message\":\"Database ping successful\"}")
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}

func setupClerk() {
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		log.Printf("{\"severity\":\"ERROR\",\"message\":\"CLERK_SECRET_KEY is required\"}")
		os.Exit(1)
	}
	clerk.SetKey(clerkSecretKey)
}

func setupRouter(userHandler *handlers.UserHandler, carpoolHandler *handlers.CarPoolHandler, inviteHandler *handlers.InviteHandler, carpoolRideHandler *handlers.CarPoolRideHandler) *mux.Router {
	r := mux.NewRouter()

	// Health check endpoint (public)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	// Welcome endpoint (public)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Welcome to Car Backend API!",
		})
	}).Methods("GET")

	// Public routes
	r.HandleFunc("/webhook/clerk", userHandler.HandleWebhook).Methods("POST")

	// Protected routes with Clerk authentication
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(clerkhttp.WithHeaderAuthorization())
	protected.HandleFunc("/profile", userHandler.CreateProfile).Methods("POST")
	protected.HandleFunc("/profile", userHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/profile", userHandler.UpdateProfile).Methods("PUT")

	// Make these carpool endpoints public for testing
	r.HandleFunc("/api/carpools", carpoolHandler.CreateCarPool).Methods("POST")
	r.HandleFunc("/api/carpools/{id}", carpoolHandler.GetCarPool).Methods("GET")

	//protected.HandleFunc("/carpools", carpoolHandler.CreateCarPool).Methods("POST")
	//protected.HandleFunc("/carpools/{id}", carpoolHandler.GetCarPool).Methods("GET")
	protected.HandleFunc("/carpools/{id}", carpoolHandler.UpdateCarPool).Methods("PUT")
	protected.HandleFunc("/carpools/{id}", carpoolHandler.DeleteCarPool).Methods("DELETE")
	protected.HandleFunc("/carpools/search", carpoolHandler.SearchCarPools).Methods("POST")

	protected.HandleFunc("/carpools/{id}/rides", carpoolRideHandler.CreateCarpoolRide).Methods("POST")
	protected.HandleFunc("/carpools/{id}/rides/{rideID}", carpoolRideHandler.GetCarpoolRide).Methods("GET")

	protected.HandleFunc("/invites", inviteHandler.CreateInvite).Methods("POST")
	protected.HandleFunc("/invites/{id}", inviteHandler.GetInvite).Methods("GET")
	return r
}

func main() {

	setupLogging()

	log.Printf("{\"severity\":\"INFO\",\"message\":\"Starting application\"}")

	// Load .env file before any other setup
	if err := godotenv.Load(); err != nil {
		log.Printf("{\"severity\":\"WARNING\",\"message\":\"Error loading .env file: %v\"}", err)
	} else {
		log.Printf("{\"severity\":\"INFO\",\"message\":\"Successfully loaded .env file\"}")
	}

	// Initialize debug mode
	debugMode = os.Getenv("DEBUG") == "true"
	debugLog("Debug mode enabled")
	debugLog("Environment: %s", os.Getenv("ENV"))
	debugLog("DB_USER: %s", os.Getenv("DB_USER"))

	// CORS middleware configuration
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		Debug:            debugMode,
	})

	checkEnvironment()
	setupClerk()

	db := setupDatabase()
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	carpoolRepo := repository.NewCarPoolRepository(db)
	inviteRepo := repository.NewInviteRepository(db)
	carpoolRideRepo := repository.NewCarPoolRideRepository(db)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo)
	carpoolHandler := handlers.NewCarPoolHandler(carpoolRepo)
	inviteHandler := handlers.NewInviteHandler(inviteRepo)
	carpoolRideHandler := handlers.NewCarPoolRideHandler(carpoolRideRepo)

	router := setupRouter(userHandler, carpoolHandler, inviteHandler, carpoolRideHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create the server with CORS-enabled handler
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware.Handler(router),
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
