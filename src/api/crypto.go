package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

// generateHandle creates a unique handle from a public key
// Format: <10-char-hash>@authgrid.net
func generateHandle(publicKey []byte) string {
	// Hash the public key
	hash := sha256.Sum256(publicKey)

	// Take first 10 characters of hex encoding for the identifier
	identifier := hex.EncodeToString(hash[:])[:10]

	// Get domain from environment or use default
	domain := getEnv("AUTHGRID_DOMAIN", "authgrid.net")

	return fmt.Sprintf("%s@%s", identifier, domain)
}

// generateChallenge creates a cryptographically random challenge
// Returns base64-encoded 32-byte random string
func generateChallenge() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// generateToken creates a simple token for authenticated sessions
// In production, use proper JWT with signing
func generateToken(handle string) (string, error) {
	// For MVP, generate a random token
	// TODO: Replace with proper JWT implementation
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Include handle in token (simplified approach)
	tokenData := append([]byte(handle+":"), bytes...)
	return base64.StdEncoding.EncodeToString(tokenData), nil
}

// verifySignature verifies a signature using the appropriate algorithm
func verifySignature(publicKeyStr string, keyType string, message []byte, signature []byte) (bool, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return false, fmt.Errorf("invalid public key encoding: %w", err)
	}

	switch keyType {
	case "ed25519":
		// For Ed25519, the public key might be in SPKI format or raw
		var publicKey ed25519.PublicKey

		// Try parsing as SPKI first
		parsedKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
		if err == nil {
			// Successfully parsed as SPKI
			var ok bool
			publicKey, ok = parsedKey.(ed25519.PublicKey)
			if !ok {
				return false, fmt.Errorf("public key is not Ed25519")
			}
		} else {
			// Assume raw Ed25519 key (32 bytes)
			if len(publicKeyBytes) != ed25519.PublicKeySize {
				return false, fmt.Errorf("invalid Ed25519 key length: got %d, want %d", len(publicKeyBytes), ed25519.PublicKeySize)
			}
			publicKey = ed25519.PublicKey(publicKeyBytes)
		}

		return ed25519.Verify(publicKey, message, signature), nil

	case "ecdsa":
		// Parse ECDSA public key from SPKI format
		parsedKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
		if err != nil {
			return false, fmt.Errorf("failed to parse ECDSA public key: %w", err)
		}

		ecdsaKey, ok := parsedKey.(*ecdsa.PublicKey)
		if !ok {
			return false, fmt.Errorf("public key is not ECDSA")
		}

		// Hash the message with SHA-256 (ECDSA requires hashed message)
		hash := sha256.Sum256(message)

		// Parse DER-encoded signature
		var sig struct {
			R, S *big.Int
		}
		_, err = asn1.Unmarshal(signature, &sig)
		if err != nil {
			return false, fmt.Errorf("failed to parse ECDSA signature: %w", err)
		}

		// Verify the signature
		return ecdsa.Verify(ecdsaKey, hash[:], sig.R, sig.S), nil

	default:
		return false, fmt.Errorf("unsupported key type: %s", keyType)
	}
}
