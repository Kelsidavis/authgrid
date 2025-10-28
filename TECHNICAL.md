# Authgrid Technical Architecture

## System Overview

Authgrid is a distributed authentication system that uses public-key cryptography to verify user identity without passwords. The system consists of three main components:

1. **Core API Server** — Handles registration, challenge generation, and signature verification
2. **Client SDKs** — Manage keypair generation and signing on user devices
3. **Federation Layer** — Enables cross-domain identity and interoperability

## Core Concepts

### Identity Handles

Format: `<identifier>@<domain>`

Example: `n5xtc4q3f3@authgrid.net`

**Structure:**
- Identifier: 10-character alphanumeric hash (base36 or bech32)
- Domain: The Authgrid node domain

**Generation:**
```
identifier = base36(sha256(public_key))[0:10]
handle = identifier + "@" + domain
```

**Properties:**
- Collision-resistant (10 chars base36 = 3.65 quadrillion possibilities)
- Looks like email (fits existing form fields)
- No PII (nothing personally identifying)
- Globally unique (when combined with domain)
- Short enough to type (if needed)

### Cryptographic Keys

**Key Types Supported:**

1. **Ed25519** (recommended for raw implementation)
   - 256-bit elliptic curve
   - Fast signing and verification
   - Small signature size (64 bytes)
   - Industry standard (used by Signal, Tor, SSH)

2. **WebAuthn** (recommended for browsers)
   - Uses device's secure enclave (TPM, Secure Enclave, etc.)
   - Biometric unlock (fingerprint, face)
   - Best user experience
   - Platform-native security

**Key Storage:**
- Private key: NEVER leaves device, stored in secure enclave
- Public key: Stored in Authgrid database, publicly readable

---

## Authentication Flow

### Registration

```
┌─────────┐                ┌──────────────┐                ┌─────────────┐
│ Client  │                │  Authgrid    │                │  Database   │
│ Device  │                │     API      │                │             │
└────┬────┘                └──────┬───────┘                └──────┬──────┘
     │                            │                               │
     │  1. Generate keypair       │                               │
     ├──────────────────────────► │                               │
     │     (public_key)            │                               │
     │                            │                               │
     │                            │  2. Generate handle           │
     │                            ├──────────────────────────────►│
     │                            │     Store public_key          │
     │                            │                               │
     │  3. Return handle          │                               │
     │◄───────────────────────────┤                               │
     │    (n5xtc4q3f3@authgrid.net)                              │
     │                            │                               │
     │  4. Store private key      │                               │
     │     in secure enclave      │                               │
     │                            │                               │
```

**API Endpoint:**
```
POST /register
Content-Type: application/json

{
  "public_key": "base64_encoded_public_key",
  "key_type": "ed25519" | "webauthn",
  "attestation": "optional_webauthn_attestation"
}

Response:
{
  "handle": "n5xtc4q3f3@authgrid.net",
  "id": "uuid",
  "created_at": "2025-01-15T10:30:00Z"
}
```

### Authentication

```
┌─────────┐          ┌──────────────┐          ┌─────────────┐
│ Client  │          │  Authgrid    │          │  Database   │
│ Device  │          │     API      │          │             │
└────┬────┘          └──────┬───────┘          └──────┬──────┘
     │                      │                         │
     │  1. Request login    │                         │
     ├─────────────────────►│                         │
     │   (handle)           │                         │
     │                      │                         │
     │                      │  2. Lookup public_key   │
     │                      ├────────────────────────►│
     │                      │                         │
     │  3. Return challenge │                         │
     │◄─────────────────────┤                         │
     │   (random 32 bytes)  │                         │
     │                      │                         │
     │  4. Sign challenge   │                         │
     │     with private key │                         │
     │                      │                         │
     │  5. Submit signature │                         │
     ├─────────────────────►│                         │
     │   (signature)        │                         │
     │                      │                         │
     │                      │  6. Verify signature    │
     │                      │     with public_key     │
     │                      │                         │
     │  7. Return token     │                         │
     │◄─────────────────────┤                         │
     │   (JWT or session)   │                         │
     │                      │                         │
```

**API Endpoints:**

```
POST /challenge
Content-Type: application/json

{
  "handle": "n5xtc4q3f3@authgrid.net"
}

Response:
{
  "challenge": "base64_encoded_random_bytes",
  "expires_at": "2025-01-15T10:35:00Z"
}
```

```
POST /verify
Content-Type: application/json

{
  "handle": "n5xtc4q3f3@authgrid.net",
  "challenge": "base64_encoded_challenge",
  "signature": "base64_encoded_signature"
}

Response:
{
  "verified": true,
  "token": "jwt_token_here",
  "expires_at": "2025-01-15T12:30:00Z"
}
```

---

## Data Models

### User Record

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  handle VARCHAR(255) UNIQUE NOT NULL,
  public_key TEXT NOT NULL,
  key_type VARCHAR(50) NOT NULL, -- 'ed25519' or 'webauthn'
  attestation JSONB, -- WebAuthn attestation data
  created_at TIMESTAMP DEFAULT NOW(),
  last_login TIMESTAMP,
  metadata JSONB, -- Extensible field for custom data

  INDEX idx_handle (handle),
  INDEX idx_created_at (created_at)
);
```

### Challenge Record

```sql
CREATE TABLE challenges (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  handle VARCHAR(255) NOT NULL,
  challenge TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  expires_at TIMESTAMP NOT NULL,
  used BOOLEAN DEFAULT FALSE,

  INDEX idx_handle (handle),
  INDEX idx_expires_at (expires_at)
);

-- Clean up expired challenges periodically
CREATE INDEX idx_cleanup ON challenges (expires_at) WHERE used = FALSE;
```

### Session Record (optional)

```sql
CREATE TABLE sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  token TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  expires_at TIMESTAMP NOT NULL,
  last_activity TIMESTAMP,
  metadata JSONB, -- Device info, IP, etc.

  INDEX idx_token (token),
  INDEX idx_user_id (user_id),
  INDEX idx_expires_at (expires_at)
);
```

---

## Security Considerations

### Threat Model

**Protected against:**
- ✅ Server breach (only public keys stored)
- ✅ Man-in-the-middle (HTTPS + challenge-response)
- ✅ Replay attacks (challenges expire, single-use)
- ✅ Credential stuffing (no passwords to stuff)
- ✅ Phishing (no secrets to phish)
- ✅ Brute force (cryptographic keys, not guessable)

**Still vulnerable to:**
- ⚠️ Device compromise (if attacker controls device)
- ⚠️ Social engineering (tricking user to approve malicious challenge)
- ⚠️ Physical device theft (mitigated by biometric/PIN lock)

### Challenge Security

**Requirements:**
- Must be cryptographically random (use `crypto.getRandomValues()` or equivalent)
- Must be at least 32 bytes (256 bits)
- Must expire (recommend 5 minutes)
- Must be single-use (mark as used after verification)

**Implementation:**
```go
func generateChallenge() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(bytes), nil
}
```

### Rate Limiting

**Recommended limits:**
- Registration: 5 per IP per hour
- Challenge requests: 10 per handle per minute
- Verify attempts: 5 per challenge (then invalidate)
- API calls: 100 per API key per minute (authenticated)

### Transport Security

- ✅ TLS 1.3 required
- ✅ HSTS headers enforced
- ✅ Certificate pinning recommended for mobile apps
- ✅ No mixed content allowed

---

## Client SDK Architecture

### JavaScript/TypeScript SDK

```typescript
import { AuthgridClient } from 'authgrid-client';

const client = new AuthgridClient({
  apiUrl: 'https://api.authgrid.net',
  domain: 'authgrid.net'
});

// Registration
async function register() {
  const result = await client.register();
  console.log(`Your handle: ${result.handle}`);
  // Store result.handle in localStorage or database
}

// Login
async function login(handle: string) {
  const token = await client.authenticate(handle);
  console.log(`Logged in! Token: ${token}`);
  // Use token for authenticated requests
}
```

**Internal flow:**
```typescript
class AuthgridClient {
  async register(): Promise<{ handle: string }> {
    // 1. Generate keypair using WebAuthn or WebCrypto
    const keypair = await this.generateKeypair();

    // 2. Store private key in IndexedDB or WebAuthn credential
    await this.storePrivateKey(keypair.privateKey);

    // 3. Send public key to server
    const response = await fetch(`${this.apiUrl}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        public_key: keypair.publicKey,
        key_type: 'webauthn'
      })
    });

    const data = await response.json();
    return { handle: data.handle };
  }

  async authenticate(handle: string): Promise<string> {
    // 1. Request challenge
    const challengeResp = await fetch(`${this.apiUrl}/challenge`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ handle })
    });
    const { challenge } = await challengeResp.json();

    // 2. Sign challenge with private key
    const signature = await this.signChallenge(challenge);

    // 3. Verify signature
    const verifyResp = await fetch(`${this.apiUrl}/verify`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ handle, challenge, signature })
    });
    const { token } = await verifyResp.json();

    return token;
  }
}
```

---

## Federation Architecture

### Multi-Domain Support

Each organization can run their own Authgrid node:

- `alice@company.com` (Company's Authgrid node)
- `bob@authgrid.net` (Main Authgrid node)
- `charlie@bank.io` (Bank's Authgrid node)

### Discovery Protocol

When a service receives a handle, it needs to find the Authgrid node:

```
1. Extract domain from handle: "company.com"
2. Lookup DNS TXT record: _authgrid.company.com
3. Record contains API endpoint: https://auth.company.com
4. Make API calls to that endpoint
```

**DNS Configuration:**
```
_authgrid.company.com. IN TXT "v=authgrid1 api=https://auth.company.com"
```

### Cross-Domain Authentication

```
┌─────────┐          ┌──────────────┐          ┌─────────────┐
│  App    │          │  Authgrid    │          │  Company    │
│ (relying│          │  Federation  │          │  Authgrid   │
│  party) │          │    Hub       │          │    Node     │
└────┬────┘          └──────┬───────┘          └──────┬──────┘
     │                      │                         │
     │  1. Login request    │                         │
     ├─────────────────────►│                         │
     │   (alice@company.com)│                         │
     │                      │                         │
     │                      │  2. Discover node       │
     │                      ├────────────────────────►│
     │                      │     (DNS lookup)        │
     │                      │                         │
     │                      │  3. Forward challenge   │
     │                      ├────────────────────────►│
     │                      │                         │
     │  4. Return challenge │                         │
     │◄─────────────────────┤                         │
     │                      │                         │
     │  5. Sign & verify    │  6. Verify signature    │
     ├─────────────────────►├────────────────────────►│
     │                      │                         │
     │  7. Return token     │                         │
     │◄─────────────────────┤                         │
     │                      │                         │
```

---

## Deployment Architecture

### Minimal Setup (Single Server)

```
┌─────────────────────────────────────────┐
│           Authgrid Server               │
│                                         │
│  ┌────────────┐       ┌──────────────┐ │
│  │   API      │──────►│  PostgreSQL  │ │
│  │  (Go/Rust) │       │              │ │
│  └────────────┘       └──────────────┘ │
│         │                               │
│         │                               │
│  ┌────────────┐                         │
│  │   Redis    │                         │
│  │  (cache)   │                         │
│  └────────────┘                         │
└─────────────────────────────────────────┘

Hardware: $5-20/mo VPS
Capacity: 10k-100k users
Latency: <100ms globally (with CDN)
```

### Production Setup (High Availability)

```
                    ┌──────────────┐
                    │  CloudFlare  │
                    │   /  CDN     │
                    └───────┬──────┘
                            │
              ┌─────────────┴─────────────┐
              │                           │
     ┌────────▼───────┐         ┌────────▼───────┐
     │  API Server 1  │         │  API Server 2  │
     │   (Go/Rust)    │         │   (Go/Rust)    │
     └────────┬───────┘         └────────┬───────┘
              │                           │
              └─────────────┬─────────────┘
                            │
              ┌─────────────▼─────────────┐
              │    PostgreSQL Cluster     │
              │   (Primary + Replicas)    │
              └───────────────────────────┘
                            │
              ┌─────────────▼─────────────┐
              │      Redis Cluster        │
              │  (Challenge cache + rate  │
              │       limiting)           │
              └───────────────────────────┘

Hardware: ~$200/mo cloud infrastructure
Capacity: 1M+ users
Latency: <50ms globally
Uptime: 99.99%
```

### Enterprise Deployment (On-Premise)

```
┌─────────────────────────────────────────────┐
│          Customer's Data Center             │
│                                             │
│  ┌────────────────────────────────────────┐ │
│  │       Authgrid Container Cluster       │ │
│  │                                        │ │
│  │  ┌──────────┐  ┌──────────┐          │ │
│  │  │   API    │  │   API    │          │ │
│  │  └────┬─────┘  └────┬─────┘          │ │
│  │       │             │                │ │
│  │  ┌────▼─────────────▼────┐           │ │
│  │  │    PostgreSQL HA      │           │ │
│  │  └───────────────────────┘           │ │
│  └────────────────────────────────────────┘ │
│                                             │
│  ┌────────────────────────────────────────┐ │
│  │   LDAP/AD Integration Bridge          │ │
│  └────────────────────────────────────────┘ │
│                                             │
│  ┌────────────────────────────────────────┐ │
│  │   Audit Logging & Compliance          │ │
│  └────────────────────────────────────────┘ │
└─────────────────────────────────────────────┘
```

---

## Performance Targets

### Latency
- Registration: <500ms
- Challenge generation: <50ms
- Signature verification: <100ms
- End-to-end login: <1s

### Throughput
- Challenge requests: 10k/sec per server
- Verification: 5k/sec per server
- Registration: 1k/sec per server

### Storage
- Per user: ~500 bytes (handle + public key + metadata)
- 1M users: ~500 MB
- 10M users: ~5 GB
- Very efficient scaling

---

## API Reference

### Core Endpoints

#### POST /register
Register a new user and get a handle.

**Request:**
```json
{
  "public_key": "base64_string",
  "key_type": "ed25519" | "webauthn",
  "attestation": {} // Optional WebAuthn data
}
```

**Response: 201 Created**
```json
{
  "handle": "n5xtc4q3f3@authgrid.net",
  "id": "uuid",
  "created_at": "ISO8601_timestamp"
}
```

**Errors:**
- 400: Invalid public key format
- 429: Rate limit exceeded

---

#### POST /challenge
Request an authentication challenge.

**Request:**
```json
{
  "handle": "n5xtc4q3f3@authgrid.net"
}
```

**Response: 200 OK**
```json
{
  "challenge": "base64_string",
  "expires_at": "ISO8601_timestamp"
}
```

**Errors:**
- 404: Handle not found
- 429: Rate limit exceeded

---

#### POST /verify
Verify a signed challenge.

**Request:**
```json
{
  "handle": "n5xtc4q3f3@authgrid.net",
  "challenge": "base64_string",
  "signature": "base64_string"
}
```

**Response: 200 OK**
```json
{
  "verified": true,
  "token": "jwt_token",
  "expires_at": "ISO8601_timestamp"
}
```

**Errors:**
- 400: Invalid signature
- 404: Challenge not found or expired
- 429: Too many verification attempts

---

#### GET /user/:handle
Get public user information (optional endpoint).

**Response: 200 OK**
```json
{
  "handle": "n5xtc4q3f3@authgrid.net",
  "public_key": "base64_string",
  "created_at": "ISO8601_timestamp",
  "metadata": {} // Optional public metadata
}
```

---

## OpenID Connect Bridge

For compatibility with existing OAuth/OIDC systems:

### Discovery Endpoint
`GET /.well-known/openid-configuration`

### Authorization Endpoint
`GET /authorize?client_id=...&redirect_uri=...&response_type=code`

### Token Endpoint
`POST /token`

This allows Authgrid to act as an identity provider for services that expect OAuth/OIDC.

---

## Next Steps

See [VISION.md](VISION.md) for the product vision and [STRATEGY.md](STRATEGY.md) for the roadmap.

Ready to start building? Check out the implementation guides in `/docs` (coming soon).
