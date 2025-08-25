# Installation Guide

## System Requirements

### Software Requirements
- **Go 1.23.0** or later
- **MySQL 8.0** or later
- **Git** for repository cloning
- **Docker** (optional)

### Hardware Requirements
- **RAM**: Minimum 512MB (recommended 1GB)
- **Storage**: Minimum 100MB for application
- **CPU**: 1 core (recommended 2 cores)

## Installation

### 1. Clone Repository
```bash
git clone <repository-url>
cd reports-api
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Database Setup

#### Create Database
```sql
CREATE DATABASE reports CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### Create User (if needed)
```sql
CREATE USER 'reports_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON reports.* TO 'reports_user'@'localhost';
FLUSH PRIVILEGES;
```

### 4. Environment Variables Setup

#### For Development
```bash
cp .env.dev .env
```

#### Edit .env file
```env
# Server Configuration
PORT=5001

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=reports_user
DB_PASSWORD=your_password
DB_NAME=reports

# JWT Configuration
JWT_SECRET=your_jwt_secret_key
TOKEN_HOUR_LIFESPAN=24

# Telegram Bot (optional)
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id

# MinIO/S3 Configuration (optional)
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=reports-files
```

### 5. Run Application

#### Development Mode
```bash
go run main.go dev
```

#### Production Mode
```bash
go run main.go prod
```

#### Default Mode
```bash
go run main.go
```

## Docker Installation

### 1. Build Docker Image
```bash
docker build -t reports-api .
```

### 2. Run Container
```bash
docker run -d \
  --name reports-api \
  -p 5001:5001 \
  --env-file .env \
  reports-api
```

### 3. Using Docker Compose (if available)
```bash
docker-compose up -d
```

## Additional Configuration

### SSL/TLS Configuration
For production, use a reverse proxy like Nginx:

```nginx
server {
    listen 443 ssl;
    server_name your-domain.com;
    
    ssl_certificate /path/to/certificate.crt;
    ssl_certificate_key /path/to/private.key;
    
    location / {
        proxy_pass http://localhost:5001;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Systemd Service (Linux)
Create file `/etc/systemd/system/reports-api.service`:

```ini
[Unit]
Description=Reports API Service
After=network.target

[Service]
Type=simple
User=reports
WorkingDirectory=/opt/reports-api
ExecStart=/opt/reports-api/reports-api prod
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable service:
```bash
sudo systemctl enable reports-api
sudo systemctl start reports-api
```

## Installation Verification

### 1. Health Check
```bash
curl http://localhost:5001/
```

Expected response:
```json
{
  "status": "OK",
  "version": "1.0.0"
}
```

### 2. Check Database Connection
Look for in logs:
```
[INFO] ✅ Database connection established
```

### 3. Test API Endpoints
```bash
curl http://localhost:5001/api/v1/problem/list
```

## Troubleshooting

### Common Issues

#### 1. Database Connection Error
```
Error: dial tcp 127.0.0.1:3306: connect: connection refused
```

**Solution:**
- Check if MySQL service is running
- Verify DB_HOST, DB_PORT in .env
- Check firewall settings

#### 2. Port Already in Use
```
Error: listen tcp :5001: bind: address already in use
```

**Solution:**
- Change PORT in .env
- Or stop process using that port

#### 3. Permission Denied
```
Error: permission denied
```

**Solution:**
- Check file and folder permissions
- Run with sudo (not recommended for production)

### Viewing Logs
```bash
# View real-time logs
tail -f /var/log/reports-api.log

# View Docker logs
docker logs -f reports-api
```

## Updates

### 1. Pull Latest from Git
```bash
git pull origin main
```

### 2. Update Dependencies
```bash
go mod tidy
```

### 3. Restart Service
```bash
# Systemd
sudo systemctl restart reports-api

# Docker
docker restart reports-api

# Manual
pkill reports-api
go run main.go prod
```

## Backup

### Database Backup
```bash
mysqldump -u reports_user -p reports > backup_$(date +%Y%m%d).sql
```

### Application Backup
```bash
tar -czf reports-api-backup-$(date +%Y%m%d).tar.gz \
  /opt/reports-api \
  /etc/systemd/system/reports-api.service
```