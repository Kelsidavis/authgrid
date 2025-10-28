# Hybrid Hosting Setup Guide

This guide shows you how to host your Authgrid frontend at **https://authgrid.org/** while running the backend **locally** on your own machine.

## Architecture Overview

```
┌──────────────────────────────────────┐
│   https://authgrid.org/              │
│   (Static files on web host)         │
│   - index.html                       │
│   - authgrid.js                      │
│   - config.js ← Configure API URL    │
└──────────────┬───────────────────────┘
               │ HTTPS requests
               ▼
┌──────────────────────────────────────┐
│   Public Tunnel                      │
│   (ngrok/Cloudflare/localtunnel)    │
│   https://abc123.ngrok-free.app      │
└──────────────┬───────────────────────┘
               │ Tunnels to
               ▼
┌──────────────────────────────────────┐
│   Your Local Machine                 │
│   Backend API: localhost:8080        │
│   PostgreSQL: localhost:5432         │
└──────────────────────────────────────┘
```

## Step 1: Configure Backend CORS

✅ **Already done!** The API in `src/api/main.go` is now configured to accept requests from:
- `https://authgrid.org`
- `https://www.authgrid.org`
- `http://localhost:3000` (local development)

## Step 2: Start Your Local Backend

Start the backend API and database locally:

```bash
# Start API + PostgreSQL with Docker Compose
make run

# Or manually:
docker-compose up -d
```

Verify it's running:
```bash
curl http://localhost:8080/health
# Should return: {"status":"healthy","time":"..."}
```

## Step 3: Expose Your Local API to the Internet

Choose one of these methods to make your local API accessible from the internet:

### Option A: Ngrok (Recommended - Easy & Reliable)

**Install:**
```bash
# macOS
brew install ngrok

# Linux
curl -s https://ngrok-agent.s3.amazonaws.com/ngrok.asc | \
  sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null && \
  echo "deb https://ngrok-agent.s3.amazonaws.com buster main" | \
  sudo tee /etc/apt/sources.list.d/ngrok.list && \
  sudo apt update && sudo apt install ngrok
```

**Run:**
```bash
# Sign up at ngrok.com and authenticate
ngrok config add-authtoken YOUR_TOKEN

# Start tunnel
ngrok http 8080
```

You'll get a URL like: `https://abc123.ngrok-free.app`

**Note:** Free ngrok URLs change each time you restart. Paid plans ($8/mo) give you a static domain.

### Option B: Cloudflare Tunnel (Free & Fast)

**Install:**
```bash
# macOS
brew install cloudflare/cloudflare/cloudflared

# Linux
wget -q https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared-linux-amd64.deb
```

**Run:**
```bash
cloudflared tunnel --url http://localhost:8080
```

You'll get a URL like: `https://your-tunnel.trycloudflare.com`

**Note:** Free, no signup required. URL changes each restart.

### Option C: Localtunnel (Quick & Simple)

**Install:**
```bash
npm install -g localtunnel
```

**Run:**
```bash
lt --port 8080 --subdomain authgrid-api
```

You'll get: `https://authgrid-api.loca.lt`

**Note:** May require password entry on first visit.

### Option D: Your Own Domain with Reverse SSH Tunnel (Advanced)

If you have your own server:

```bash
# On your local machine
ssh -R 8080:localhost:8080 user@your-server.com

# On your server, configure nginx to proxy to localhost:8080
```

## Step 4: Configure Frontend

Edit `examples/demo/config.js` and set your tunnel URL:

```javascript
// Update this line with your tunnel URL
window.AUTHGRID_API_URL = 'https://abc123.ngrok-free.app';
```

**Examples:**
```javascript
// Ngrok
window.AUTHGRID_API_URL = 'https://abc123.ngrok-free.app';

// Cloudflare Tunnel
window.AUTHGRID_API_URL = 'https://authgrid-tunnel.trycloudflare.com';

// Localtunnel
window.AUTHGRID_API_URL = 'https://authgrid-api.loca.lt';
```

## Step 5: Deploy Frontend to Your Web Host

Upload these files to your web host (https://authgrid.org/):

```
examples/demo/
├── index.html      ← Main page
├── authgrid.js     ← Client SDK (from src/client/authgrid.js)
└── config.js       ← API configuration (with your tunnel URL)
```

### Deployment Options:

#### Option 1: Traditional Web Host (cPanel, FTP)
1. Copy `src/client/authgrid.js` to `examples/demo/`
2. Upload `examples/demo/*` to your web host's public_html folder
3. Access https://authgrid.org/

#### Option 2: Static Hosting (Netlify, Vercel, GitHub Pages)

**Netlify/Vercel:**
```bash
cd examples/demo
# Deploy via drag-and-drop or CLI
netlify deploy --prod --dir=.
```

**GitHub Pages:**
```bash
# Push examples/demo/ to gh-pages branch
git subtree push --prefix examples/demo origin gh-pages
```

#### Option 3: CloudFlare Pages
1. Connect your GitHub repo
2. Set build directory: `examples/demo`
3. Deploy

## Step 6: Test Your Setup

1. **Visit https://authgrid.org/**
2. **Click "Register Now"**
   - Should generate keypair and register with your local backend
3. **Check backend logs:**
   ```bash
   make logs
   # Should see: POST /register 200 OK
   ```
4. **Click "Login"**
   - Should authenticate using the registered keypair

## Troubleshooting

### CORS Errors
**Error:** `Access-Control-Allow-Origin` error in browser console

**Fix:** Make sure your tunnel URL is added to CORS in `src/api/main.go:57-61`

```go
allowedOrigins := []string{
    "https://authgrid.org",
    "https://www.authgrid.org",
    "http://localhost:3000",
    // Add your tunnel URL if testing with custom domain
}
```

Rebuild after changes:
```bash
docker-compose down
docker-compose up --build -d
```

### Connection Refused
**Error:** Frontend can't reach API

**Check:**
1. Backend is running: `curl http://localhost:8080/health`
2. Tunnel is active: `curl https://your-tunnel.url/health`
3. `config.js` has correct tunnel URL

### Mixed Content Warning (HTTP/HTTPS)
**Error:** "Mixed content blocked" when loading HTTP resources from HTTPS site

**Fix:** Always use HTTPS tunnels (ngrok, Cloudflare provide HTTPS by default)

### Ngrok Free Plan Interruptions
**Error:** "Visit Site" button appears on ngrok URLs

**Workaround:**
- Users need to click "Visit Site" once
- Upgrade to ngrok paid plan ($8/mo) for no interruptions
- Use Cloudflare Tunnel (free, no interruptions)

## Production Considerations

### Security
- **HTTPS Required:** Always use HTTPS tunnels in production
- **Rate Limiting:** Already configured (100 req/sec)
- **Private Keys:** Never leave user's device (stored in browser localStorage)

### Reliability
- **Tunnel Uptime:** Keep tunnel running (use systemd service or screen/tmux)
- **Static Domains:** Use paid ngrok or your own domain for consistent URLs
- **Database Backups:** Backup PostgreSQL regularly
  ```bash
  docker exec authgrid-postgres-1 pg_dump -U authgrid authgrid > backup.sql
  ```

### Monitoring
```bash
# Watch API logs in real-time
make logs

# Check health endpoint
curl https://your-tunnel.url/health

# Monitor tunnel status
ngrok http 8080 --log=stdout
```

### Keeping Tunnel Running
Use a process manager to keep the tunnel running:

**systemd service (Linux):**
```bash
# /etc/systemd/system/authgrid-tunnel.service
[Unit]
Description=Authgrid Ngrok Tunnel
After=network.target

[Service]
Type=simple
User=youruser
ExecStart=/usr/local/bin/ngrok http 8080
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable:
```bash
sudo systemctl enable authgrid-tunnel
sudo systemctl start authgrid-tunnel
```

**tmux/screen (Quick):**
```bash
# Start detached session
tmux new-session -d -s authgrid 'ngrok http 8080'

# Reattach later
tmux attach -t authgrid
```

## Alternative: Fully Local Setup

If you prefer everything local (no tunnel), you can:

1. Host frontend locally too:
   ```bash
   cd examples/demo
   python3 -m http.server 3000
   # Access at http://localhost:3000
   ```

2. Or use the included Docker setup:
   ```bash
   make run  # Frontend on :3000, API on :8080
   ```

## Summary

✅ **What's hosted remotely:** Static frontend files (HTML, JS, CSS)
✅ **What's local:** Backend API, PostgreSQL database, user data
✅ **Connection:** Secure HTTPS tunnel from cloud to your machine
✅ **Cost:** Free (or $8/mo for static ngrok domain)
✅ **Control:** Full control over backend and data

---

**Questions?** Check the [main README](README.md) or open an issue on GitHub.
