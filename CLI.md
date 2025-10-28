# Authgrid CLI

Command-line tool for Authgrid passwordless authentication.

## Installation

### Build from source

```bash
make build-cli
```

This creates the `./authgrid` binary in the project root.

### Install system-wide (optional)

```bash
make install-cli
# or manually:
sudo cp authgrid /usr/local/bin/
```

## Usage

```bash
authgrid [command] [flags]
```

## Commands

### `authgrid register`

Register a new user and receive a handle.

**Example:**
```bash
$ authgrid register
Registering new user...

✅ Registration successful!
   Handle: c4af5d15cd@authgrid.net
   Keystore: /home/user/.authgrid

To login:
   authgrid login --handle c4af5d15cd@authgrid.net
```

**What happens:**
1. Generates an Ed25519 keypair
2. Sends public key to Authgrid API
3. Receives your unique handle
4. Saves keypair to `~/.authgrid/` directory

**Security:**
- Private key never leaves your machine
- Stored with 0600 permissions (user-only access)
- Public key is sent to the server

---

### `authgrid login --handle <handle>`

Authenticate with a registered handle.

**Example:**
```bash
$ authgrid login --handle c4af5d15cd@authgrid.net
Logging in as c4af5d15cd@authgrid.net...

✅ Login successful!
   Handle: c4af5d15cd@authgrid.net
   Token: YzRhZjVkMTVjZEBhdXRoZ3JpZC5uZXQ6+bGScgWM...
```

**What happens:**
1. Loads your private key from keystore
2. Requests a challenge from the API
3. Signs the challenge with your private key
4. Sends signature to API for verification
5. Receives authentication token

**Usage in scripts:**
```bash
# Store token in variable
TOKEN=$(authgrid login --handle c4af5d15cd@authgrid.net | grep "Token:" | awk '{print $2}')

# Use token in API calls
curl -H "Authorization: Bearer $TOKEN" https://api.example.com/profile
```

---

### `authgrid list`

List all handles stored in your keystore.

**Example:**
```bash
$ authgrid list
Stored handles (3):
  • c4af5d15cd@authgrid.net
  • f8a9b3c2d1@authgrid.net
  • 1234567890@authgrid.net

To login:
  authgrid login --handle <handle>
```

---

### `authgrid version`

Show CLI version.

**Example:**
```bash
$ authgrid version
authgrid-cli version 0.1.1-alpha
```

---

### `authgrid help`

Show help message with all commands.

---

## Flags

### Global Flags

These can be used with any command:

**`--api <URL>`**
- Authgrid API URL
- Default: `http://localhost:8080`
- Example: `authgrid register --api https://authgrid.example.com`

**`--keystore <DIR>`**
- Directory to store keypairs
- Default: `~/.authgrid`
- Example: `authgrid register --keystore /tmp/test-keys`

---

## Examples

### Basic registration and login

```bash
# 1. Register
authgrid register

# 2. Login (use the handle from step 1)
authgrid login --handle YOUR_HANDLE@authgrid.net
```

### Using a custom API

```bash
# Connect to production API
authgrid register --api https://api.authgrid.com

# Login to production
authgrid login --handle YOUR_HANDLE@authgrid.com --api https://api.authgrid.com
```

### Multiple handles

```bash
# Register multiple handles
authgrid register  # Gets handle1@authgrid.net
authgrid register  # Gets handle2@authgrid.net
authgrid register  # Gets handle3@authgrid.net

# List all handles
authgrid list

# Login with specific handle
authgrid login --handle handle2@authgrid.net
```

### Custom keystore

```bash
# Use a different keystore directory
authgrid register --keystore ~/.my-custom-keys
authgrid login --handle HANDLE@authgrid.net --keystore ~/.my-custom-keys
```

### Automation scripts

```bash
#!/bin/bash
# auto-login.sh

HANDLE="your-handle@authgrid.net"

# Login and extract token
OUTPUT=$(authgrid login --handle $HANDLE)

if echo "$OUTPUT" | grep -q "Login successful"; then
    TOKEN=$(echo "$OUTPUT" | grep "Token:" | awk '{print $2}')
    echo "Authenticated! Token: $TOKEN"

    # Use the token
    curl -H "Authorization: Bearer $TOKEN" \
         http://localhost:3000/api/protected
else
    echo "Login failed!"
    exit 1
fi
```

---

## File Structure

### Keystore Directory

Default location: `~/.authgrid/`

```
~/.authgrid/
├── c4af5d15cd@authgrid.net.key
├── f8a9b3c2d1@authgrid.net.key
└── 1234567890@authgrid.net.key
```

### Key File Format

Each `.key` file contains:
- Line 1: Base64-encoded private key
- Line 2: Base64-encoded public key

**Example:**
```
MC4CAQAwBQYDK2VwBCIEIJ... (private key)
MCowBQYDK2VwAyEAGb9ECW... (public key)
```

**Security:**
- Files have 0600 permissions (read/write by owner only)
- Never share these files
- Back them up securely if needed

---

## Troubleshooting

### "Error loading keypair"

**Problem:** Handle not found in keystore

**Solution:**
```bash
# Check which handles you have
authgrid list

# Register if you haven't yet
authgrid register
```

### "Error connecting to API"

**Problem:** Can't reach the Authgrid API

**Solutions:**
1. Check if API is running:
   ```bash
   curl http://localhost:8080/health
   ```

2. Verify API URL:
   ```bash
   authgrid register --api http://localhost:8080
   ```

3. Check Docker containers:
   ```bash
   docker compose ps
   ```

### "API error: Handle not found"

**Problem:** Handle doesn't exist on server

**Solution:**
- Register first: `authgrid register`
- Use the correct handle from registration
- Check you're using the right API URL

### Permission denied on keystore

**Problem:** Can't write to keystore directory

**Solution:**
```bash
# Check permissions
ls -ld ~/.authgrid

# Fix permissions
chmod 700 ~/.authgrid
```

---

## Security Best Practices

### 1. Protect Your Keystore

```bash
# Ensure proper permissions
chmod 700 ~/.authgrid
chmod 600 ~/.authgrid/*.key
```

### 2. Back Up Your Keys

```bash
# Create encrypted backup
tar -czf authgrid-backup.tar.gz ~/.authgrid
gpg --encrypt authgrid-backup.tar.gz

# Store authgrid-backup.tar.gz.gpg in a safe place
```

### 3. Use Different Handles for Different Purposes

```bash
# Personal handle
authgrid register  # personal@authgrid.net

# Work handle
authgrid register --keystore ~/.authgrid-work  # work@authgrid.net

# Test handle
authgrid register --keystore /tmp/test-keys  # test@authgrid.net
```

### 4. Rotate Keys Periodically

```bash
# Generate new handle
authgrid register

# Update services to use new handle
# ...

# Remove old handle
rm ~/.authgrid/old-handle@authgrid.net.key
```

---

## Integration with Scripts

### Shell Scripts

```bash
#!/bin/bash
# Example: Automated deployment script with Authgrid auth

HANDLE="deploy@authgrid.net"

# Authenticate
echo "Authenticating..."
authgrid login --handle $HANDLE > /tmp/auth-output.txt

if grep -q "Login successful" /tmp/auth-output.txt; then
    echo "✅ Authenticated"

    # Extract token
    TOKEN=$(grep "Token:" /tmp/auth-output.txt | awk '{print $2}')

    # Deploy with authentication
    curl -H "Authorization: Bearer $TOKEN" \
         -X POST \
         https://api.example.com/deploy
else
    echo "❌ Authentication failed"
    exit 1
fi
```

### CI/CD Pipelines

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Setup Authgrid CLI
        run: |
          wget https://github.com/Kelsidavis/authgrid/releases/latest/download/authgrid
          chmod +x authgrid

      - name: Authenticate
        run: |
          # Restore keypair from secrets
          mkdir -p ~/.authgrid
          echo "${{ secrets.AUTHGRID_KEY }}" > ~/.authgrid/deploy@authgrid.net.key

          # Login
          ./authgrid login --handle deploy@authgrid.net
```

---

## Comparison with Web Demo

| Feature | CLI | Web Demo |
|---------|-----|----------|
| Registration | ✅ `authgrid register` | ✅ Click button |
| Login | ✅ `authgrid login` | ✅ Click button |
| Key Storage | ~/.authgrid/ directory | Browser localStorage |
| Automation | ✅ Perfect for scripts | ❌ Not scriptable |
| UI | Command-line only | ✅ Visual interface |
| Cross-device | Copy keystore manually | Separate per browser |

**Use CLI when:**
- Automating authentication in scripts
- CI/CD pipelines
- Server-side authentication
- DevOps workflows

**Use Web Demo when:**
- Interactive testing
- Demonstrating to users
- Browser-based applications

---

## Advanced Usage

### Custom API Integration

```bash
# Development
authgrid register --api http://localhost:8080

# Staging
authgrid register --api https://staging.authgrid.com

# Production
authgrid register --api https://api.authgrid.com
```

### Handle Migration

```bash
# Export old keystore
cp -r ~/.authgrid ~/.authgrid.backup

# Register on new server
authgrid register --api https://new-server.com

# Test new handle
authgrid login --handle NEW_HANDLE@authgrid.net --api https://new-server.com
```

### Multi-Environment Setup

```bash
# Development keystore
export DEV_KEYSTORE=~/.authgrid-dev
authgrid register --keystore $DEV_KEYSTORE --api http://localhost:8080

# Production keystore
export PROD_KEYSTORE=~/.authgrid-prod
authgrid register --keystore $PROD_KEYSTORE --api https://api.authgrid.com

# Use in scripts
authgrid login --handle $DEV_HANDLE --keystore $DEV_KEYSTORE --api http://localhost:8080
```

---

## Building from Source

```bash
# Clone repository
git clone https://github.com/Kelsidavis/authgrid
cd authgrid

# Build CLI
make build-cli

# Test it
./authgrid help

# Install (optional)
sudo cp authgrid /usr/local/bin/
```

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on:
- Reporting bugs
- Suggesting features
- Contributing code
- Testing guidelines

---

## License

MIT License - See [LICENSE](LICENSE) file

---

**Get started:** `authgrid register`

*Passwords die here.*
