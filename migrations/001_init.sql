-- Authgrid Database Schema

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    handle VARCHAR(255) UNIQUE NOT NULL,
    public_key TEXT NOT NULL,
    key_type VARCHAR(50) NOT NULL DEFAULT 'ed25519',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_login TIMESTAMP,
    metadata JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX idx_users_handle ON users(handle);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Challenges table
CREATE TABLE IF NOT EXISTS challenges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    handle VARCHAR(255) NOT NULL,
    challenge TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_challenges_handle ON challenges(handle);
CREATE INDEX idx_challenges_expires_at ON challenges(expires_at);
CREATE INDEX idx_challenges_cleanup ON challenges(expires_at) WHERE used = FALSE;

-- Sessions table (optional, for token management)
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    token TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP,
    metadata JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Function to clean up expired challenges
CREATE OR REPLACE FUNCTION cleanup_expired_challenges()
RETURNS void AS $$
BEGIN
    DELETE FROM challenges WHERE expires_at < NOW() - INTERVAL '1 hour';
END;
$$ LANGUAGE plpgsql;

-- Comments for documentation
COMMENT ON TABLE users IS 'Stores user public keys and handles';
COMMENT ON TABLE challenges IS 'Stores authentication challenges (ephemeral)';
COMMENT ON TABLE sessions IS 'Stores authentication sessions and tokens';
COMMENT ON COLUMN users.handle IS 'Unique email-like identifier (e.g., abc123@authgrid.net)';
COMMENT ON COLUMN users.public_key IS 'Base64-encoded Ed25519 public key';
COMMENT ON COLUMN challenges.challenge IS 'Base64-encoded random challenge bytes';
COMMENT ON COLUMN challenges.used IS 'Whether this challenge has been used for authentication';
