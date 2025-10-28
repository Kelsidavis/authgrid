# Getting Started with Authgrid

This guide will help you add Authgrid passwordless authentication to your application in under 5 minutes.

## Overview

Authgrid replaces traditional email/password authentication with cryptographic keys. Users get a handle like `n5xtc4q3f3@authgrid.net` that works across any service using Authgrid.

**What you'll build:**
- User registration (generates handle + keypair)
- User login (challenge-response flow)
- Protected routes (session management)

---

## Prerequisites

- Node.js 18+ (or Python 3.9+, Go 1.21+)
- A modern browser with WebAuthn support
- 5 minutes of your time

---

## Quick Start (JavaScript/Node.js)

### 1. Install the SDK

```bash
npm install authgrid-client
```

### 2. Initialize Authgrid

```javascript
// server.js
import express from 'express';
import { Authgrid } from 'authgrid-client';

const app = express();
app.use(express.json());

const authgrid = new Authgrid({
  apiUrl: 'https://api.authgrid.net',
  domain: 'authgrid.net'
  // For custom domain (Pro tier): domain: 'yourapp.com'
});
```

### 3. Add Registration Endpoint

```javascript
// POST /register - Register a new user
app.post('/register', async (req, res) => {
  try {
    const result = await authgrid.register();

    // Store the handle in your database
    // await db.users.create({ handle: result.handle });

    res.json({
      success: true,
      handle: result.handle,
      message: 'Registration successful! Save this handle.'
    });
  } catch (error) {
    res.status(400).json({ error: error.message });
  }
});
```

### 4. Add Login Endpoint

```javascript
// POST /login - Authenticate a user
app.post('/login', async (req, res) => {
  const { handle } = req.body;

  try {
    const token = await authgrid.authenticate(handle);

    // Create session in your database
    // await db.sessions.create({ handle, token });

    res.json({
      success: true,
      token,
      message: 'Login successful!'
    });
  } catch (error) {
    res.status(401).json({ error: 'Authentication failed' });
  }
});
```

### 5. Add Authentication Middleware

```javascript
// Middleware to protect routes
async function requireAuth(req, res, next) {
  const token = req.headers.authorization?.replace('Bearer ', '');

  if (!token) {
    return res.status(401).json({ error: 'No token provided' });
  }

  try {
    const verified = await authgrid.verifyToken(token);
    req.user = verified;
    next();
  } catch (error) {
    res.status(401).json({ error: 'Invalid token' });
  }
}

// Protected route example
app.get('/profile', requireAuth, (req, res) => {
  res.json({
    handle: req.user.handle,
    message: 'This is a protected route!'
  });
});
```

### 6. Add Client-Side Code

```html
<!-- register.html -->
<!DOCTYPE html>
<html>
<head>
  <title>Register - Authgrid Demo</title>
</head>
<body>
  <h1>Register</h1>
  <button id="registerBtn">Register with Authgrid</button>
  <div id="result"></div>

  <script type="module">
    import { AuthgridClient } from 'authgrid-client/browser';

    const client = new AuthgridClient({
      apiUrl: 'https://api.authgrid.net'
    });

    document.getElementById('registerBtn').addEventListener('click', async () => {
      try {
        const result = await client.register();

        document.getElementById('result').innerHTML = `
          <p>Success! Your handle: <strong>${result.handle}</strong></p>
          <p>Save this handle to log in later.</p>
        `;

        // Store handle in localStorage
        localStorage.setItem('authgrid_handle', result.handle);
      } catch (error) {
        document.getElementById('result').innerHTML = `
          <p style="color: red;">Error: ${error.message}</p>
        `;
      }
    });
  </script>
</body>
</html>
```

```html
<!-- login.html -->
<!DOCTYPE html>
<html>
<head>
  <title>Login - Authgrid Demo</title>
</head>
<body>
  <h1>Login</h1>
  <input type="text" id="handleInput" placeholder="your-handle@authgrid.net" />
  <button id="loginBtn">Login</button>
  <div id="result"></div>

  <script type="module">
    import { AuthgridClient } from 'authgrid-client/browser';

    const client = new AuthgridClient({
      apiUrl: 'https://api.authgrid.net'
    });

    // Pre-fill handle if saved
    const savedHandle = localStorage.getItem('authgrid_handle');
    if (savedHandle) {
      document.getElementById('handleInput').value = savedHandle;
    }

    document.getElementById('loginBtn').addEventListener('click', async () => {
      const handle = document.getElementById('handleInput').value;

      try {
        const token = await client.authenticate(handle);

        document.getElementById('result').innerHTML = `
          <p style="color: green;">Login successful!</p>
        `;

        // Store token for authenticated requests
        localStorage.setItem('authgrid_token', token);

        // Redirect to dashboard
        setTimeout(() => {
          window.location.href = '/dashboard';
        }, 1000);
      } catch (error) {
        document.getElementById('result').innerHTML = `
          <p style="color: red;">Login failed: ${error.message}</p>
        `;
      }
    });
  </script>
</body>
</html>
```

### 7. Run Your Server

```bash
node server.js
```

Visit `http://localhost:3000/register.html` to test registration!

---

## What Just Happened?

1. **User clicks "Register"**
   - Browser generates a keypair using WebAuthn
   - Private key stored in browser's secure storage
   - Public key sent to Authgrid API
   - User receives a handle: `n5xtc4q3f3@authgrid.net`

2. **User logs in later**
   - User enters their handle
   - Authgrid API sends a challenge (random bytes)
   - Browser signs the challenge with private key
   - Authgrid verifies the signature
   - User receives an auth token

3. **User makes authenticated requests**
   - Include token in `Authorization: Bearer <token>` header
   - Your middleware verifies the token
   - User accesses protected routes

---

## Next Steps

### Customize the Experience

```javascript
const authgrid = new Authgrid({
  apiUrl: 'https://api.authgrid.net',
  domain: 'authgrid.net',
  options: {
    // Require biometric authentication
    userVerification: 'required',

    // Custom timeout for WebAuthn
    timeout: 60000,

    // Custom challenge expiry
    challengeExpiry: 300 // 5 minutes
  }
});
```

### Store User Data

```javascript
app.post('/register', async (req, res) => {
  const { displayName, email } = req.body; // Optional metadata

  const result = await authgrid.register();

  // Store in your database
  await db.users.create({
    handle: result.handle,
    displayName,
    email, // Optional - can map handle to email for recovery
    createdAt: new Date()
  });

  res.json({ handle: result.handle });
});
```

### Add Session Management

```javascript
import session from 'express-session';

app.use(session({
  secret: 'your-secret-key',
  resave: false,
  saveUninitialized: false,
  cookie: { secure: true, httpOnly: true }
}));

app.post('/login', async (req, res) => {
  const { handle } = req.body;
  const token = await authgrid.authenticate(handle);

  // Store in session
  req.session.handle = handle;
  req.session.token = token;

  res.json({ success: true });
});

function requireAuth(req, res, next) {
  if (!req.session.token) {
    return res.status(401).json({ error: 'Not authenticated' });
  }
  next();
}
```

### Handle Recovery

```javascript
// Enable social recovery during registration
app.post('/register', async (req, res) => {
  const { recoveryContacts } = req.body; // Array of email/handle

  const result = await authgrid.register({
    recovery: {
      type: 'social',
      contacts: recoveryContacts,
      threshold: 2 // Need 2 out of N contacts to recover
    }
  });

  res.json({ handle: result.handle });
});
```

---

## Framework-Specific Guides

### Express.js (Detailed)

See [docs/integrations/express.md](../integrations/express.md)

### Next.js

See [docs/integrations/nextjs.md](../integrations/nextjs.md)

### Django (Python)

See [docs/integrations/django.md](../integrations/django.md)

### Flask (Python)

See [docs/integrations/flask.md](../integrations/flask.md)

### Go

See [docs/integrations/go.md](../integrations/go.md)

---

## Production Checklist

Before deploying to production:

- [ ] Use HTTPS everywhere (required for WebAuthn)
- [ ] Set up proper session management
- [ ] Add rate limiting to auth endpoints
- [ ] Configure CORS properly
- [ ] Enable CSP headers
- [ ] Set up monitoring and alerting
- [ ] Add error tracking (Sentry, etc.)
- [ ] Test recovery flows
- [ ] Set up database backups
- [ ] Review security best practices

---

## Troubleshooting

### "WebAuthn not supported"

**Solution:** WebAuthn requires HTTPS (or localhost). Use `ngrok` or deploy to a staging server with SSL.

### "Challenge expired"

**Solution:** Challenges expire after 5 minutes by default. Make sure the user completes authentication quickly, or increase the timeout.

### "Invalid signature"

**Solution:** This usually means:
1. Private key doesn't match public key (user registered from different device)
2. Challenge was tampered with
3. Time skew between client and server (check system clocks)

### "Handle already exists"

**Solution:** Each public key generates a unique handle. This error shouldn't occur unless there's a hash collision (astronomically unlikely) or database corruption.

---

## API Reference

See [docs/api-reference.md](../api-reference.md) for complete API documentation.

---

## Support

- **GitHub Issues:** [github.com/Kelsidavis/authgrid](https://github.com/Kelsidavis/authgrid)
- **Discord:** Coming soon
- **Email:** support@authgrid.net (coming soon)

---

## Examples

Check out complete example applications:

- [examples/express-basic](../../examples/express-basic) — Minimal Express app
- [examples/nextjs-app](../../examples/nextjs-app) — Next.js with App Router
- [examples/python-flask](../../examples/python-flask) — Flask API
- [examples/go-server](../../examples/go-server) — Go HTTP server

---

*Welcome to passwordless authentication. You're going to love it.*
