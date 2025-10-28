# Authgrid FAQ

## General Questions

### What is Authgrid?

Authgrid is a passwordless authentication system that uses cryptographic keys instead of passwords. Users get an email-like handle (e.g., `n5xtc4q3f3@authgrid.net`) that works across any service using Authgrid — no inbox, no password, no recovery spam.

### How is this different from passwords?

With passwords, you memorize a secret that's stored on the server (hashed, hopefully). If the server is breached, your password might be compromised.

With Authgrid, you have a private key on your device (never shared) and a public key on the server. You prove you own the private key without revealing it. Even if the server is breached, attackers only get public keys — which are useless without the private keys on your device.

### How is this different from WebAuthn?

Authgrid uses WebAuthn under the hood for browsers, but adds a critical feature: **cross-site identity**.

With raw WebAuthn, each site has a separate credential. With Authgrid, you have one handle that works everywhere — like email addresses work for many services, but with cryptographic security.

### How is this different from "Sign in with Google"?

With social login, Google (or Facebook, etc.) tracks where you log in and can revoke your access.

With Authgrid, you own your identity. No one can track you, no one can revoke your access (except the specific service you're using), and your identity works even if Authgrid goes offline (thanks to federation).

### Is this blockchain/crypto/web3?

No. Authgrid uses established cryptography (Ed25519, WebAuthn) but no blockchain, no cryptocurrency, no NFTs. It's just clean, proven cryptography applied to authentication.

---

## Technical Questions

### What happens if I lose my device?

You have several recovery options:

1. **Social recovery** — Trusted contacts can help you recover (N-of-M approval)
2. **Seed phrase backup** — Write down 12-24 words during setup
3. **Device-to-device transfer** — Transfer keys when setting up a new phone
4. **Custodial backup** — Organizations can offer optional backup services

You choose which recovery methods to enable during registration.

### Can Authgrid see my private keys?

No. Your private keys are generated on your device and never leave it. They're stored in your device's secure enclave (Secure Enclave on iOS, KeyStore on Android, TPM on Windows).

### What if Authgrid the company goes away?

Authgrid is an open protocol. Anyone can run an Authgrid node (like email servers). If the main authgrid.net service disappeared:

1. Your keys are still on your device
2. Services could run their own Authgrid nodes
3. You could migrate to another Authgrid provider
4. It's open source — the community could fork it

This is why federation is built-in from day one.

### Does this work offline?

Partially. You can:
- Sign challenges offline (no internet needed)
- Provision keys via QR codes
- Use NFC for local authentication

You do need connectivity to verify signatures with the Authgrid server (or the service's server if they run verification locally).

### What about privacy?

Authgrid is privacy-first by design:

- Handles contain no personally identifying information
- No tracking across sites (each service verifies independently)
- Optional anonymous handles (no name, email, or phone required)
- Open protocol means no vendor lock-in

You control your identity and what information you share.

### Is this secure?

Security is the foundation:

- Private keys never leave your device
- Biometric unlock (not transmitted anywhere)
- Challenge-response prevents replay attacks
- HTTPS prevents man-in-the-middle attacks
- Even if the server is breached, attackers only get public keys

We're planning security audits, bug bounties, and SOC2 certification.

### What cryptographic algorithms does it use?

- **Ed25519** — Elliptic curve for signatures (fast, secure, widely used)
- **WebAuthn** — Browser standard for hardware-backed keys
- **HTTPS/TLS 1.3** — Transport security
- **Random challenges** — 256-bit cryptographically secure random challenges

All industry-standard, battle-tested cryptography.

---

## Developer Questions

### How do I integrate Authgrid into my app?

For JavaScript/Node.js:
```bash
npm install authgrid-client
```

```javascript
import { authgrid } from 'authgrid-client'

app.post('/login', authgrid.authenticate())
```

See [integration guides](docs/) for other languages and frameworks.

### Do I need to run my own server?

No. You can use the managed Authgrid service (authgrid.net). But you *can* self-host if you want — it's open source.

### How much does it cost?

**Free tier:**
- 10,000 monthly active users
- Standard handles (@authgrid.net)
- Community support

**Pro tier: $49/month**
- 100,000 MAUs
- Custom domain (@yourapp.com)
- Email support

**Enterprise: Custom pricing**
- Unlimited users, SLA, compliance features

See [BUSINESS.md](BUSINESS.md) for full pricing details.

### Can I use my own domain?

Yes! Pro and Enterprise tiers let you use custom domains like `@yourapp.com`. You'll need to:

1. Configure DNS records
2. Point to your Authgrid node (or we can host it for you)
3. Configure your API settings

### What frameworks/languages are supported?

**Currently planned:**
- JavaScript/TypeScript (Node.js, browsers)
- Go
- Python (Flask, Django)
- Rust
- Java/Kotlin (Android)
- Swift (iOS)

More coming soon. It's just HTTP/JSON, so any language that can make HTTP requests can integrate.

### How does federation work?

When you use a custom domain (e.g., `alice@company.com`):

1. Service receives the handle
2. Extracts domain: `company.com`
3. Looks up DNS TXT record: `_authgrid.company.com`
4. Finds the Authgrid node URL
5. Makes API calls to that node

Each organization can run their own node while staying interoperable.

---

## Comparison Questions

### Authgrid vs. Auth0/Okta?

| Feature | Authgrid | Auth0/Okta |
|---------|----------|------------|
| Cost | $49-$2k/mo | $500-$5k/mo |
| Setup time | 5 minutes | Hours to days |
| Vendor lock-in | None (open protocol) | High |
| Self-hostable | Yes | No (Okta) / Limited (Auth0) |
| Passwords | Never | Optional |
| Cross-site identity | Yes | No |

### Authgrid vs. Firebase Auth?

| Feature | Authgrid | Firebase Auth |
|---------|----------|---------------|
| Owner | Independent | Google |
| Data portability | Full | Limited |
| Offline | Partial | No |
| Custom backend | Any | Firebase only |
| Pricing | Transparent | Opaque (bundled) |

### Authgrid vs. AWS Cognito?

| Feature | Authgrid | Cognito |
|---------|----------|---------|
| Complexity | Simple | Complex |
| Developer experience | Excellent | Poor |
| Cloud lock-in | None | AWS only |
| Documentation | Clear | Confusing |
| Pricing | Simple | Complicated |

---

## Business Questions

### How do you make money if it's open source?

We make money from:
1. **Managed hosting** — Most people prefer not to run servers
2. **Enterprise features** — SSO bridges, compliance tools, support
3. **Professional services** — Migration help, custom integrations
4. **White-label licensing** — Large orgs buying turnkey solutions

See [BUSINESS.md](BUSINESS.md) for the full business model.

### Why should I trust a new authentication system?

Fair question. Here's why Authgrid is trustworthy:

1. **Open source** — Anyone can audit the code
2. **Open protocol** — Not dependent on one company
3. **Standard cryptography** — No custom crypto, only proven algorithms
4. **Federation** — You can run your own node
5. **Security audits** — Planned for Q2 2025
6. **Bug bounty** — Coming soon

Plus, the architecture means even we can't access your private keys.

### What's your long-term vision?

Replace passwords everywhere with cryptographic identity. Make authentication:
- More secure (cryptographic keys vs. passwords)
- More private (no tracking, no PII required)
- More accessible (works offline, low bandwidth)
- More open (federated, not controlled by one company)

Think of it as "email for authentication" — a simple, open, universal standard.

---

## Usage Questions

### Can I use this for my side project?

Yes! The free tier supports 10,000 monthly active users. Perfect for side projects.

### Can I use this for my startup?

Absolutely. Many startups use Authgrid to avoid building authentication from scratch. Start on the free tier, upgrade to Pro ($49/mo) as you grow.

### Can I use this in production?

We're working toward production-readiness with:
- Security audits (Q2 2025)
- SOC2 certification (Q3 2025)
- 99.99% SLA for enterprise
- Bug bounty program
- Regular penetration testing

Check the [ROADMAP.md](ROADMAP.md) for current status.

### Can I migrate from Auth0/Firebase/Cognito?

Yes! We're building migration tools:
```bash
authgrid migrate --from auth0 --output ./users.json
authgrid migrate --from firebase --project your-project
```

See [docs/migration/](docs/migration/) for guides. We also offer professional migration services for enterprises.

---

## Community Questions

### How can I contribute?

We'd love your help! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Ways to contribute:
- Code (features, bug fixes, tests)
- Documentation (guides, translations)
- Community (answering questions, writing tutorials)
- Design (UI/UX improvements)

### Is there a community forum/chat?

Coming soon:
- Discord server (for real-time chat)
- GitHub Discussions (for long-form questions)
- Reddit community (for news and discussions)

For now, use GitHub Issues for questions.

### Can I use Authgrid in my open source project?

Yes! Authgrid is MIT licensed, so you can use it in any project (commercial or open source).

### How can I stay updated?

- **GitHub:** Star the repo for updates
- **Twitter/X:** @authgrid (coming soon)
- **Blog:** blog.authgrid.net (coming soon)
- **Newsletter:** Subscribe at authgrid.net (coming soon)

---

## Still have questions?

- Open an issue on [GitHub](https://github.com/Kelsidavis/authgrid)
- Email us at hello@authgrid.net (coming soon)
- Check the [documentation](docs/)

---

*Passwords die here. Let's build the future together.*
