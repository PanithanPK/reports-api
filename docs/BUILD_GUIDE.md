# 🐳 Docker Build Guide

## Quick Build Commands

### Development Build
```bash
docker build -t <your-dockerhub-username>/reports-api:dev -f Dockerfile.dev .
docker push <your-dockerhub-username>/reports-api:dev
```

### Production Build
```bash
docker build -t <your-dockerhub-username>/reports-api:prod -f Dockerfile .
docker push <your-dockerhub-username>/reports-api:prod
```

## Interactive Build Script

Run the automated build script:
```bash
./build2.sh
```

### Build Options:
1. **Development Environment** - Uses `Dockerfile.dev`
2. **Production Environment** - Uses `Dockerfile`  
3. **Build Both** - Builds both environments

### Script Features:
- ✅ Interactive menu selection
- ✅ Confirmation prompts for safety
- ✅ Automatic Docker Hub login
- ✅ Progress indicators
- ✅ Error handling with rollback
- ✅ Automatic cleanup

## GitLab CI/CD Pipeline

### Pipeline Triggers:
- **Push to `develop`** → Auto-build dev image
- **Push to `main`** → Manual trigger for prod image

### Required Variables:
Add these in GitLab **Settings** → **CI/CD** → **Variables**:

| Variable | Value | Settings |
|----------|-------|----------|
| `DOCKER_TOKEN_DEV` | `<your-dev-docker-token>` | Masked ✅ |
| `DOCKER_TOKEN_PROD` | `<your-prod-docker-token>` | Masked ✅ |

**Important:** Uncheck "Protected" for development builds!

## Docker Images

### Available Tags:
- `<your-dockerhub-username>/reports-api:dev` - Development version (port 5000)
- `<your-dockerhub-username>/reports-api:prod` - Production version (port 5001)

### Running Containers:

**Development:**
```bash
docker run -d \
  --name reports-api-dev \
  -p 5001:5000 \
  -e DB_HOST=localhost \
  -e DB_USER=reports_user \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=reports \
  <your-dockerhub-username>/reports-api:dev
```

**Production:**
```bash
docker run -d \
  --name reports-api-prod \
  -p 5001:5001 \
  -e DB_HOST=localhost \
  -e DB_USER=reports_user \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=reports \
  <your-dockerhub-username>/reports-api:prod
```

## Environment Variables

### Required:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=reports_user
DB_PASSWORD=your_password
DB_NAME=reports
PORT=5001
JWT_SECRET=your_jwt_secret
TOKEN_HOUR_LIFESPAN=24
```

### Optional:
```env
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=reports-files
```

## Troubleshooting

### Common Issues:

**1. Docker login fails:**
```bash
# Check if tokens are set correctly
echo $DOCKER_TOKEN_DEV
# Should show masked value in GitLab CI
```

**2. Build fails - .env not found:**
- ✅ Fixed: Dockerfiles no longer copy .env files
- ✅ Use environment variables instead

**3. GitLab CI variables not working:**
- ✅ Uncheck "Protected" in variable settings
- ✅ Make sure variable names match exactly

**4. Port conflicts:**
- Dev image: Internal port 5000, map to any external port
- Prod image: Internal port 5001, map to any external port

## Quick Start

1. **Local Development:**
   ```bash
   go run main.go dev
   ```

2. **Docker Development:**
   ```bash
   docker build -t reports-api-dev -f Dockerfile.dev .
   docker run -p 5001:5000 reports-api-dev
   ```

3. **Health Check:**
   ```bash
   curl http://localhost:5001/
   ```

## Notes

- 🔒 Never commit .env files to repository
- 🐳 Use environment variables in containers
- 🚀 GitLab CI/CD handles automatic builds
- 📝 Manual script provides interactive experience
- ✅ Health check endpoint: `GET /`