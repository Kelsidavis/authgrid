# Authgrid Quickstart

Get Authgrid running on your machine in 5 minutes.

## Prerequisites

- Docker and Docker Compose
- Make (optional, for convenience commands)

## Option 1: Quick Start with Docker (Recommended)

### 1. Clone and start services

```bash
# Navigate to authgrid directory
cd authgrid

# Start all services (API, Database, Demo)
make run
```

Or without Make:

```bash
docker-compose up -d
```

### 2. Access the services

- **Demo UI:** http://localhost:3000
- **API Server:** http://localhost:8080
- **Health Check:** http://localhost:8080/health

### 3. Try it out!

1. Open http://localhost:3000 in your browser
2. Click **"Register Now"** to create a new user
3. You'll receive a handle like `abc123def4@authgrid.net`
4. Click **"Login"** to authenticate using your handle

That's it! You're running Authgrid.

---

## Option 2: Manual Setup (Development)

### 1. Start PostgreSQL

```bash
docker-compose up -d postgres
```

Wait for PostgreSQL to be ready:

```bash
docker-compose logs -f postgres
# Wait for "database system is ready to accept connections"
```

### 2. Run the API server

```bash
cd src/api

# Install dependencies
go mod download

# Set environment variables
export DATABASE_URL="postgres://authgrid:authgrid@localhost:5432/authgrid?sslmode=disable"
export PORT="8080"

# Run the server
go run .
```

### 3. Serve the demo frontend

```bash
# In a new terminal
cd examples/demo
python3 -m http.server 3000
```

Or use any static file server.

### 4. Access the demo

Open http://localhost:3000 in your browser.

---

## Useful Commands

### Using Make

```bash
make help          # Show all available commands
make build         # Build Docker containers
make run           # Start all services
make stop          # Stop all services
make logs          # View all logs
make logs-api      # View API logs only
make db-shell      # Open PostgreSQL shell
make clean         # Stop and remove all data
```

### Using Docker Compose directly

```bash
docker-compose up -d        # Start services in background
docker-compose down         # Stop services
docker-compose logs -f      # Follow logs
docker-compose ps           # Show running services
docker-compose restart api  # Restart API server
```

---

## Testing the API

### Using curl

**Register a user:**

First, generate a keypair (you'll need a proper Ed25519 library for this, or use the JavaScript client). For testing, you can use this example public key:

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "public_key": "MCowBQYDK2VwAyEAGb9ECWmEzf6FQbrBZ9w7lshQhqoWtWLLxPEo1Qk1234=",
    "key_type": "ed25519"
  }'
```

**Request a challenge:**

```bash
curl -X POST http://localhost:8080/challenge \
  -H "Content-Type: application/json" \
  -d '{
    "handle": "abc123def4@authgrid.net"
  }'
```

**Health check:**

```bash
curl http://localhost:8080/health
```

---

## Troubleshooting

### Port already in use

If port 8080 or 3000 is already in use, edit `docker-compose.yml` to change the port mappings:

```yaml
ports:
  - "8081:8080"  # Change 8081 to any available port
```

### Database connection failed

Check if PostgreSQL is running:

```bash
docker-compose ps
```

View database logs:

```bash
make logs-db
```

Reset the database (WARNING: deletes all data):

```bash
make db-reset
```

### CORS errors in browser

Make sure the API is running on the expected URL. The demo is configured for `http://localhost:8080`.

If you change the API port, update it in `examples/demo/index.html`:

```javascript
const authgrid = new AuthgridClient({
    apiUrl: 'http://localhost:YOUR_PORT'
});
```

### WebCrypto not available

WebCrypto API requires HTTPS or localhost. Make sure you're accessing the demo via:
- `http://localhost:3000` (✓ works)
- Not `http://127.0.0.1:3000` (might not work in some browsers)
- Not `http://your-ip:3000` (won't work without HTTPS)

---

## What's Next?

### Explore the codebase

- **API Server:** `src/api/` — Go backend with Ed25519 verification
- **Client SDK:** `src/client/` — JavaScript library for browsers
- **Demo:** `examples/demo/` — Working example application
- **Database:** `migrations/` — PostgreSQL schema

### Read the documentation

- [VISION.md](VISION.md) — Core concept and architecture
- [TECHNICAL.md](TECHNICAL.md) — Technical deep-dive
- [docs/getting-started.md](docs/getting-started.md) — Integration guide

### Integrate into your app

1. Copy `src/client/authgrid.js` into your project
2. Include it in your HTML:
   ```html
   <script src="authgrid.js"></script>
   ```
3. Use the client:
   ```javascript
   const authgrid = new AuthgridClient({
       apiUrl: 'http://localhost:8080'
   });

   // Register
   const { handle } = await authgrid.register();

   // Login
   const { token } = await authgrid.authenticate(handle);
   ```

See [docs/getting-started.md](docs/getting-started.md) for complete integration examples.

---

## Production Deployment

This quickstart is for **development only**. For production:

1. Use proper HTTPS/TLS
2. Configure production database with SSL
3. Set secure CORS origins
4. Use proper JWT tokens (not random strings)
5. Add monitoring and logging
6. Set up backups
7. Configure rate limiting per-user
8. Review security checklist in [TECHNICAL.md](TECHNICAL.md)

---

## Getting Help

- **Issues:** [GitHub Issues](https://github.com/Kelsidavis/authgrid/issues)
- **Documentation:** See `/docs` folder
- **Security:** security@authgrid.net (coming soon)

---

**Welcome to Authgrid!** 🎉

Passwords die here. Let's build the future of authentication.
