package main

import (
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// RegisterRequest represents a registration request
type RegisterRequest struct {
	PublicKey string `json:"public_key"` // base64 encoded
	KeyType   string `json:"key_type"`   // "ed25519"
}

// RegisterResponse represents a registration response
type RegisterResponse struct {
	Handle    string    `json:"handle"`
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// ChallengeRequest represents a challenge request
type ChallengeRequest struct {
	Handle string `json:"handle"`
}

// ChallengeResponse represents a challenge response
type ChallengeResponse struct {
	Challenge string    `json:"challenge"`
	ExpiresAt time.Time `json:"expires_at"`
}

// VerifyRequest represents a verification request
type VerifyRequest struct {
	Handle    string `json:"handle"`
	Challenge string `json:"challenge"`
	Signature string `json:"signature"`
}

// VerifyResponse represents a verification response
type VerifyResponse struct {
	Verified  bool      `json:"verified"`
	Token     string    `json:"token,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// registerHandler handles user registration
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate key type
	if req.KeyType != "ed25519" && req.KeyType != "ecdsa" {
		respondError(w, http.StatusBadRequest, "Only ed25519 and ecdsa key types are supported")
		return
	}

	// Decode public key
	publicKeyBytes, err := base64.StdEncoding.DecodeString(req.PublicKey)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid public key encoding")
		return
	}

	// Validate public key length based on key type
	if req.KeyType == "ed25519" {
		// Ed25519 raw public keys are 32 bytes
		if len(publicKeyBytes) != ed25519.PublicKeySize && len(publicKeyBytes) < 30 {
			respondError(w, http.StatusBadRequest, "Invalid Ed25519 public key length")
			return
		}
	} else if req.KeyType == "ecdsa" {
		// ECDSA P-256 public keys in SPKI format are ~91 bytes
		// We accept keys between 60-120 bytes for ECDSA
		if len(publicKeyBytes) < 60 || len(publicKeyBytes) > 120 {
			respondError(w, http.StatusBadRequest, "Invalid ECDSA public key length")
			return
		}
	}

	// Generate handle from public key
	handle := generateHandle(publicKeyBytes)

	// Check if handle already exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE handle = $1)", handle).Scan(&exists)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if exists {
		respondError(w, http.StatusConflict, "Handle already exists")
		return
	}

	// Insert user into database
	var id string
	var createdAt time.Time
	err = db.QueryRow(`
		INSERT INTO users (handle, public_key, key_type, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`, handle, req.PublicKey, req.KeyType).Scan(&id, &createdAt)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondJSON(w, http.StatusCreated, RegisterResponse{
		Handle:    handle,
		ID:        id,
		CreatedAt: createdAt,
	})
}

// challengeHandler generates an authentication challenge
func challengeHandler(w http.ResponseWriter, r *http.Request) {
	var req ChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Handle == "" {
		respondError(w, http.StatusBadRequest, "Handle is required")
		return
	}

	// Check if user exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE handle = $1)", req.Handle).Scan(&exists)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if !exists {
		respondError(w, http.StatusNotFound, "Handle not found")
		return
	}

	// Generate challenge
	challenge, err := generateChallenge()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate challenge")
		return
	}

	expiresAt := time.Now().Add(5 * time.Minute)

	// Store challenge in database
	_, err = db.Exec(`
		INSERT INTO challenges (handle, challenge, created_at, expires_at, used)
		VALUES ($1, $2, NOW(), $3, FALSE)
	`, req.Handle, challenge, expiresAt)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to store challenge")
		return
	}

	respondJSON(w, http.StatusOK, ChallengeResponse{
		Challenge: challenge,
		ExpiresAt: expiresAt,
	})
}

// verifyHandler verifies a signed challenge
func verifyHandler(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Handle == "" || req.Challenge == "" || req.Signature == "" {
		respondError(w, http.StatusBadRequest, "Handle, challenge, and signature are required")
		return
	}

	// Get user's public key and key type
	var publicKeyStr string
	var keyType string
	err := db.QueryRow("SELECT public_key, key_type FROM users WHERE handle = $1", req.Handle).Scan(&publicKeyStr, &keyType)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Handle not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Check if challenge exists and is valid
	var challengeID string
	var expiresAt time.Time
	var used bool
	err = db.QueryRow(`
		SELECT id, expires_at, used
		FROM challenges
		WHERE handle = $1 AND challenge = $2
	`, req.Handle, req.Challenge).Scan(&challengeID, &expiresAt, &used)

	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Challenge not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Check if challenge is expired
	if time.Now().After(expiresAt) {
		respondError(w, http.StatusBadRequest, "Challenge expired")
		return
	}

	// Check if challenge was already used
	if used {
		respondError(w, http.StatusBadRequest, "Challenge already used")
		return
	}

	// Decode challenge and signature
	challengeBytes, err := base64.StdEncoding.DecodeString(req.Challenge)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid challenge encoding")
		return
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid signature encoding")
		return
	}

	// Verify signature based on key type
	valid, err := verifySignature(publicKeyStr, keyType, challengeBytes, signatureBytes)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Signature verification error: "+err.Error())
		return
	}
	if !valid {
		respondError(w, http.StatusUnauthorized, "Invalid signature")
		return
	}

	// Mark challenge as used
	_, err = db.Exec("UPDATE challenges SET used = TRUE WHERE id = $1", challengeID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to mark challenge as used")
		return
	}

	// Update last login time
	_, err = db.Exec("UPDATE users SET last_login = NOW() WHERE handle = $1", req.Handle)
	if err != nil {
		// Non-critical error, log but continue
		// In production, use proper logging
	}

	// Generate token (simplified - in production use proper JWT)
	token, err := generateToken(req.Handle)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	tokenExpiry := time.Now().Add(24 * time.Hour)

	respondJSON(w, http.StatusOK, VerifyResponse{
		Verified:  true,
		Token:     token,
		ExpiresAt: tokenExpiry,
	})
}

// getUserHandler returns public user information
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	handle := vars["handle"]

	var publicKey string
	var createdAt time.Time
	err := db.QueryRow(`
		SELECT public_key, created_at
		FROM users
		WHERE handle = $1
	`, handle).Scan(&publicKey, &createdAt)

	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Handle not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"handle":     handle,
		"public_key": publicKey,
		"created_at": createdAt,
	})
}
