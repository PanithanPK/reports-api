# Deployment Guide

## Production Deployment

### 1. Server Preparation

#### System Requirements
- **OS**: Ubuntu 20.04 LTS or later
- **RAM**: Minimum 2GB
- **Storage**: Minimum 10GB
- **Network**: Static IP address

#### Install Dependencies
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install MySQL
sudo apt install mysql-server -y
sudo mysql_secure_installation

# Install Nginx
sudo apt install nginx -y

# Install Git
sudo apt install git -y
```

### 2. Application Setup

#### Create Application User
```bash
sudo useradd -m -s /bin/bash reports
sudo mkdir -p /opt/reports-api
sudo chown reports:reports /opt/reports-api
```

#### Clone and Build Application
```bash
sudo -u reports git clone <repository-url> /opt/reports-api
cd /opt/reports-api
sudo -u reports go mod tidy
sudo -u reports go build -o reports-api main.go
```

#### Environment Configuration
```bash
sudo -u reports cp .env.prod .env
sudo -u reports nano .env
```

### 3. Database Setup

#### Create Database and User
```sql
CREATE DATABASE reports CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'reports_user'@'localhost' IDENTIFIED BY 'strong_password_here';
GRANT ALL PRIVILEGES ON reports.* TO 'reports_user'@'localhost';
FLUSH PRIVILEGES;
```

#### Import Schema (if available)
```bash
mysql -u reports_user -p reports < database_schema.sql
```

### 4. Systemd Service

#### Create Service File
```bash
sudo nano /etc/systemd/system/reports-api.service
```

```ini
[Unit]
Description=Reports API Service
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=reports
Group=reports
WorkingDirectory=/opt/reports-api
ExecStart=/opt/reports-api/reports-api prod
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/reports-api

# Environment
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

#### Enable Service
```bash
sudo systemctl daemon-reload
sudo systemctl enable reports-api
sudo systemctl start reports-api
sudo systemctl status reports-api
```

### 5. Nginx Configuration

#### Create Virtual Host
```bash
sudo nano /etc/nginx/sites-available/reports-api
```

```nginx
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;
    
    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com www.your-domain.com;
    
    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # Security Headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    # Gzip Compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
    
    # Rate Limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    
    location / {
        limit_req zone=api burst=20 nodelay;
        
        proxy_pass http://127.0.0.1:5001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # Static files (if available)
    location /static/ {
        alias /opt/reports-api/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # Health check endpoint
    location /health {
        access_log off;
        proxy_pass http://127.0.0.1:5001/;
    }
}
```

#### Enable Site
```bash
sudo ln -s /etc/nginx/sites-available/reports-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 6. SSL Certificate (Let's Encrypt)

#### Install Certbot
```bash
sudo apt install certbot python3-certbot-nginx -y
```

#### Create Certificate
```bash
sudo certbot --nginx -d your-domain.com -d www.your-domain.com
```

#### Setup Auto-renewal
```bash
sudo crontab -e
```

Add line:
```
0 12 * * * /usr/bin/certbot renew --quiet
```

### 7. Firewall Configuration

#### Setup UFW
```bash
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 'Nginx Full'
sudo ufw enable
```

### 8. Monitoring and Logging

#### Setup Log Rotation
```bash
sudo nano /etc/logrotate.d/reports-api
```

```
/var/log/reports-api/*.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 644 reports reports
    postrotate
        systemctl reload reports-api
    endscript
}
```

#### Install Monitoring Tools
```bash
# htop for system monitoring
sudo apt install htop -y

# fail2ban for brute force protection
sudo apt install fail2ban -y
```

### 9. Backup Strategy

#### Create Backup Script
```bash
sudo nano /opt/backup-reports.sh
```

```bash
#!/bin/bash
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup Database
mysqldump -u reports_user -p'password' reports > $BACKUP_DIR/db_backup_$DATE.sql

# Backup Application
tar -czf $BACKUP_DIR/app_backup_$DATE.tar.gz /opt/reports-api

# Remove old backups (older than 30 days)
find $BACKUP_DIR -name "*.sql" -mtime +30 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete

echo "Backup completed: $DATE"
```

#### Setup Cron Job for Backup
```bash
sudo chmod +x /opt/backup-reports.sh
sudo crontab -e
```

Add line:
```
0 2 * * * /opt/backup-reports.sh >> /var/log/backup.log 2>&1
```

### 10. Deployment Updates

#### Create Deploy Script
```bash
sudo nano /opt/deploy-reports.sh
```

```bash
#!/bin/bash
cd /opt/reports-api

echo "Starting deployment..."

# Backup current version
cp reports-api reports-api.backup

# Pull latest code
sudo -u reports git pull origin main

# Build new version
sudo -u reports go build -o reports-api.new main.go

# Test new version
if ./reports-api.new --version; then
    echo "Build successful"
    mv reports-api.new reports-api
    sudo systemctl restart reports-api
    echo "Deployment completed"
else
    echo "Build failed, rolling back"
    mv reports-api.backup reports-api
    exit 1
fi
```

### 11. Health Checks

#### Create Health Check Script
```bash
sudo nano /opt/health-check.sh
```

```bash
#!/bin/bash
ENDPOINT="http://localhost:5001/"
EXPECTED_STATUS="200"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" $ENDPOINT)

if [ "$STATUS" = "$EXPECTED_STATUS" ]; then
    echo "$(date): Health check passed"
    exit 0
else
    echo "$(date): Health check failed - Status: $STATUS"
    # Restart service
    sudo systemctl restart reports-api
    exit 1
fi
```

#### Setup Cron Job for Health Check
```bash
sudo chmod +x /opt/health-check.sh
sudo crontab -e
```

Add line:
```
*/5 * * * * /opt/health-check.sh >> /var/log/health-check.log 2>&1
```

## Docker Deployment

### 1. Docker Compose Setup
```yaml
version: '3.8'

services:
  reports-api:
    build: .
    ports:
      - "5001:5001"
    environment:
      - DB_HOST=mysql
      - DB_USER=reports_user
      - DB_PASSWORD=strong_password
      - DB_NAME=reports
    depends_on:
      - mysql
    restart: unless-stopped
    volumes:
      - ./uploads:/app/uploads

  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=root_password
      - MYSQL_DATABASE=reports
      - MYSQL_USER=reports_user
      - MYSQL_PASSWORD=strong_password
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - reports-api
    restart: unless-stopped

volumes:
  mysql_data:
```

### 2. Deploy with Docker
```bash
docker-compose up -d
```

## Post-Deployment Verification

### 1. Check Service Status
```bash
sudo systemctl status reports-api
sudo systemctl status nginx
sudo systemctl status mysql
```

### 2. Check Logs
```bash
sudo journalctl -u reports-api -f
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### 3. Test API
```bash
curl -k https://your-domain.com/
curl -k https://your-domain.com/api/v1/problem/list
```

### 4. Verify SSL
```bash
openssl s_client -connect your-domain.com:443 -servername your-domain.com
```