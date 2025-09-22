# Nginx Setup Guide

## Overview
This guide covers Nginx configuration for the Reports API, including reverse proxy setup, SSL/TLS configuration, load balancing, and security best practices.

## Basic Nginx Configuration

### Main Configuration (nginx.conf)
```nginx
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
    use epoll;
    multi_accept on;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # Logging format
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    # Basic settings
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    server_tokens off;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/xml+rss
        application/atom+xml
        image/svg+xml;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=login:10m rate=1r/s;

    # Include server configurations
    include /etc/nginx/conf.d/*.conf;
}
```

## Server Configurations

### HTTP to HTTPS Redirect
```nginx
# /etc/nginx/conf.d/redirect.conf
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;
    
    # Redirect all HTTP requests to HTTPS
    return 301 https://$server_name$request_uri;
}
```

### Main Server Configuration
```nginx
# /etc/nginx/conf.d/reports-api.conf
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL Configuration
    ssl_certificate /etc/nginx/ssl/your-domain.crt;
    ssl_certificate_key /etc/nginx/ssl/your-domain.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Client max body size for file uploads
    client_max_body_size 100M;

    # Proxy settings
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Forwarded-Host $server_name;
    proxy_set_header X-Forwarded-Port $server_port;

    # Timeouts
    proxy_connect_timeout 60s;
    proxy_send_timeout 60s;
    proxy_read_timeout 60s;

    # Main API location
    location / {
        limit_req zone=api burst=20 nodelay;
        proxy_pass http://reports-api-backend;
        
        # Handle WebSocket upgrades if needed
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # API endpoints with specific rate limiting
    location /api/authEntry/login {
        limit_req zone=login burst=5 nodelay;
        proxy_pass http://reports-api-backend;
    }

    # File upload endpoints
    location ~ ^/api/v1/(problem|resolution|progress)/.*$ {
        client_max_body_size 100M;
        proxy_pass http://reports-api-backend;
        
        # Increase timeouts for file uploads
        proxy_connect_timeout 300s;
        proxy_send_timeout 300s;
        proxy_read_timeout 300s;
    }

    # Static files (if served by Nginx)
    location /static/ {
        alias /var/www/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Health check endpoint
    location /health {
        access_log off;
        proxy_pass http://reports-api-backend;
    }

    # Swagger documentation (development only)
    location /swagger/ {
        # Restrict access in production
        # allow 192.168.1.0/24;
        # deny all;
        proxy_pass http://reports-api-backend;
    }

    # Error pages
    error_page 404 /404.html;
    error_page 500 502 503 504 /50x.html;
    
    location = /50x.html {
        root /usr/share/nginx/html;
    }
}
```

## Upstream Configuration

### Single Backend
```nginx
# /etc/nginx/conf.d/upstream.conf
upstream reports-api-backend {
    server reports-api:5001;
    
    # Health check (nginx plus only)
    # health_check;
    
    # Connection settings
    keepalive 32;
    keepalive_requests 100;
    keepalive_timeout 60s;
}
```

### Load Balancing (Multiple Instances)
```nginx
upstream reports-api-backend {
    # Load balancing methods:
    # - round_robin (default)
    # - least_conn
    # - ip_hash
    # - hash $request_uri consistent
    
    least_conn;
    
    server reports-api-1:5001 weight=3 max_fails=3 fail_timeout=30s;
    server reports-api-2:5001 weight=3 max_fails=3 fail_timeout=30s;
    server reports-api-3:5001 weight=2 max_fails=3 fail_timeout=30s backup;
    
    keepalive 32;
}
```

## SSL/TLS Configuration

### Let's Encrypt Setup
```bash
# Install Certbot
sudo apt-get update
sudo apt-get install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### Self-Signed Certificate (Development)
```bash
# Generate private key
openssl genrsa -out /etc/nginx/ssl/your-domain.key 2048

# Generate certificate
openssl req -new -x509 -key /etc/nginx/ssl/your-domain.key -out /etc/nginx/ssl/your-domain.crt -days 365

# Set permissions
chmod 600 /etc/nginx/ssl/your-domain.key
chmod 644 /etc/nginx/ssl/your-domain.crt
```

### Advanced SSL Configuration
```nginx
server {
    listen 443 ssl http2;
    
    # SSL certificates
    ssl_certificate /etc/nginx/ssl/your-domain.crt;
    ssl_certificate_key /etc/nginx/ssl/your-domain.key;
    
    # SSL protocols and ciphers
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    # SSL session settings
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;
    ssl_session_tickets off;
    
    # OCSP stapling
    ssl_stapling on;
    ssl_stapling_verify on;
    ssl_trusted_certificate /etc/nginx/ssl/chain.crt;
    resolver 8.8.8.8 8.8.4.4 valid=300s;
    resolver_timeout 5s;
    
    # HSTS
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload" always;
}
```

## Security Configuration

### Rate Limiting
```nginx
http {
    # Define rate limit zones
    limit_req_zone $binary_remote_addr zone=general:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=login:10m rate=1r/s;
    limit_req_zone $binary_remote_addr zone=api:10m rate=20r/s;
    limit_req_zone $binary_remote_addr zone=upload:10m rate=2r/s;
    
    # Connection limiting
    limit_conn_zone $binary_remote_addr zone=conn_limit_per_ip:10m;
    limit_conn_zone $server_name zone=conn_limit_per_server:10m;
}

server {
    # Apply rate limits
    location / {
        limit_req zone=general burst=20 nodelay;
        limit_conn conn_limit_per_ip 10;
    }
    
    location /api/authEntry/login {
        limit_req zone=login burst=5 nodelay;
    }
    
    location /api/ {
        limit_req zone=api burst=50 nodelay;
    }
    
    location ~ /upload {
        limit_req zone=upload burst=10 nodelay;
        client_max_body_size 100M;
    }
}
```

### IP Whitelisting/Blacklisting
```nginx
# Create IP lists
# /etc/nginx/conf.d/ip-whitelist.conf
geo $whitelist {
    default 0;
    192.168.1.0/24 1;
    10.0.0.0/8 1;
    172.16.0.0/12 1;
}

# /etc/nginx/conf.d/ip-blacklist.conf
geo $blacklist {
    default 0;
    192.168.100.100 1;
    10.0.0.50 1;
}

server {
    # Block blacklisted IPs
    if ($blacklist) {
        return 403;
    }
    
    # Admin area - whitelist only
    location /admin {
        if ($whitelist = 0) {
            return 403;
        }
        proxy_pass http://reports-api-backend;
    }
}
```

### DDoS Protection
```nginx
# /etc/nginx/conf.d/ddos-protection.conf
limit_req_zone $binary_remote_addr zone=ddos:10m rate=1r/s;

server {
    location / {
        # Basic DDoS protection
        limit_req zone=ddos burst=5 nodelay;
        
        # Block requests with no User-Agent
        if ($http_user_agent = "") {
            return 403;
        }
        
        # Block common attack patterns
        if ($request_uri ~* "(\<|%3C).*script.*(\>|%3E)") {
            return 403;
        }
        
        if ($query_string ~* "[;'\x22\x27\x3C\x3E\x00\x10\x0B\x0D\x0A\x1A]") {
            return 403;
        }
    }
}
```

## Caching Configuration

### Static Content Caching
```nginx
server {
    # Cache static files
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|pdf|txt)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
        add_header Vary Accept-Encoding;
        access_log off;
    }
    
    # Cache API responses (careful with dynamic content)
    location /api/v1/department/listall {
        proxy_pass http://reports-api-backend;
        proxy_cache api_cache;
        proxy_cache_valid 200 5m;
        proxy_cache_key "$scheme$request_method$host$request_uri";
        add_header X-Cache-Status $upstream_cache_status;
    }
}

# Cache configuration
proxy_cache_path /var/cache/nginx/api levels=1:2 keys_zone=api_cache:10m max_size=100m inactive=60m use_temp_path=off;
```

### Microcaching for Dynamic Content
```nginx
server {
    location /api/v1/dashboard/data {
        proxy_pass http://reports-api-backend;
        
        # Microcaching - cache for 1 second
        proxy_cache api_cache;
        proxy_cache_valid 200 1s;
        proxy_cache_lock on;
        proxy_cache_use_stale updating;
    }
}
```

## Monitoring and Logging

### Access Log Configuration
```nginx
http {
    # Custom log format
    log_format detailed '$remote_addr - $remote_user [$time_local] '
                       '"$request" $status $bytes_sent '
                       '"$http_referer" "$http_user_agent" '
                       '$request_time $upstream_response_time '
                       '$upstream_addr $upstream_status';
    
    # Separate logs for different endpoints
    map $request_uri $log_file {
        ~^/api/authEntry/ auth;
        ~^/api/v1/problem/ problem;
        ~^/api/v1/dashboard/ dashboard;
        default main;
    }
    
    access_log /var/log/nginx/access.log detailed;
    access_log /var/log/nginx/auth.log detailed if=$log_file=auth;
    access_log /var/log/nginx/problem.log detailed if=$log_file=problem;
    access_log /var/log/nginx/dashboard.log detailed if=$log_file=dashboard;
}
```

### Error Handling
```nginx
server {
    # Custom error pages
    error_page 400 /error/400.html;
    error_page 401 /error/401.html;
    error_page 403 /error/403.html;
    error_page 404 /error/404.html;
    error_page 500 502 503 504 /error/50x.html;
    
    location ^~ /error/ {
        internal;
        root /var/www/html;
    }
    
    # Log 4xx and 5xx errors
    access_log /var/log/nginx/error_access.log detailed if=$status~^[45];
}
```

### Health Check Endpoint
```nginx
server {
    # Nginx status (for monitoring)
    location /nginx_status {
        stub_status on;
        access_log off;
        allow 127.0.0.1;
        allow 192.168.1.0/24;
        deny all;
    }
    
    # Application health check
    location /health {
        proxy_pass http://reports-api-backend;
        access_log off;
        
        # Return 503 if backend is down
        proxy_intercept_errors on;
        error_page 502 503 504 = @maintenance;
    }
    
    location @maintenance {
        return 503 "Service temporarily unavailable";
        add_header Content-Type text/plain;
    }
}
```

## Docker Integration

### Nginx Dockerfile
```dockerfile
FROM nginx:alpine

# Copy configuration files
COPY nginx.conf /etc/nginx/nginx.conf
COPY conf.d/ /etc/nginx/conf.d/

# Copy SSL certificates
COPY ssl/ /etc/nginx/ssl/

# Copy static files
COPY static/ /var/www/static/

# Create cache directory
RUN mkdir -p /var/cache/nginx/api

# Set permissions
RUN chown -R nginx:nginx /var/cache/nginx

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
```

### Docker Compose with Nginx
```yaml
version: '3.8'

services:
  nginx:
    build:
      context: ./nginx
      dockerfile: Dockerfile
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/conf.d:/etc/nginx/conf.d:ro
      - ./ssl:/etc/nginx/ssl:ro
      - nginx_cache:/var/cache/nginx
    depends_on:
      - reports-api
    networks:
      - reports-network
    restart: unless-stopped

  reports-api:
    image: reports-api:latest
    expose:
      - "5001"
    networks:
      - reports-network

volumes:
  nginx_cache:

networks:
  reports-network:
    driver: bridge
```

## Performance Optimization

### Worker Process Optimization
```nginx
# Automatically detect number of CPU cores
worker_processes auto;

# Increase worker connections
events {
    worker_connections 2048;
    use epoll;
    multi_accept on;
}

# Optimize file operations
http {
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    
    # Increase buffer sizes
    client_body_buffer_size 128k;
    client_max_body_size 100m;
    client_header_buffer_size 1k;
    large_client_header_buffers 4 4k;
    output_buffers 1 32k;
    postpone_output 1460;
}
```

### Connection Optimization
```nginx
upstream reports-api-backend {
    server reports-api:5001;
    
    # Keep connections alive
    keepalive 32;
    keepalive_requests 1000;
    keepalive_timeout 60s;
}

server {
    # Enable HTTP/2
    listen 443 ssl http2;
    
    # Optimize proxy connections
    proxy_http_version 1.1;
    proxy_set_header Connection "";
    
    # Buffer responses
    proxy_buffering on;
    proxy_buffer_size 4k;
    proxy_buffers 8 4k;
    proxy_busy_buffers_size 8k;
}
```

## Troubleshooting

### Common Issues

#### 1. 502 Bad Gateway
```bash
# Check if backend is running
docker ps | grep reports-api

# Check Nginx error logs
docker logs nginx

# Test backend connectivity
docker exec nginx curl http://reports-api:5001/
```

#### 2. SSL Certificate Issues
```bash
# Test SSL certificate
openssl s_client -connect your-domain.com:443

# Check certificate expiry
openssl x509 -in /etc/nginx/ssl/your-domain.crt -text -noout | grep "Not After"

# Verify certificate chain
nginx -t
```

#### 3. Rate Limiting Issues
```bash
# Check rate limit zones
grep "limiting requests" /var/log/nginx/error.log

# Adjust rate limits in configuration
limit_req zone=api burst=100 nodelay;
```

### Debugging Commands
```bash
# Test Nginx configuration
nginx -t

# Reload Nginx configuration
nginx -s reload

# View Nginx processes
ps aux | grep nginx

# Check listening ports
netstat -tlnp | grep nginx

# Monitor real-time logs
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log
```

This Nginx setup guide provides comprehensive configuration for production deployment of the Reports API with security, performance, and monitoring best practices.