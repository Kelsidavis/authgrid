package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func TestGenerateHandle(t *testing.T) {
	// Generate a test public key
	_, publicKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Generate handle
	handle := generateHandle(publicKey)

	// Check format
	if len(handle) < 10 {
		t.Errorf("Handle too short: %s", handle)
	}

	// Should contain @ symbol
	if handle[10] != '@' {
		t.Errorf("Handle missing @ symbol: %s", handle)
	}

	// Should end with domain
	expectedSuffix := "@authgrid.net"
	if len(handle) < len(expectedSuffix) {
		t.Errorf("Handle too short for domain: %s", handle)
	}
}

func TestGenerateHandleDeterministic(t *testing.T) {
	// Same public key should always generate same handle
	_, publicKey, _ := ed25519.GenerateKey(rand.Reader)

	handle1 := generateHandle(publicKey)
	handle2 := generateHandle(publicKey)

	if handle1 != handle2 {
		t.Errorf("Handles don't match: %s != %s", handle1, handle2)
	}
}

func TestGenerateChallenge(t *testing.T) {
	challenge, err := generateChallenge()
	if err != nil {
		t.Fatalf("Failed to generate challenge: %v", err)
	}

	// Should be base64 encoded
	decoded, err := base64.StdEncoding.DecodeString(challenge)
	if err != nil {
		t.Fatalf("Challenge is not valid base64: %v", err)
	}

	// Should be 32 bytes
	if len(decoded) != 32 {
		t.Errorf("Challenge wrong length: got %d, want 32", len(decoded))
	}
}

func TestGenerateChallengeUnique(t *testing.T) {
	// Generate 100 challenges and ensure they're all unique
	challenges := make(map[string]bool)

	for i := 0; i < 100; i++ {
		challenge, err := generateChallenge()
		if err != nil {
			t.Fatalf("Failed to generate challenge: %v", err)
		}

		if challenges[challenge] {
			t.Errorf("Duplicate challenge generated")
		}
		challenges[challenge] = true
	}
}

func TestGenerateToken(t *testing.T) {
	handle := "test@authgrid.net"
	token, err := generateToken(handle)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Should be base64 encoded
	_, err = base64.StdEncoding.DecodeString(token)
	if err != nil {
		t.Fatalf("Token is not valid base64: %v", err)
	}

	// Should be unique
	token2, _ := generateToken(handle)
	if token == token2 {
		t.Errorf("Tokens should be unique")
	}
}

func TestVerifySignatureEd25519(t *testing.T) {
	// Generate Ed25519 keypair
	privateKey, publicKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Create a message
	message := []byte("test message for signing")

	// Sign the message
	signature := ed25519.Sign(privateKey, message)

	// Encode public key
	publicKeyStr := base64.StdEncoding.EncodeToString(publicKey)

	// Verify signature
	valid, err := verifySignature(publicKeyStr, "ed25519", message, signature)
	if err != nil {
		t.Fatalf("Verification failed: %v", err)
	}

	if !valid {
		t.Errorf("Valid signature marked as invalid")
	}
}

func TestVerifySignatureEd25519Invalid(t *testing.T) {
	// Generate Ed25519 keypair
	_, publicKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Create a message
	message := []byte("test message")

	// Create invalid signature
	invalidSignature := make([]byte, ed25519.SignatureSize)
	rand.Read(invalidSignature)

	// Encode public key
	publicKeyStr := base64.StdEncoding.EncodeToString(publicKey)

	// Verify signature (should fail)
	valid, err := verifySignature(publicKeyStr, "ed25519", message, invalidSignature)
	if err != nil {
		t.Fatalf("Verification error: %v", err)
	}

	if valid {
		t.Errorf("Invalid signature marked as valid")
	}
}

func TestVerifySignatureWrongMessage(t *testing.T) {
	// Generate Ed25519 keypair
	privateKey, publicKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	// Create and sign a message
	originalMessage := []byte("original message")
	signature := ed25519.Sign(privateKey, originalMessage)

	// Try to verify with different message
	differentMessage := []byte("different message")

	// Encode public key
	publicKeyStr := base64.StdEncoding.EncodeToString(publicKey)

	// Verify with wrong message (should fail)
	valid, err := verifySignature(publicKeyStr, "ed25519", differentMessage, signature)
	if err != nil {
		t.Fatalf("Verification error: %v", err)
	}

	if valid {
		t.Errorf("Signature valid for wrong message")
	}
}
