# Reports API

## Overview
Reports API is a backend service for managing problems, tasks, IP phones, programs, departments, and branches with authentication system for users and administrators.

## Technology Stack

### Backend Framework
- **Go 1.23.0** - Programming language
- **Fiber v2** - Web framework for Go
- **MySQL** - Database

### Libraries & Dependencies
- **github.com/go-sql-driver/mysql** - MySQL driver
- **github.com/gofiber/fiber/v2** - Web framework
- **github.com/joho/godotenv** - Environment variables management
- **github.com/go-telegram-bot-api/telegram-bot-api/v5** - Telegram bot integration
- **github.com/minio/minio-go/v7** - Object storage client

### Features
- ✅ Problem Management
- ✅ Task Management
- ✅ IP Phone Management
- ✅ Program Management
- ✅ Department Management
- ✅ Branch Management
- ✅ User Authentication (login/logout)
- ✅ File Upload Support
- ✅ Telegram Integration
- ✅ CORS Support
- ✅ Environment-based Configuration

## Project Structure
```
reports-api/
├── db/                 # Database connection
├── handlers/           # HTTP handlers
├── middleware/         # Custom middleware
├── models/            # Data models
├── routes/            # Route definitions
├── utils/             # Utility functions
├── .env.dev           # Development environment
├── .env.prod          # Production environment
├── Dockerfile         # Docker configuration
└── main.go           # Application entry point
```

## Quick Start

### Prerequisites
- Go 1.23.0 or later
- MySQL Database
- Git

### Installation
```bash
# Clone repository
git clone <repository-url>
cd reports-api

# Install dependencies
go mod tidy

# Copy environment file
cp .env.dev .env

# Edit environment variables as needed
# Configure .env file with appropriate values

# Run application
go run main.go
```

### Environment Options
```bash
# Development
go run main.go dev
# or
go run main.go -d

# Production
go run main.go prod
# or
go run main.go -p
```

## Documentation
For detailed API usage and installation instructions, please refer to:
- [API Usage Guide](./docs/API_USAGE.md) - How to use the API
- [Installation Guide](./docs/INSTALLATION.md) - Installation instructions
- [Deployment Guide](./docs/DEPLOYMENT.md) - Deployment guide

## API Documentation
Swagger UI is available when the server is running at:
```
http://localhost:5001/docs
```

## Health Check
When the server is running, you can check its status at:
```
GET http://localhost:5001/
```

## License
MIT License
