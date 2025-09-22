# Docker Setup Guide

## Overview
This guide covers Docker setup for the Reports API, including development and production configurations, Docker Compose setup, and best practices.

## Docker Images

### Development Image (Dockerfile.dev)
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env.dev .
EXPOSE 5000
CMD ["./main", "dev"]
```

### Production Image (Dockerfile)
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env.prod .
EXPOSE 5001
CMD ["./main", "prod"]
```

## Docker Compose Configuration

### Basic Docker Compose (docker-compose.yml)
```yaml
version: '3.8'

services:
  reports-api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "5001:5001"
    environment:
      - PORT=5001
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=reports_user
      - DB_PASS=reports_password
      - DB_NAME=reports_db
    depends_on:
      - mysql
      - minio
    networks:
      - reports-network
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=root_password
      - MYSQL_DATABASE=reports_db
      - MYSQL_USER=reports_user
      - MYSQL_PASSWORD=reports_password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - reports-network
    restart: unless-stopped

  minio:
    image: minio/minio:latest
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      - MINIO_ROOT_USER=minioadmin
      - MINIO_ROOT_PASSWORD=minioadmin123
    volumes:
      - minio_data:/data
    command: server /data --console-address ":9001"
    networks:
      - reports-network
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
    networks:
      - reports-network
    restart: unless-stopped

volumes:
  mysql_data:
  minio_data:

networks:
  reports-network:
    driver: bridge
```

### Development Docker Compose (docker-compose.dev.yml)
```yaml
version: '3.8'

services:
  reports-api-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "5000:5000"
    environment:
      - PORT=5000
      - DB_HOST=mysql-dev
      - DB_PORT=3306
      - DB_USER=reports_user
      - DB_PASS=reports_password
      - DB_NAME=reports_db_dev
    volumes:
      - .:/app
      - /app/vendor
    depends_on:
      - mysql-dev
      - minio-dev
    networks:
      - reports-dev-network
    restart: unless-stopped

  mysql-dev:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=root_password
      - MYSQL_DATABASE=reports_db_dev
      - MYSQL_USER=reports_user
      - MYSQL_PASSWORD=reports_password
    ports:
      - "3307:3306"
    volumes:
      - mysql_dev_data:/var/lib/mysql
    networks:
      - reports-dev-network
    restart: unless-stopped

  minio-dev:
    image: minio/minio:latest
    ports:
      - "9002:9000"
      - "9003:9001"
    environment:
      - MINIO_ROOT_USER=minioadmin
      - MINIO_ROOT_PASSWORD=minioadmin123
    volumes:
      - minio_dev_data:/data
    command: server /data --console-address ":9001"
    networks:
      - reports-dev-network
    restart: unless-stopped

volumes:
  mysql_dev_data:
  minio_dev_data:

networks:
  reports-dev-network:
    driver: bridge
```

## Build Scripts

### Interactive Build Script (build2.sh)
```bash
#!/bin/bash

echo "=== Reports API Docker Build Script ==="
echo "Select build environment:"
echo "1) Development"
echo "2) Production"
echo "3) Both"
read -p "Enter choice (1-3): " choice

DOCKER_USERNAME=""
read -p "Enter Docker Hub username: " DOCKER_USERNAME

case $choice in
    1)
        echo "Building development image..."
        docker build -t $DOCKER_USERNAME/reports-api:dev -f Dockerfile.dev .
        
        read -p "Push to Docker Hub? (y/n): " push_choice
        if [ "$push_choice" = "y" ]; then
            docker login
            docker push $DOCKER_USERNAME/reports-api:dev
        fi
        ;;
    2)
        echo "Building production image..."
        docker build -t $DOCKER_USERNAME/reports-api:latest .
        docker build -t $DOCKER_USERNAME/reports-api:prod .
        
        read -p "Push to Docker Hub? (y/n): " push_choice
        if [ "$push_choice" = "y" ]; then
            docker login
            docker push $DOCKER_USERNAME/reports-api:latest
            docker push $DOCKER_USERNAME/reports-api:prod
        fi
        ;;
    3)
        echo "Building both images..."
        docker build -t $DOCKER_USERNAME/reports-api:dev -f Dockerfile.dev .
        docker build -t $DOCKER_USERNAME/reports-api:latest .
        docker build -t $DOCKER_USERNAME/reports-api:prod .
        
        read -p "Push to Docker Hub? (y/n): " push_choice
        if [ "$push_choice" = "y" ]; then
            docker login
            docker push $DOCKER_USERNAME/reports-api:dev
            docker push $DOCKER_USERNAME/reports-api:latest
            docker push $DOCKER_USERNAME/reports-api:prod
        fi
        ;;
    *)
        echo "Invalid choice"
        exit 1
        ;;
esac

echo "Build completed!"
```

## Docker Commands

### Basic Commands
```bash
# Build development image
docker build -t reports-api:dev -f Dockerfile.dev .

# Build production image
docker build -t reports-api:prod .

# Run development container
docker run -d --name reports-api-dev -p 5000:5000 reports-api:dev

# Run production container
docker run -d --name reports-api-prod -p 5001:5001 reports-api:prod

# View logs
docker logs -f reports-api-dev

# Stop container
docker stop reports-api-dev

# Remove container
docker rm reports-api-dev

# Remove image
docker rmi reports-api:dev
```

### Docker Compose Commands
```bash
# Start all services
docker-compose up -d

# Start development services
docker-compose -f docker-compose.dev.yml up -d

# View logs
docker-compose logs -f reports-api

# Stop all services
docker-compose down

# Rebuild and start
docker-compose up -d --build

# Remove volumes (careful!)
docker-compose down -v
```

## Environment Variables

### Required Environment Variables
```bash
# Server Configuration
PORT=5001

# Database Configuration
DB_HOST=mysql
DB_PORT=3306
DB_USER=reports_user
DB_PASS=reports_password
DB_NAME=reports_db

# MinIO Configuration
End_POINT=minio:9000
ACCESS_KEY=minioadmin
SECRET_ACCESSKEY=minioadmin123
BUCKET_NAME=reports-bucket

# Telegram Configuration
BOT_TOKEN=your_telegram_bot_token
CHAT_ID=your_telegram_chat_id

# JWT Configuration
JWT_SECRET=your_jwt_secret_key
TOKEN_HOUR_LIFESPAN=24
```

### Docker Environment File (.env)
```bash
# Copy environment template
cp .env.example .env

# Edit environment variables
nano .env
```

## Health Checks

### Application Health Check
```bash
# Check if API is running
curl http://localhost:5001/

# Expected response
{"status":"OK","version":"1.0.0"}
```

### Docker Health Check
```dockerfile
# Add to Dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:5001/ || exit 1
```

## Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Find process using port
lsof -i :5001

# Kill process
kill -9 <PID>

# Or use different port
docker run -p 5002:5001 reports-api:prod
```

#### 2. Database Connection Error
```bash
# Check MySQL container
docker logs mysql

# Check network connectivity
docker exec reports-api ping mysql

# Verify environment variables
docker exec reports-api env | grep DB_
```

#### 3. MinIO Connection Error
```bash
# Check MinIO container
docker logs minio

# Access MinIO console
http://localhost:9001

# Verify MinIO credentials
docker exec reports-api env | grep MINIO
```

#### 4. Build Failures
```bash
# Clear Docker cache
docker system prune -a

# Rebuild without cache
docker build --no-cache -t reports-api:dev .

# Check Dockerfile syntax
docker build --dry-run -t reports-api:dev .
```

### Debugging

#### Container Debugging
```bash
# Access container shell
docker exec -it reports-api sh

# Check container resources
docker stats reports-api

# Inspect container
docker inspect reports-api

# Check container processes
docker exec reports-api ps aux
```

#### Log Analysis
```bash
# View real-time logs
docker logs -f --tail 100 reports-api

# Export logs to file
docker logs reports-api > app.log 2>&1

# Filter logs by level
docker logs reports-api 2>&1 | grep ERROR
```

## Performance Optimization

### Resource Limits
```yaml
# In docker-compose.yml
services:
  reports-api:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 512M
        reservations:
          cpus: '1.0'
          memory: 256M
```

### Multi-stage Build Optimization
```dockerfile
# Use specific Go version
FROM golang:1.23-alpine AS builder

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main .
```

## Security Best Practices

### Image Security
```dockerfile
# Use non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
USER appuser

# Update packages
RUN apk update && apk upgrade && apk add --no-cache ca-certificates

# Remove unnecessary packages
RUN apk del build-dependencies
```

### Network Security
```yaml
# Use custom networks
networks:
  reports-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### Secrets Management
```yaml
# Use Docker secrets
secrets:
  db_password:
    file: ./secrets/db_password.txt
  jwt_secret:
    file: ./secrets/jwt_secret.txt

services:
  reports-api:
    secrets:
      - db_password
      - jwt_secret
```

## Monitoring

### Container Monitoring
```bash
# Monitor resource usage
docker stats

# Monitor specific container
docker stats reports-api

# Export metrics
docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"
```

### Log Monitoring
```yaml
# Add logging driver
services:
  reports-api:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## Backup and Recovery

### Database Backup
```bash
# Backup MySQL data
docker exec mysql mysqldump -u root -p reports_db > backup.sql

# Restore MySQL data
docker exec -i mysql mysql -u root -p reports_db < backup.sql
```

### Volume Backup
```bash
# Backup volumes
docker run --rm -v mysql_data:/data -v $(pwd):/backup alpine tar czf /backup/mysql_backup.tar.gz /data

# Restore volumes
docker run --rm -v mysql_data:/data -v $(pwd):/backup alpine tar xzf /backup/mysql_backup.tar.gz -C /
```

This Docker setup guide provides comprehensive instructions for containerizing and deploying the Reports API with all necessary services and configurations.