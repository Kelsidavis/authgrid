# Authgrid API Server

Core API server for Authgrid passwordless authentication.

## Features

- ✅ Ed25519 public key cryptography
- ✅ Challenge-response authentication
- ✅ PostgreSQL storage
- ✅ Rate limiting
- ✅ CORS support
- ✅ RESTful JSON API

## API Endpoints

### POST /register

Register a new user and get a handle.

**Request:**
```json
{
  "public_key": "base64_encoded_ed25519_public_key",
  "key_type": "ed25519"
}
```

**Response: 201 Created**
```json
{
  "handle": "abc123def4@authgrid.net",
  "id": "uuid",
  "created_at": "2025-01-15T10:30:00Z"
}
```

---

### POST /challenge

Request an authentication challenge for a handle.

**Request:**
```json
{
  "handle": "abc123def4@authgrid.net"
}
```

**Response: 200 OK**
```json
{
  "challenge": "base64_encoded_random_bytes",
  "expires_at": "2025-01-15T10:35:00Z"
}
```

---

### POST /verify

Verify a signed challenge to complete authentication.

**Request:**
```json
{
  "handle": "abc123def4@authgrid.net",
  "challenge": "base64_encoded_challenge",
  "signature": "base64_encoded_signature"
}
```

**Response: 200 OK**
```json
{
  "verified": true,
  "token": "authentication_token",
  "expires_at": "2025-01-16T10:30:00Z"
}
```

---

### GET /user/:handle

Get public information about a user (optional endpoint).

**Response: 200 OK**
```json
{
  "handle": "abc123def4@authgrid.net",
  "public_key": "base64_encoded_public_key",
  "created_at": "2025-01-15T10:30:00Z"
}
```

---

### GET /health

Health check endpoint.

**Response: 200 OK**
```json
{
  "status": "healthy",
  "time": "2025-01-15T10:30:00Z"
}
```

## Development

### Prerequisites

- Go 1.21+
- PostgreSQL 15+

### Setup

1. Install dependencies:
```bash
go mod download
```

2. Set up environment variables:
```bash
export DATABASE_URL="postgres://authgrid:authgrid@localhost:5432/authgrid?sslmode=disable"
export PORT="8080"
export AUTHGRID_DOMAIN="authgrid.net"
```

3. Run the server:
```bash
go run .
```

### Testing

Run tests:
```bash
go test -v ./...
```

## Deployment

### Docker

Build and run with Docker:
```bash
docker build -t authgrid-api .
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  authgrid-api
```

### Production

For production deployment:

1. Use proper TLS/HTTPS
2. Set secure DATABASE_URL with SSL
3. Configure CORS allowed origins
4. Set up monitoring and logging
5. Use proper JWT implementation for tokens
6. Enable database connection pooling
7. Add request ID tracking
8. Implement proper error logging

## Configuration

Environment variables:

- `DATABASE_URL` - PostgreSQL connection string
- `PORT` - Server port (default: 8080)
- `AUTHGRID_DOMAIN` - Domain for handle generation (default: authgrid.net)

## Security

- Private keys never reach the server
- Challenges expire after 5 minutes
- Challenges are single-use only
- Rate limiting: 100 req/sec globally
- HTTPS required in production
- Database credentials should be rotated regularly

## License

MIT License - See LICENSE file for details
