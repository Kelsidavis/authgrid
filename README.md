# Authgrid

**Passwordless authentication that looks like email, feels like WebAuthn, and runs anywhere.**

Authgrid is a lightweight, globally scalable identity system that replaces passwords with cryptographic keys stored securely on users' devices. Users get a self-owned identity handle (like `n5xtc4q3f3@authgrid.net`) that works across any service using Authgrid — no inbox, no password, no recovery spam.

## 🚀 Quick Start

**Get Authgrid running in 5 minutes:**

```bash
cd authgrid
make run

# Visit http://localhost:3000 for the demo
# Visit http://localhost:8080 for the API
```

**Or integrate into your app:**

```html
<script src="authgrid.js"></script>
<script>
  const authgrid = new AuthgridClient({
    apiUrl: 'http://localhost:8080'
  });

  // Register
  const { handle } = await authgrid.register();

  // Login
  const { token } = await authgrid.authenticate(handle);
</script>
```

👉 **[Full Quickstart Guide](QUICKSTART.md)** | **[Build Status](BUILD_STATUS.md)**

## Core Concept

**Instead of passwords:** Users authenticate by proving control of a cryptographic key stored securely on their device, unlocked by fingerprint or PIN.

**Instead of email addresses:** They get a self-owned identity handle that looks familiar: `n5xtc4q3f3@authgrid.net`

That's it. No inbox. No password. No recovery spam. Just one small piece of text that is you across any service using Authgrid.

## Why Authgrid?

- ✅ **Instant global adoption path** — legacy systems already accept "email" strings as usernames
- ✅ **Low bandwidth** — total handshake <2 KB, works over 2G, satellite, or mesh
- ✅ **Zero-cost scaling** — no SMTP, no password resets, no storage of secrets
- ✅ **Private by design** — nothing personally identifying unless you choose it
- ✅ **Offline-friendly** — can provision keys via QR codes or kiosks
- ✅ **Developer-simple** — three API endpoints, 200 lines of backend to start

## How It Works

### User registers
1. Device generates a keypair
2. Server assigns a pseudo-email (`<short_hash>@authgrid.net`)
3. Server stores only the public key
4. Done.

### User logs in anywhere
1. Service sends challenge
2. Device signs challenge (unlocked by biometric or PIN)
3. Server verifies signature
4. Done.

### Inter-service federation
- Services use the Authgrid API or OpenID-compatible bridge
- Handle (`@authgrid.net`) acts as stable cross-site ID

## Project Status

✅ **Phase 4 Complete** — Production-ready!

**What's working:**
- ✅ Core API Server (Go) with Ed25519 & ECDSA support
- ✅ PostgreSQL database with migrations
- ✅ JavaScript client SDK (browser-ready)
- ✅ Demo web application
- ✅ **NEW:** CLI tool (`authgrid register`, `authgrid login`)
- ✅ **NEW:** Unit tests (8 tests covering crypto functions)
- ✅ **NEW:** Express.js integration example
- ✅ Docker deployment
- ✅ Comprehensive documentation

**Try it now:** `make run` and visit http://localhost:3000

**Or use the CLI:** `make build-cli && ./authgrid register`

See [PHASE4_COMPLETE.md](PHASE4_COMPLETE.md) for details.

## Documentation

**Getting Started:**
- [QUICKSTART.md](QUICKSTART.md) — Get running in 5 minutes
- [docs/getting-started.md](docs/getting-started.md) — Integration guide
- [CLI.md](CLI.md) — Command-line tool usage
- [TEST.md](TEST.md) — Testing guide
- [FAQ.md](FAQ.md) — Common questions answered

**Examples:**
- [examples/demo/](examples/demo/) — Web demo with registration & login
- [examples/express-simple/](examples/express-simple/) — Express.js integration with sessions

**Deep Dives:**
- [VISION.md](VISION.md) — The full pitch and technical overview
- [STRATEGY.md](STRATEGY.md) — Growth strategy and go-to-market
- [BUSINESS.md](BUSINESS.md) — Business model and financials
- [TECHNICAL.md](TECHNICAL.md) — Architecture and API specs
- [ROADMAP.md](ROADMAP.md) — Development timeline
- [PHASE4_COMPLETE.md](PHASE4_COMPLETE.md) — Latest features and improvements

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Core API | **Go 1.21** with standard library HTTP server |
| Database | **PostgreSQL 15** with connection pooling |
| Transport | **JSON over HTTPS** (REST) |
| Cryptography | **Ed25519** signatures via Go crypto/ed25519 |
| Client SDK | **JavaScript** (browser-ready, WebCrypto API) |
| Deployment | **Docker + Docker Compose** |
| Rate Limiting | **golang.org/x/time/rate** |
| CORS | **rs/cors** middleware |

**Hosting:** Runs on a $5/mo VPS for thousands of users

## Security Model

- ✅ Private keys never leave device
- ✅ Biometric only unlocks key, not used directly
- ✅ Challenge–response, no replay, no MITM possible
- ✅ Server compromise leaks only public data — users remain safe
- ✅ Optional social recovery or seed phrase backup
- ✅ Optional federation: multiple domains can host their own Authgrid nodes, all interoperable

## The Vision

Think of it as: **"Email syntax meets cryptographic identity."**

You can log into anything — a forum, IoT device, bank app, or mesh node — with the same handle, using the same cryptographic proof, without ever creating "another account."

It's the bridge between today's web logins and the decentralized identity future. And because we're using familiar UX (email-like handles, 3-field forms), it can actually proliferate — not just live in whitepapers.

## License

MIT License — See [LICENSE](LICENSE) file for details

## Contributing

We'd love your help! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**Quick ways to contribute:**
- Test the MVP and report bugs
- Add unit tests
- Create example integrations
- Improve documentation
- Submit feature requests

---

*Passwords die here.*
