# Contributing to Authgrid

Thank you for your interest in contributing to Authgrid! We're building the future of passwordless authentication, and we'd love your help.

## How to Contribute

### Reporting Issues

- Check if the issue already exists
- Use the issue template (if available)
- Include steps to reproduce
- Specify your environment (OS, browser, versions)

### Suggesting Features

- Open an issue with the "feature request" label
- Explain the use case and motivation
- Consider how it fits with Authgrid's philosophy (simple, secure, scalable)

### Submitting Code

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes**
   - Follow the code style (see below)
   - Add tests for new functionality
   - Update documentation as needed
4. **Commit your changes**
   ```bash
   git commit -m "Add feature: description"
   ```
5. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```
6. **Open a Pull Request**
   - Describe what you changed and why
   - Reference any related issues
   - Ensure CI passes

## Code Style

### Go
- Follow standard Go conventions
- Run `gofmt` before committing
- Use `golangci-lint` for linting
- Keep functions small and focused

### Rust
- Follow Rust style guidelines
- Run `cargo fmt` before committing
- Run `cargo clippy` and fix warnings
- Write idiomatic Rust

### JavaScript/TypeScript
- Use Prettier for formatting
- Follow Airbnb style guide
- Use TypeScript for type safety
- Write JSDoc comments for public APIs

### Python
- Follow PEP 8
- Use Black for formatting
- Type hints for all function signatures
- Docstrings for all public functions

## Testing

- Write tests for all new features
- Ensure existing tests pass
- Aim for >80% code coverage
- Include integration tests where appropriate

**Running tests:**
```bash
# Go
go test ./...

# Rust
cargo test

# JavaScript
npm test

# Python
pytest
```

## Documentation

- Update README.md if needed
- Add inline code comments for complex logic
- Update API documentation
- Include examples for new features

## Security

**Please do not open public issues for security vulnerabilities.**

Instead, email security@authgrid.net with:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We'll acknowledge within 48 hours and work with you on a fix.

## Development Setup

### Prerequisites
- Go 1.21+ or Rust 1.70+
- PostgreSQL 14+
- Node.js 18+ (for client SDK)
- Docker (optional, for local testing)

### Setup
```bash
# Clone the repository
git clone https://github.com/Kelsidavis/authgrid.git
cd authgrid

# Install dependencies
make install

# Set up database
make db-setup

# Run tests
make test

# Start development server
make dev
```

See individual component READMEs for specific setup instructions.

## Code Review Process

1. **Automated checks** — CI must pass (tests, linting, formatting)
2. **Peer review** — At least one maintainer must approve
3. **Discussion** — Address feedback and comments
4. **Merge** — Squash and merge when approved

## Community Guidelines

- Be respectful and inclusive
- Assume good intentions
- Focus on ideas, not individuals
- Help others learn and grow
- Follow our [Code of Conduct](CODE_OF_CONDUCT.md)

## Areas We Need Help

### Core Development
- Performance optimization
- Security hardening
- New integrations (frameworks, languages)

### Documentation
- Tutorials and guides
- API reference improvements
- Translation to other languages

### Community
- Answering questions (GitHub, Discord, forums)
- Writing blog posts and articles
- Creating video tutorials

### Testing
- Writing more tests
- Load testing
- Security testing

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Acknowledged in release notes
- Invited to the contributors Discord channel

Significant contributions may result in:
- Maintainer status
- Speaking opportunities at conferences
- Authgrid swag

## Questions?

- **Discord:** [Join our community](https://discord.gg/authgrid) (coming soon)
- **Email:** contribute@authgrid.net (coming soon)
- **Discussions:** Use GitHub Discussions for questions

---

Thank you for helping build the future of authentication! 🔐
