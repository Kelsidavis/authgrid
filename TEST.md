# Testing Authgrid

Quick guide to test the Authgrid MVP.

## Prerequisites

Authgrid should be running:
```bash
docker compose ps
# All 3 containers should be "Up" or "Healthy"
```

If not running:
```bash
make run
```

---

## Test 1: API Health Check

```bash
curl http://localhost:8080/health
```

**Expected:**
```json
{
  "status": "healthy",
  "time": "2025-10-28T..."
}
```

✅ **Pass** if you see the healthy status.

---

## Test 2: Demo Web Interface

1. **Open the demo:**
   ```
   http://localhost:3000
   ```

2. **Register a new user:**
   - Click "Register Now" button
   - Wait ~1 second
   - You should see: "✓ Registration successful!"
   - Your handle will be displayed (e.g., `f8a9b3c2d1@authgrid.net`)

3. **Login with your handle:**
   - Your handle should be pre-filled in the login form
   - Click "Login" button
   - You should see: "✓ Login successful!"
   - A token will be displayed

✅ **Pass** if both registration and login work.

---

## Test 3: Multiple Handles

1. Open the demo in a **private/incognito window**
2. Click "Register Now" again
3. You'll get a different handle
4. Login with the new handle

✅ **Pass** if you can have multiple handles and login with each.

---

## Test 4: Browser Developer Console

Open browser DevTools (F12) and check:

1. **Console tab** — Should have no errors
2. **Network tab** — Check the API calls:
   - `POST /register` → 201 Created
   - `POST /challenge` → 200 OK
   - `POST /verify` → 200 OK

✅ **Pass** if all API calls return success.

---

## Test 5: Database Verification

Check if users are stored in the database:

```bash
docker compose exec postgres psql -U authgrid -d authgrid -c "SELECT handle, key_type, created_at FROM users;"
```

**Expected:** List of registered users with their handles.

✅ **Pass** if you see your registered handles.

---

## Test 6: Key Storage

Check browser localStorage:

1. Open DevTools (F12) → Application tab → Local Storage
2. Look for keys starting with `authgrid_keypair_`
3. You should see one entry per registered handle

✅ **Pass** if your keypairs are stored locally.

---

## Test 7: Cross-Browser Testing

Test in different browsers:

- Chrome/Edge
- Firefox
- Safari (Mac)

All should work, though they may use different algorithms:
- Modern browsers: Ed25519
- Safari/older browsers: ECDSA P-256

✅ **Pass** if registration and login work in all browsers.

---

## Test 8: API Direct Testing

Test the API endpoints directly:

### Test Registration (requires Go or Python)

This is complex because you need to generate a proper keypair. Use the web demo instead, or:

```bash
# Register via web demo, then test challenge/verify manually
HANDLE="your-handle@authgrid.net"  # Use your actual handle from demo

# Request challenge
curl -X POST http://localhost:8080/challenge \
  -H "Content-Type: application/json" \
  -d "{\"handle\":\"$HANDLE\"}"

# You'll get a challenge back
# To verify, you need to sign it with your private key (use the web demo)
```

✅ **Pass** if challenge is returned.

---

## Test 9: Error Handling

### Test with invalid handle:
```bash
curl -X POST http://localhost:8080/challenge \
  -H "Content-Type: application/json" \
  -d '{"handle":"nonexistent@authgrid.net"}'
```

**Expected:**
```json
{
  "error": "Handle not found"
}
```

✅ **Pass** if you get a proper error message.

---

## Test 10: Rate Limiting

Send 200 requests rapidly:
```bash
for i in {1..200}; do
  curl -s http://localhost:8080/health > /dev/null &
done
wait
```

Some requests may fail with 429 (rate limit).

✅ **Pass** if rate limiting is working.

---

## Common Issues

### "Registration failed: Invalid public key length"
- This was fixed in v0.1.1
- Make sure you rebuilt: `docker compose build api && docker compose up -d`

### "Failed to connect to database"
- Check: `docker compose logs postgres`
- Restart: `docker compose restart postgres`

### "Port already in use"
- Stop Authgrid: `make stop`
- Find what's using the port: `lsof -i :8080`
- Change port in `docker-compose.yml` if needed

### Demo page won't load
- Check: `docker compose logs demo`
- Verify: `curl http://localhost:3000`
- Restart: `docker compose restart demo`

---

## Quick Troubleshooting

```bash
# View all logs
make logs

# View API logs only
make logs-api

# View database logs
make logs-db

# Restart everything
make stop && make run

# Reset database (WARNING: deletes all data)
make clean && make run
```

---

## Success Criteria

All tests should pass:
- ✅ API health check works
- ✅ Demo registration works
- ✅ Demo login works
- ✅ Multiple handles work
- ✅ No console errors
- ✅ Data stored in database
- ✅ Keys stored in browser
- ✅ Works in multiple browsers
- ✅ API endpoints respond correctly
- ✅ Error handling works

---

## Next Steps After Testing

If all tests pass:

1. **Read the documentation:**
   - [QUICKSTART.md](QUICKSTART.md) — Integration guide
   - [TECHNICAL.md](TECHNICAL.md) — Architecture
   - [docs/getting-started.md](docs/getting-started.md) — Developer guide

2. **Try integrating:**
   - Copy `src/client/authgrid.js` to your project
   - Add it to your HTML
   - Use the SDK to add passwordless auth

3. **Give feedback:**
   - Report bugs via GitHub issues
   - Suggest features
   - Contribute code (see [CONTRIBUTING.md](CONTRIBUTING.md))

---

**Happy Testing!** 🚀

If you find any issues, see [BUGFIX_ECDSA.md](BUGFIX_ECDSA.md) for the most recent fixes.
