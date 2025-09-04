# Security Policy

## Supported Versions

We actively support the following versions of Reports API with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.13.x  | :white_check_mark: |
| 1.12.x  | :white_check_mark: |
| 1.11.x  | :x:                |
| < 1.11  | :x:                |

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

### For Users

- Always use the latest supported version
- Keep Go dependencies updated (`go mod tidy && go mod download`)
- Use environment variables for sensitive data
- Enable HTTPS in production
- Regular security audits
- Monitor application logs

### Environment Variables

Never commit sensitive information to the repository. Use environment variables for:

```bash
# Database Configuration
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_USER=your_user
DATABASE_PASSWORD=your_secure_password
DATABASE_NAME=reports_db

# Application Configuration
PORT=8080
APP_ENV=prod

# Telegram Bot (if used)
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id

# MinIO/S3 Configuration (if used)
MINIO_ENDPOINT=your_endpoint
MINIO_ACCESS_KEY=your_access_key
MINIO_SECRET_KEY=your_secret_key
```

### Database Security

- Use strong passwords for MySQL connections
- Enable SSL/TLS connections
- Implement proper access controls
- Regular backups with encryption
- Monitor database access logs
- Use connection pooling securely

### API Security

- Input validation and sanitization implemented
- CORS properly configured for specific origins
- Security headers enabled
- File upload restrictions (100MB limit)
- Request timeouts configured (30 seconds)

## Known Security Features

### Implemented Security Measures

- **Security Headers**:
  - `X-Content-Type-Options: nosniff`
  - `X-Frame-Options: DENY`
  - `X-XSS-Protection: 1; mode=block`
  - `Content-Type: application/json`

- **Application Security**:
  - Memory limit: 384MB
  - CPU cores limited to 2
  - Request body limit: 100MB
  - Read/Write timeouts: 30 seconds
  - Panic recovery middleware

- **Development Security**:
  - Swagger UI only available in development mode
  - Environment-based configuration
  - Structured logging with levels

## Dependencies Security

### Key Dependencies

- **Go**: 1.23.0+ (keep updated)
- **Fiber**: v2.52.5+ (web framework)
- **MySQL Driver**: v1.9.3+
- **Telegram Bot API**: v5.5.1+
- **MinIO**: v7.0.95+ (object storage)
- **Swagger**: v1.16.6+

### Security Recommendations

1. Regularly update dependencies: `go get -u ./...`
2. Run security audits: `go list -json -deps | nancy sleuth`
3. Use `go mod verify` to verify dependencies
4. Monitor for CVEs in used packages

## Deployment Security

### Docker Security

- Use multi-stage builds
- Run as non-root user
- Scan images for vulnerabilities
- Keep base images updated
- Use `.dockerignore` to exclude sensitive files

### Production Checklist

- [ ] Use HTTPS/TLS encryption
- [ ] Set `APP_ENV=prod`
- [ ] Configure proper firewall rules
- [ ] Enable database SSL connections
- [ ] Set up monitoring and alerting
- [ ] Regular security updates
- [ ] Backup and disaster recovery plan

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