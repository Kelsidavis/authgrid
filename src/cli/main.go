package main

import (
	"bufio"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	version = "0.1.1-alpha"
)

var (
	apiURL      string
	keystoreDir string
)

func main() {
	// Global flags
	flag.StringVar(&apiURL, "api", "http://localhost:8080", "Authgrid API URL")
	flag.StringVar(&keystoreDir, "keystore", getDefaultKeystoreDir(), "Keystore directory")

	// Subcommands
	registerCmd := flag.NewFlagSet("register", flag.ExitOnError)
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)

	// Login flags
	loginHandle := loginCmd.String("handle", "", "Handle to authenticate with")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Parse global flags
	flag.Parse()

	// Get subcommand
	subcommand := os.Args[1]

	switch subcommand {
	case "register":
		registerCmd.Parse(os.Args[2:])
		handleRegister()

	case "login":
		loginCmd.Parse(os.Args[2:])
		if *loginHandle == "" {
			fmt.Println("Error: --handle flag is required")
			loginCmd.PrintDefaults()
			os.Exit(1)
		}
		handleLogin(*loginHandle)

	case "list":
		listCmd.Parse(os.Args[2:])
		handleList()

	case "version":
		versionCmd.Parse(os.Args[2:])
		fmt.Printf("authgrid-cli version %s\n", version)

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Authgrid CLI - Passwordless authentication tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  authgrid [command] [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  register          Register a new user and get a handle")
	fmt.Println("  login             Authenticate with a handle")
	fmt.Println("  list              List stored handles")
	fmt.Println("  version           Show version information")
	fmt.Println("  help              Show this help message")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --api URL         Authgrid API URL (default: http://localhost:8080)")
	fmt.Println("  --keystore DIR    Keystore directory (default: ~/.authgrid)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  authgrid register")
	fmt.Println("  authgrid login --handle abc123@authgrid.net")
	fmt.Println("  authgrid list")
	fmt.Println()
}

func handleRegister() {
	fmt.Println("Registering new user...")

	// Generate Ed25519 keypair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("Error generating keypair: %v\n", err)
		os.Exit(1)
	}

	// Encode public key
	publicKeyB64 := base64.StdEncoding.EncodeToString(publicKey)

	// Register with API
	reqBody := map[string]string{
		"public_key": publicKeyB64,
		"key_type":   "ed25519",
	}

	resp, err := makeRequest("POST", apiURL+"/register", reqBody)
	if err != nil {
		fmt.Printf("Error registering: %v\n", err)
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		os.Exit(1)
	}

	handle, ok := result["handle"].(string)
	if !ok {
		fmt.Println("Error: invalid response from server")
		os.Exit(1)
	}

	// Save keypair to keystore
	if err := saveKeypair(handle, privateKey, publicKey); err != nil {
		fmt.Printf("Error saving keypair: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("✅ Registration successful!")
	fmt.Printf("   Handle: %s\n", handle)
	fmt.Printf("   Keystore: %s\n", keystoreDir)
	fmt.Println()
	fmt.Println("To login:")
	fmt.Printf("   authgrid login --handle %s\n", handle)
	fmt.Println()
}

func handleLogin(handle string) {
	fmt.Printf("Logging in as %s...\n", handle)

	// Load keypair
	privateKey, publicKey, err := loadKeypair(handle)
	if err != nil {
		fmt.Printf("Error loading keypair: %v\n", err)
		fmt.Println("Have you registered this handle? Try: authgrid register")
		os.Exit(1)
	}

	// Request challenge
	reqBody := map[string]string{
		"handle": handle,
	}

	resp, err := makeRequest("POST", apiURL+"/challenge", reqBody)
	if err != nil {
		fmt.Printf("Error requesting challenge: %v\n", err)
		os.Exit(1)
	}

	var challengeResp map[string]interface{}
	if err := json.Unmarshal(resp, &challengeResp); err != nil {
		fmt.Printf("Error parsing challenge: %v\n", err)
		os.Exit(1)
	}

	challengeB64, ok := challengeResp["challenge"].(string)
	if !ok {
		fmt.Println("Error: invalid challenge response")
		os.Exit(1)
	}

	// Decode challenge
	challenge, err := base64.StdEncoding.DecodeString(challengeB64)
	if err != nil {
		fmt.Printf("Error decoding challenge: %v\n", err)
		os.Exit(1)
	}

	// Sign challenge
	signature := ed25519.Sign(privateKey, challenge)
	signatureB64 := base64.StdEncoding.EncodeToString(signature)

	// Verify signature
	verifyBody := map[string]string{
		"handle":    handle,
		"challenge": challengeB64,
		"signature": signatureB64,
	}

	resp, err = makeRequest("POST", apiURL+"/verify", verifyBody)
	if err != nil {
		fmt.Printf("Error verifying signature: %v\n", err)
		os.Exit(1)
	}

	var verifyResp map[string]interface{}
	if err := json.Unmarshal(resp, &verifyResp); err != nil {
		fmt.Printf("Error parsing verify response: %v\n", err)
		os.Exit(1)
	}

	verified, ok := verifyResp["verified"].(bool)
	if !ok || !verified {
		fmt.Println("❌ Authentication failed")
		os.Exit(1)
	}

	token, _ := verifyResp["token"].(string)

	fmt.Println()
	fmt.Println("✅ Login successful!")
	fmt.Printf("   Handle: %s\n", handle)
	fmt.Printf("   Token: %s...\n", token[:40])
	fmt.Println()

	// For demonstration, verify the public key
	_ = publicKey // Use the variable to avoid unused warning
}

func handleList() {
	// Ensure keystore exists
	if _, err := os.Stat(keystoreDir); os.IsNotExist(err) {
		fmt.Println("No handles registered yet.")
		fmt.Println("Try: authgrid register")
		return
	}

	// List files in keystore
	files, err := os.ReadDir(keystoreDir)
	if err != nil {
		fmt.Printf("Error reading keystore: %v\n", err)
		os.Exit(1)
	}

	var handles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".key") {
			handle := strings.TrimSuffix(file.Name(), ".key")
			handles = append(handles, handle)
		}
	}

	if len(handles) == 0 {
		fmt.Println("No handles registered yet.")
		fmt.Println("Try: authgrid register")
		return
	}

	fmt.Printf("Stored handles (%d):\n", len(handles))
	for _, handle := range handles {
		fmt.Printf("  • %s\n", handle)
	}
	fmt.Println()
	fmt.Println("To login:")
	fmt.Printf("  authgrid login --handle <handle>\n")
	fmt.Println()
}

// Helper functions

func makeRequest(method, url string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = strings.NewReader(string(jsonData))
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errResp map[string]string
		json.Unmarshal(respBody, &errResp)
		if msg, ok := errResp["error"]; ok {
			return nil, fmt.Errorf("API error: %s", msg)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func getDefaultKeystoreDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".authgrid"
	}
	return filepath.Join(home, ".authgrid")
}

func saveKeypair(handle string, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) error {
	// Create keystore directory
	if err := os.MkdirAll(keystoreDir, 0700); err != nil {
		return err
	}

	// Save keypair
	filename := filepath.Join(keystoreDir, handle+".key")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(base64.StdEncoding.EncodeToString(privateKey) + "\n")
	writer.WriteString(base64.StdEncoding.EncodeToString(publicKey) + "\n")
	return writer.Flush()
}

func loadKeypair(handle string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	filename := filepath.Join(keystoreDir, handle+".key")
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read private key
	if !scanner.Scan() {
		return nil, nil, fmt.Errorf("invalid keyfile")
	}
	privateKeyB64 := scanner.Text()
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		return nil, nil, err
	}

	// Read public key
	if !scanner.Scan() {
		return nil, nil, fmt.Errorf("invalid keyfile")
	}
	publicKeyB64 := scanner.Text()
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return nil, nil, err
	}

	return ed25519.PrivateKey(privateKey), ed25519.PublicKey(publicKey), nil
}
