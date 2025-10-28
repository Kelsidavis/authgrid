const express = require('express');
const session = require('express-session');
const fetch = require('node-fetch');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 3001;
const AUTHGRID_API = process.env.AUTHGRID_API || 'http://localhost:8080';

// Middleware
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(express.static('public'));

// Session configuration
app.use(session({
  secret: 'your-secret-key-change-in-production',
  resave: false,
  saveUninitialized: false,
  cookie: {
    secure: false, // Set to true in production with HTTPS
    maxAge: 24 * 60 * 60 * 1000 // 24 hours
  }
}));

// Middleware to check authentication
function requireAuth(req, res, next) {
  if (!req.session.authgridHandle) {
    return res.status(401).json({ error: 'Not authenticated' });
  }
  next();
}

// Routes

// Serve main page
app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

// Get current session info
app.get('/api/session', (req, res) => {
  if (req.session.authgridHandle) {
    res.json({
      authenticated: true,
      handle: req.session.authgridHandle,
      loginTime: req.session.loginTime
    });
  } else {
    res.json({ authenticated: false });
  }
});

// Register endpoint (proxy to Authgrid API)
app.post('/api/register', async (req, res) => {
  try {
    const response = await fetch(`${AUTHGRID_API}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(req.body)
    });

    const data = await response.json();

    if (!response.ok) {
      return res.status(response.status).json(data);
    }

    res.json(data);
  } catch (error) {
    console.error('Registration error:', error);
    res.status(500).json({ error: 'Registration failed' });
  }
});

// Challenge endpoint (proxy to Authgrid API)
app.post('/api/challenge', async (req, res) => {
  try {
    const response = await fetch(`${AUTHGRID_API}/challenge`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(req.body)
    });

    const data = await response.json();

    if (!response.ok) {
      return res.status(response.status).json(data);
    }

    res.json(data);
  } catch (error) {
    console.error('Challenge error:', error);
    res.status(500).json({ error: 'Challenge request failed' });
  }
});

// Verify endpoint (proxy to Authgrid API + create session)
app.post('/api/verify', async (req, res) => {
  try {
    const response = await fetch(`${AUTHGRID_API}/verify`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(req.body)
    });

    const data = await response.json();

    if (!response.ok) {
      return res.status(response.status).json(data);
    }

    // Create session
    if (data.verified) {
      req.session.authgridHandle = req.body.handle;
      req.session.authgridToken = data.token;
      req.session.loginTime = new Date().toISOString();

      // Save session
      req.session.save((err) => {
        if (err) {
          console.error('Session save error:', err);
          return res.status(500).json({ error: 'Session creation failed' });
        }

        res.json({
          verified: true,
          handle: req.body.handle,
          message: 'Authentication successful'
        });
      });
    } else {
      res.json(data);
    }
  } catch (error) {
    console.error('Verification error:', error);
    res.status(500).json({ error: 'Verification failed' });
  }
});

// Logout
app.post('/api/logout', (req, res) => {
  req.session.destroy((err) => {
    if (err) {
      return res.status(500).json({ error: 'Logout failed' });
    }
    res.json({ message: 'Logged out successfully' });
  });
});

// Protected route example
app.get('/api/profile', requireAuth, (req, res) => {
  res.json({
    handle: req.session.authgridHandle,
    loginTime: req.session.loginTime,
    message: 'This is a protected route!'
  });
});

// Protected route example - Dashboard
app.get('/api/dashboard', requireAuth, (req, res) => {
  res.json({
    handle: req.session.authgridHandle,
    data: {
      widgets: ['Analytics', 'Reports', 'Settings'],
      notifications: 3,
      lastActivity: req.session.loginTime
    }
  });
});

// Start server
app.listen(PORT, () => {
  console.log(`
╔══════════════════════════════════════════╗
║   Authgrid Express.js Example           ║
╠══════════════════════════════════════════╣
║                                          ║
║   Server: http://localhost:${PORT}        ║
║   Authgrid API: ${AUTHGRID_API}          ║
║                                          ║
║   Endpoints:                             ║
║   - GET  /                               ║
║   - POST /api/register                   ║
║   - POST /api/challenge                  ║
║   - POST /api/verify                     ║
║   - POST /api/logout                     ║
║   - GET  /api/session                    ║
║   - GET  /api/profile (protected)        ║
║   - GET  /api/dashboard (protected)      ║
║                                          ║
╚══════════════════════════════════════════╝
  `);
});
