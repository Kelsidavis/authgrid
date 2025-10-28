# Authgrid Express.js Example

Simple Express.js application demonstrating Authgrid passwordless authentication with server-side sessions.

## Features

- ✅ User registration
- ✅ Passwordless login
- ✅ Session management
- ✅ Protected routes
- ✅ Logout functionality
- ✅ Clean, modern UI

## Prerequisites

- Node.js 14+ installed
- Authgrid API running (see main README.md)

## Quick Start

### 1. Install Dependencies

```bash
npm install
```

### 2. Start Authgrid API

In the main authgrid directory:
```bash
make run
```

This starts the Authgrid API on http://localhost:8080

### 3. Run the Express Server

```bash
npm start
```

The server will start on http://localhost:3001

### 4. Open in Browser

Visit: http://localhost:3001

## Usage

### Register

1. Click "Register" button
2. A handle will be generated (e.g., `abc123@authgrid.net`)
3. The handle is pre-filled in the login form

### Login

1. Enter your handle (or use the pre-filled one from registration)
2. Click "Login"
3. Your browser signs a challenge with your private key
4. A server-side session is created

### Access Protected Routes

Once logged in:
- Click "Get Profile" to fetch your profile data
- Click "Get Dashboard" to fetch dashboard data
- Both routes require authentication

### Logout

Click "Logout" to end your session.

## Project Structure

```
express-simple/
├── server.js           # Express server with Authgrid integration
├── package.json        # Dependencies
├── public/
│   ├── index.html     # Frontend
│   └── authgrid.js    # Authgrid client SDK
└── README.md          # This file
```

## API Endpoints

### Public Endpoints

**POST /api/register**
- Proxies to Authgrid API
- Registers new user
- Returns handle

**POST /api/challenge**
- Proxies to Authgrid API
- Requests authentication challenge
- Returns challenge and expiry

**POST /api/verify**
- Proxies to Authgrid API
- Verifies signed challenge
- Creates server-side session
- Returns verification result

**GET /api/session**
- Returns current session information
- No authentication required

### Protected Endpoints

**GET /api/profile**
- Requires authentication
- Returns user profile data

**GET /api/dashboard**
- Requires authentication
- Returns dashboard data

**POST /api/logout**
- Requires authentication
- Destroys session

## How It Works

### Registration Flow

```
1. User clicks "Register"
2. Browser generates Ed25519 keypair
3. Public key sent to Express server
4. Express proxies to Authgrid API
5. Authgrid generates handle
6. Handle returned to browser
7. Private key stored in browser localStorage
```

### Login Flow

```
1. User enters handle and clicks "Login"
2. Browser requests challenge from Express
3. Express proxies challenge request to Authgrid
4. Challenge returned to browser
5. Browser signs challenge with private key
6. Signature sent to Express
7. Express proxies verification to Authgrid
8. On success, Express creates session
9. User is authenticated
```

### Accessing Protected Routes

```
1. User clicks "Get Profile"
2. Browser makes request to /api/profile
3. Express checks session middleware
4. If authenticated, returns data
5. If not authenticated, returns 401
```

## Session Configuration

Sessions are configured in `server.js`:

```javascript
app.use(session({
  secret: 'your-secret-key-change-in-production',
  resave: false,
  saveUninitialized: false,
  cookie: {
    secure: false,  // Set to true in production with HTTPS
    maxAge: 24 * 60 * 60 * 1000  // 24 hours
  }
}));
```

**In production:**
- Change the `secret` to a random string
- Set `cookie.secure` to `true`
- Use HTTPS
- Consider using a session store (Redis, MongoDB, etc.)

## Environment Variables

```bash
PORT=3001                          # Server port
AUTHGRID_API=http://localhost:8080 # Authgrid API URL
```

## Customization

### Adding More Protected Routes

```javascript
app.get('/api/my-route', requireAuth, (req, res) => {
  res.json({
    handle: req.session.authgridHandle,
    message: 'This is my protected route'
  });
});
```

### Custom Session Data

```javascript
app.post('/api/verify', async (req, res) => {
  // ... verification code ...

  if (data.verified) {
    req.session.authgridHandle = req.body.handle;
    req.session.authgridToken = data.token;
    req.session.customData = { /* your data */ };
  }
});
```

### Adding Authorization

```javascript
// Role-based access control
function requireRole(role) {
  return (req, res, next) => {
    if (!req.session.authgridHandle) {
      return res.status(401).json({ error: 'Not authenticated' });
    }

    // Check user role from database
    const userRole = getUserRole(req.session.authgridHandle);
    if (userRole !== role) {
      return res.status(403).json({ error: 'Forbidden' });
    }

    next();
  };
}

// Admin-only route
app.get('/api/admin', requireRole('admin'), (req, res) => {
  res.json({ message: 'Admin dashboard' });
});
```

## Production Deployment

### Security Checklist

- [ ] Change session secret
- [ ] Enable HTTPS
- [ ] Set `cookie.secure` to `true`
- [ ] Use session store (not memory)
- [ ] Add rate limiting
- [ ] Add CORS properly
- [ ] Add helmet.js
- [ ] Add input validation
- [ ] Add logging
- [ ] Set proper error handling

### Example Production Config

```javascript
const session = require('express-session');
const RedisStore = require('connect-redis')(session);
const helmet = require('helmet');
const rateLimit = require('express-rate-limit');

// Security headers
app.use(helmet());

// Rate limiting
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100 // limit each IP to 100 requests per windowMs
});
app.use('/api/', limiter);

// Redis session store
app.use(session({
  store: new RedisStore({ client: redisClient }),
  secret: process.env.SESSION_SECRET,
  resave: false,
  saveUninitialized: false,
  cookie: {
    secure: true,
    httpOnly: true,
    maxAge: 24 * 60 * 60 * 1000,
    sameSite: 'strict'
  }
}));
```

## Testing

### Manual Testing

1. Start Authgrid API: `make run`
2. Start Express server: `npm start`
3. Open browser to http://localhost:3001
4. Test registration, login, protected routes, logout

### With cURL

```bash
# Check session (should not be authenticated)
curl http://localhost:3001/api/session

# Try protected route (should fail)
curl http://localhost:3001/api/profile

# After logging in via browser, try again
curl -b cookies.txt http://localhost:3001/api/profile
```

## Troubleshooting

### "Cannot connect to Authgrid API"

Make sure Authgrid API is running:
```bash
curl http://localhost:8080/health
```

If not running:
```bash
cd ../../
make run
```

### "Session not persisting"

- Check that cookies are enabled in browser
- Check browser console for errors
- Try clearing cookies and trying again

### Port already in use

Change the port:
```bash
PORT=3002 npm start
```

## Integration with Your App

To integrate Authgrid into your Express app:

1. **Install dependencies:**
   ```bash
   npm install express-session node-fetch
   ```

2. **Copy these files:**
   - `public/authgrid.js` → Your public directory
   - Relevant parts of `server.js` → Your server

3. **Add the middleware:**
   ```javascript
   app.use(session({ /* config */ }));
   ```

4. **Add the routes:**
   - `/api/register`, `/api/challenge`, `/api/verify`
   - Protected routes with `requireAuth`

5. **Update your frontend:**
   - Include `authgrid.js`
   - Use `authgrid.register()` and `authgrid.authenticate()`

## Further Reading

- [Authgrid Main Documentation](../../README.md)
- [Authgrid Technical Architecture](../../TECHNICAL.md)
- [Authgrid API Reference](../../docs/getting-started.md)

## License

MIT License - See main Authgrid LICENSE file

---

**Questions?** Check the [FAQ](../../FAQ.md) or open an issue on GitHub.
