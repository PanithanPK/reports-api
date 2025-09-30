# Security Policy

## Supported Versions

We actively support the following versions of Reports API with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in Reports API, please report it responsibly.

### How to Report

1. **GitHub Security Advisories**: [Create a security advisory](https://github.com/PanithanPK/reports-api/security/advisories) (preferred)
2. **GitHub Issues**: For non-sensitive security discussions
3. **Direct Contact**: Contact [@PanithanPK](https://github.com/PanithanPK) directly

### What to Include

Please include the following information in your report:

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact assessment
- Affected versions
- Suggested fix (if available)
- Your contact information

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Timeline**: Critical issues within 30 days

## Security Best Practices

### Environment Variables

Create `.env` file (never commit to repository):

```bash
# Database Configuration
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_USER=your_user
DATABASE_PASSWORD=your_secure_password
DATABASE_NAME=reports_db
MYSQL_ROOT_PASSWORD=your_root_password
MYSQL_PASSWORD=your_mysql_password

# Application Configuration
PORT=8080
APP_ENV=prod

# Telegram Bot
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id

# MinIO/S3 Configuration
MINIO_ENDPOINT=your_endpoint
MINIO_ACCESS_KEY=your_access_key
MINIO_SECRET_KEY=your_secret_key

# Docker Registry
DOCKER_TOKEN=your_docker_token
```

### Database Security

- Use strong, randomly generated passwords (minimum 16 characters)
- Enable SSL/TLS connections with proper certificates
- Implement proper access controls and least privilege
- Regular backups with encryption
- Monitor database access logs
- Use parameterized queries to prevent SQL injection

### API Security

- Input validation and sanitization for all endpoints
- CORS properly configured for specific origins
- Security headers implementation
- File upload restrictions (100MB limit)
- Request timeouts configured (30 seconds)
- Use crypto/rand for session ID generation

## Security Implementation Status

### ✅ Implemented Security Measures

- **Application Limits**:
  - Memory limit: 384MB
  - CPU cores limited to 2
  - Request body limit: 100MB
  - Read/Write timeouts: 30 seconds
  - Panic recovery middleware

- **Development Security**:
  - Swagger UI configuration
  - Environment-based configuration
  - Structured logging framework

### 🔧 Required Security Headers

Implement these security headers:

```go
// Add to Fiber middleware
app.Use(func(c *fiber.Ctx) error {
    c.Set("X-Content-Type-Options", "nosniff")
    c.Set("X-Frame-Options", "DENY")
    c.Set("X-XSS-Protection", "1; mode=block")
    c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
    c.Set("Content-Security-Policy", "default-src 'self'")
    return c.Next()
})
```

## Dependencies Security

### Key Dependencies

- **Go**: 1.23.0+ (keep updated)
- **Fiber**: v2.52.5+ (web framework)
- **MySQL Driver**: v1.9.3+
- **Telegram Bot API**: v5.5.1+
- **MinIO**: v7.0.95+ (object storage)
- **Swagger**: v1.16.6+
- **Crypto**: golang.org/x/crypto v0.40.0

### Security Audit Commands

```bash
# Update all dependencies
go get -u ./...

# Verify module integrity
go mod verify

# Check for known vulnerabilities
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Alternative security scanning
go install github.com/sonatypecommunity/nancy@latest
go list -json -deps | nancy sleuth
```

### Dependency Monitoring

1. Enable GitHub Dependabot alerts
2. Regular dependency updates (monthly)
3. Monitor CVE databases
4. Use `go mod tidy` to clean unused dependencies

## Deployment Security

### Docker Security

- Use multi-stage builds
- Run as non-root user
- Scan images for vulnerabilities: `docker scan reports-api:latest`
- Keep base images updated
- Use `.dockerignore` to exclude sensitive files
- Use environment variables for sensitive data

### Production Security Checklist

#### High Priority
- [ ] Implement security headers
- [ ] Enable database SSL connections
- [ ] Use HTTPS/TLS encryption
- [ ] Set `APP_ENV=prod`
- [ ] Input validation and sanitization
- [ ] Use crypto/rand for session generation

#### Medium Priority
- [ ] Configure proper firewall rules
- [ ] Set up monitoring and alerting
- [ ] Regular security updates
- [ ] Backup and disaster recovery plan
- [ ] Implement rate limiting
- [ ] Add request logging and monitoring

### Secure Configuration Template

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  app:
    environment:
      - APP_ENV=prod
      - DATABASE_PASSWORD=${DATABASE_PASSWORD}
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
    secrets:
      - db_password
      - telegram_token
  
  mysql:
    environment:
      - MYSQL_ROOT_PASSWORD_FILE=/run/secrets/mysql_root_password
    secrets:
      - mysql_root_password

secrets:
  db_password:
    external: true
  telegram_token:
    external: true
  mysql_root_password:
    external: true
```

## Security Updates

Security updates will be:

- Released as patch versions
- Documented in [CHANGELOG.md](CHANGELOG.md)
- Announced via GitHub releases
- Tagged with security labels
- Include migration guides if needed

## Acknowledgments

We appreciate security researchers and users who report vulnerabilities responsibly. Contributors will be acknowledged in our security advisories (with permission).

## Contact

For security-related questions or concerns:

- **GitHub Security**: [Security Advisories](https://github.com/PanithanPK/reports-api/security/advisories)
- **GitHub Issues**: For general security questions (non-sensitive)
- **Maintainer**: [@PanithanPK](https://github.com/PanithanPK)
- **Repository**: [reports-api](https://github.com/PanithanPK/reports-api)

---

**Note**: This security policy is subject to updates. Please check regularly for the latest version.
