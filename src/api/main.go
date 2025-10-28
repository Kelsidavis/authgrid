package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

var (
	db      *sql.DB
	limiter *rate.Limiter
)

func main() {
	// Initialize database
	var err error
	dbURL := getEnv("DATABASE_URL", "postgres://authgrid:authgrid@localhost:5432/authgrid?sslmode=disable")
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Successfully connected to database")

	// Initialize rate limiter (100 requests per second)
	limiter = rate.NewLimiter(100, 200)

	// Setup router
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", healthHandler).Methods("GET")

	// Core endpoints
	r.HandleFunc("/register", rateLimitMiddleware(registerHandler)).Methods("POST")
	r.HandleFunc("/challenge", rateLimitMiddleware(challengeHandler)).Methods("POST")
	r.HandleFunc("/verify", rateLimitMiddleware(verifyHandler)).Methods("POST")

	// User lookup (optional, for public key retrieval)
	r.HandleFunc("/user/{handle}", getUserHandler).Methods("GET")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // In production, specify exact origins
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(r)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Authgrid API server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// rateLimitMiddleware applies rate limiting to endpoints
func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			respondError(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}
		next(w, r)
	}
}

// healthHandler returns server health status
func healthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
