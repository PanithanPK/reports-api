# Reports API - Project Overview

## Project Overview
Reports API is a comprehensive Backend API system developed with Go (Golang) and Fiber Framework for managing organizational Problem Reporting System. It features MySQL database integration, Telegram notification system, file storage with MinIO, and comprehensive task management with resolution tracking.

## Technologies Used
- **Backend**: Go 1.23 + Fiber Framework v2
- **Database**: MySQL with Connection Pooling
- **Authentication**: Session-based Authentication
- **File Storage**: MinIO Object Storage
- **Notification**: Telegram Bot API
- **Documentation**: Swagger/OpenAPI
- **Deployment**: Docker + Docker Compose
- **CI/CD**: GitLab CI/CD Pipeline

## Project Structure

### 1. Main Application (`main.go`)
- Supports dev/prod/default environments
- Memory and CPU optimization management
- CORS configuration
- Static file serving
- Custom logging system

### 2. Database Layer (`db/`)
- **db.go**: MySQL connection with Connection Pool
- Supports Environment Variables for Database Config
- Connection pooling optimization

### 3. Models (`models/`)
- **reportProblemModel.go**: Task/Problem related models with file upload support
- **resolutionModels.go**: Resolution tracking models
- **progressModels.go**: Progress tracking models
- **userModels.go**: User authentication models  
- **departmentModels.go**: Department management models
- **branchModels.go**: Branch management models
- **listphoneModel.go**: IP Phone management models
- **programModels.go**: System/Program models
- **scoresModel.go**: Department scoring models
- **dashboardModel.go**: Dashboard data models
- **datasummaryModel.go**: Data summary and export models
- **telegramModels.go**: Telegram integration models
- **assignto.go**: Assignment tracking models
- **paginationModels.go**: Pagination utilities

### 4. Handlers (`handlers/`)
- **reportProblem.go**: Task/Problem CRUD operations with file upload
- **resolution.go**: Resolution management and tracking
- **progress.go**: Progress tracking and updates
- **user.go**: Authentication & User management
- **department.go**: Department management
- **branche.go**: Branch management
- **listphones.go**: IP Phone management
- **program.go**: System/Program management
- **scores.go**: Department scoring system
- **dashboard.go**: Dashboard data aggregation
- **datasummary.go**: Data export and summary
- **common/**: Common services
  - **fundamentalService.go**: Core utility functions
  - **minioService.go**: File storage operations
  - **telegramService.go**: Telegram integration

### 5. Routes (`routes/`)
- **routes.go**: API routes registration
- Separate Authentication routes and Business logic routes

### 6. Utils (`utils/`)
- **pagination.go**: Pagination utilities
- **datetime.go**: Date/time utility functions

### 7. Configuration (`config/`)
- **config.go**: Application configuration management

### 8. Middleware (`middleware/`)
- **middleware.go**: Rate limiting and security middleware
- **SECURITY_RECOMMENDATIONS.md**: Security guidelines

### 9. Documentation (`docs/`)
- **swagger.json/yaml**: API documentation
- **API_USAGE.md**: API usage guide
- **DEPLOYMENT.md**: Deployment instructions
- **INSTALLATION.md**: Installation guide

## Main Features

### 1. Task/Problem Management
- Create/Edit/Delete Tasks with file attachments
- Auto-generate Ticket Number system (TK-DDMMYYYY-XXXX)
- Full-text search and advanced filtering
- Pagination support with query parameters
- Status tracking (0=Pending, 1=In Progress, 2=Resolved)
- Assignment tracking with notification
- Overdue calculation and monitoring
- Issue type categorization

### 2. Resolution Management
- Create and update task resolutions
- File attachment support for solutions
- Resolution tracking with timestamps
- Integration with task status updates

### 3. Progress Tracking
- Add progress updates to tasks
- File attachment support for progress entries
- Timeline tracking of task progress
- Update and delete progress entries

### 4. File Storage & Management
- MinIO object storage integration
- Image upload support (JPEG, PNG)
- Automatic file naming and organization
- File deletion capabilities
- Secure file access URLs

### 5. IP Phone Management
- Comprehensive IP Phone database
- Department linking and organization
- Advanced search and filtering
- CSV export functionality

### 6. Department & Branch Management
- Hierarchical organization structure
- Department scoring system
- Performance tracking and analytics
- CSV export capabilities

### 7. User Authentication & Authorization
- Session-based authentication
- Role-based access control (user/admin)
- User registration/login/logout
- Password management with bcrypt

### 8. Telegram Integration
- Real-time notifications for new tasks
- Status update notifications
- Assignment change notifications
- Image sharing in notifications
- Message update/delete synchronization

### 9. Dashboard & Analytics
- Comprehensive dashboard data
- Department performance metrics
- Task statistics and trends
- Data export functionality (CSV)

### 10. Data Export & Summary
- CSV export for all major entities
- Data summary and analytics
- Comprehensive reporting capabilities

### 11. Scoring System
- Automated department scoring
- Monthly performance calculation
- Score deduction based on task volume
- Performance tracking over time

## API Endpoints

### Problem/Task Routes
```
GET    /api/v1/problem/list                    - List tasks with pagination
GET    /api/v1/problem/list/:query             - List tasks with search query
GET    /api/v1/problem/list/:column/:query     - List tasks with column-specific search
GET    /api/v1/problem/list/sort/:column/:query - List tasks with sorting
POST   /api/v1/problem/create                  - Create new task with file upload
GET    /api/v1/problem/:id                     - Get task detail
PUT    /api/v1/problem/update/:id              - Update task
DELETE /api/v1/problem/delete/:id             - Delete task
PUT    /api/v1/problem/update/assignto/:id     - Update task assignment
```

### Resolution Routes
```
GET    /api/v1/resolution/:id          - Get task resolution
POST   /api/v1/resolution/create/:id   - Create resolution with file upload
PUT    /api/v1/resolution/update/:id   - Update resolution
DELETE /api/v1/resolution/delete/:id  - Delete resolution
```

### Progress Routes
```
GET    /api/v1/progress/:id                - Get task progress
POST   /api/v1/progress/create/:id         - Create progress entry
PUT    /api/v1/progress/update/:id/:pgid   - Update progress entry
DELETE /api/v1/progress/delete/:id/:pgid  - Delete progress entry
```

### IP Phone Routes
```
GET    /api/v1/ipphone/list           - List IP phones with pagination
GET    /api/v1/ipphone/list/:query    - List IP phones with search
GET    /api/v1/ipphone/listall        - List all IP phones
GET    /api/v1/ipphone/:id            - Get IP phone detail
POST   /api/v1/ipphone/create         - Create IP phone
PUT    /api/v1/ipphone/update/:id     - Update IP phone
DELETE /api/v1/ipphone/delete/:id     - Delete IP phone
```

### Department Routes
```
GET    /api/v1/department/list           - List departments with pagination
GET    /api/v1/department/list/:query    - List departments with search
GET    /api/v1/department/listall        - List all departments
POST   /api/v1/department/create         - Create department
GET    /api/v1/department/:id            - Get department detail
PUT    /api/v1/department/update/:id     - Update department
DELETE /api/v1/department/delete/:id     - Delete department
```

### Branch Routes
```
GET    /api/v1/branch/list           - List branches with pagination
GET    /api/v1/branch/list/:query    - List branches with search
POST   /api/v1/branch/create         - Create branch
GET    /api/v1/branch/:id            - Get branch detail
PUT    /api/v1/branch/update/:id     - Update branch
DELETE /api/v1/branch/delete/:id     - Delete branch
```

### Program Routes
```
GET    /api/v1/program/list              - List programs with pagination
GET    /api/v1/program/list/:query       - List programs with search
POST   /api/v1/program/create            - Create program
GET    /api/v1/program/:id               - Get program detail
PUT    /api/v1/program/update/:id        - Update program
DELETE /api/v1/program/delete/:id       - Delete program
GET    /api/v1/program/type/list         - List program types
GET    /api/v1/program/type/list/:query  - List program types with search
POST   /api/v1/program/type/create       - Create program type
POST   /api/v1/program/type/update/:id   - Update program type
DELETE /api/v1/program/type/delete/:id  - Delete program type
```

### Authentication Routes
```
POST   /api/authEntry/login           - User login
POST   /api/authEntry/registerUser    - Register user
POST   /api/authEntry/registerAdmin   - Register admin
PUT    /api/authEntry/updateUser      - Update user
DELETE /api/authEntry/deleteUser      - Delete user
POST   /api/authEntry/logout          - User logout
```

### Dashboard & Analytics Routes
```
GET    /api/v1/dashboard/data              - Get dashboard data
GET    /api/v1/dashboard/data/csv          - Export tasks to CSV
GET    /api/v1/dashboard/data/phonecsv     - Export IP phones to CSV
GET    /api/v1/dashboard/data/departmetcsv - Export departments to CSV
GET    /api/v1/dashboard/data/branchcsv    - Export branches to CSV
GET    /api/v1/dashboard/data/systemcsv    - Export systems to CSV
```

### Scoring Routes
```
GET    /api/v1/scores/list            - List department scores
GET    /api/v1/scores/:id             - Get score detail
PUT    /api/v1/scores/update/:id      - Update score
DELETE /api/v1/scores/delete/:id     - Delete score
```

### Responsibility Routes
```
GET    /api/v1/respons/list           - List responsibilities
GET    /api/v1/respons/:id            - Get responsibility detail
POST   /api/v1/respons/create         - Create responsibility
PUT    /api/v1/respons/update/:id     - Update responsibility
DELETE /api/v1/respons/delete/:id    - Delete responsibility
```

## Database Schema (Main)

### tasks table
- id, ticket_no, phone_id, system_id, department_id
- issue_type, issue_else, text, status
- assignedto_id, assign_to, reported_by
- file_paths (JSON), message_id
- resolved_at, created_at, updated_at
- created_by, updated_by

### resolutions table
- id, task_id, solution, file_paths (JSON)
- telegram_id, telegram_user, message_id
- url, assignto, previous_assignto
- assignedto_id, ticket_no
- resolved_at, created_at

### progress table
- id, task_id, text, file_paths (JSON)
- created_at, updated_at

### ip_phones table  
- id, number, name, department_id
- created_at, updated_at, deleted_at

### departments table
- id, name, branch_id
- created_at, updated_at, deleted_at

### branches table
- id, name
- created_at, updated_at, deleted_at

### systems_program table
- id, name, type
- created_at, updated_at, deleted_at

### issue_types table
- id, name
- created_at, updated_at, deleted_at

### responsibilities table
- id, name
- created_at, updated_at, deleted_at

### users table
- id, username, password, role
- created_at, updated_at, deleted_at

### scores table
- id, department_id, year, month, score

### telegram_chat table
- id, chat_id, username
- created_at, updated_at

## System Workflow

### 1. Task Creation Flow
1. Receive Task data and files from Frontend
2. Auto-generate Ticket Number (TK-DDMMYYYY-XXXX format)
3. Upload files to MinIO storage
4. Save task data to Database with file paths
5. Update Department Score
6. Send Notification via Telegram with images
7. Return task details with file URLs

### 2. Resolution Management Flow
1. Create resolution for completed tasks
2. Upload solution files to MinIO
3. Update task status to resolved
4. Send Telegram notification with solution
5. Track resolution timestamp

### 3. Progress Tracking Flow
1. Add progress updates during task lifecycle
2. Upload progress files to MinIO
3. Maintain progress timeline
4. Allow updates and deletions

### 4. File Management Flow
1. Validate file types (JPEG, PNG)
2. Generate unique object names
3. Upload to MinIO with proper metadata
4. Return secure access URLs
5. Handle file deletion when needed

### 5. Telegram Integration
- Send formatted messages for new tasks
- Include image attachments in notifications
- Update messages when tasks are modified
- Notify on assignment changes
- Sync message updates/deletions
- Store message_id for tracking

### 6. Assignment Management
- Track task assignments with history
- Send notifications on assignment changes
- Update Telegram messages with new assignee
- Maintain assignment audit trail

### 7. Scoring System
- Calculate monthly score per Department
- Start with base score of 100 points
- Deduct points based on task volume
- Track performance over time

### 8. Search and Filtering
- Support multiple search criteria
- Column-specific filtering
- Pagination with configurable limits
- Sorting capabilities

### 9. Data Export System
- Generate CSV exports for all entities
- Real-time data export
- Comprehensive reporting capabilities

## Environment Configuration

### Development (.env.dev)
```
PORT=5001
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=password
DB_NAME=report_db
End_POINT=minio.sys9.co
ACCESS_KEY=minio_access_key
SECRET_ACCESSKEY=minio_secret_key
BUCKET_NAME=reports-bucket
BOT_TOKEN=telegram_bot_token
CHAT_ID=telegram_chat_id
```

### Production (.env.prod)
```
PORT=5000
DB_HOST=production_host
DB_PORT=3306
DB_USER=prod_user
DB_PASS=prod_password
DB_NAME=report_db_prod
End_POINT=minio.sys9.co
ACCESS_KEY=prod_minio_access_key
SECRET_ACCESSKEY=prod_minio_secret_key
BUCKET_NAME=reports-bucket-prod
BOT_TOKEN=prod_telegram_bot_token
CHAT_ID=prod_telegram_chat_id
```

## Deployment

### Docker Deployment
```bash
# Build development image
docker build -t reports-api:dev -f Dockerfile.dev .

# Build production image
docker build -t reports-api:prod .

# Run with Docker Compose (Full Stack)
docker-compose up -d

# Run with development configuration
docker-compose -f docker-compose.dev.yml up -d

# Run container manually
docker run -p 5000:5000 --env-file .env.prod reports-api:prod
```

### Production Deployment with Nginx
```bash
# Full production stack with Nginx, MySQL, MinIO
docker-compose -f docker-compose.prod.yml up -d

# Scale API instances for load balancing
docker-compose up -d --scale reports-api=3

# Update services without downtime
docker-compose up -d --no-deps reports-api
```

### Kubernetes Deployment
```yaml
# Basic Kubernetes deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: reports-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: reports-api
  template:
    metadata:
      labels:
        app: reports-api
    spec:
      containers:
      - name: reports-api
        image: reports-api:latest
        ports:
        - containerPort: 5001
        env:
        - name: DB_HOST
          value: "mysql-service"
        - name: PORT
          value: "5001"
```

### Manual Run
```bash
# Development mode
go run main.go dev
# or
go run main.go -d

# Production mode
go run main.go prod
# or
go run main.go -p

# Default mode
go run main.go
```

### Build Scripts
```bash
# Interactive build script
./build2.sh

# Alternative build scripts
./build.sh
./build3.sh
```

## Security Features
- CORS configuration with specific origins
- Rate limiting middleware
- Input validation and sanitization
- SQL injection prevention
- Session-based authentication with secure cookies
- Role-based access control (user/admin)
- Password hashing with bcrypt
- Soft delete implementation
- File type validation for uploads
- Secure file storage with MinIO
- Environment-based configuration

## Performance Optimization
- Database connection pooling with MySQL
- Memory limit (384MB) with runtime controls
- CPU limit (2 cores) with GOMAXPROCS
- Garbage collection optimization (50% threshold)
- Query optimization with proper indexing
- Efficient pagination implementation
- File streaming for large uploads
- Concurrent file processing
- Response compression
- Static file serving optimization

## API Documentation
- Comprehensive Swagger/OpenAPI documentation
- Interactive API testing interface
- Request/response examples
- Authentication documentation
- Error code documentation
- Available at `/api/v1/swagger/index.html` (dev mode only)

## Usage with ChatGPT
When you need ChatGPT to help develop or modify code, reference this file to:
1. Understand project structure
2. Know existing API endpoints
3. Understand Database schema
4. Know Business logic and workflow
5. Understand Technology stack used

## Version History
- **v1.14.0**: Config system, MinIO integration, data summary, resolution management
- **v1.13.0**: Resolution update features, Telegram formatting improvements
- **v1.12.0**: Resolution creation system
- **v1.11.0**: File upload system, issue types, assignment notifications
- **v1.10.0**: Advanced search, image uploads, progress tracking
- **v1.9.0**: Telegram image sharing
- **v1.8.0**: Assignment system, file uploads, system types
- **v1.7.0**: Telegram updates and assignment features
- **v1.6.0**: Assignment functionality
- **v1.5.0**: Ticket number system
- **v1.4.0**: Unit testing framework
- **v1.3.0**: Pagination system
- **v1.2.0**: Telegram integration
- **v1.1.0**: Initial release

## Infrastructure Architecture

### Container Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Reports website │───>│     Nginx       │───>│   Reports API   │───>│     MySQL       │
│   ( next js )   │    │  (Reverse Proxy)│    │   (Go + Fiber)  │    │   (Database)    │
│                 │<───│   Load Balancer │<───│                 │<───│                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                       │
                              │                       │───────────────────────┐
                              ▼                       ▼                       ▼
                      ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
                      │   Static Files  │    │     MinIO       │    │   Telegram Bot  │
                      │   (CSS/JS/IMG)  │    │ (File Storage)  │    │ (Notifications) │
                      └─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Network Architecture
- **Frontend Network**: Public-facing Nginx (ports 80/443)
- **Backend Network**: Internal API communication
- **Database Network**: Isolated database access
- **Storage Network**: MinIO file storage access

### Scaling Strategy
- **Horizontal Scaling**: Multiple API instances behind load balancer
- **Database Scaling**: Read replicas and connection pooling
- **File Storage**: Distributed MinIO cluster
- **Caching**: Nginx caching and Redis integration

### Security Layers
1. **Network Security**: Firewall and VPN access
2. **Application Security**: JWT authentication and RBAC
3. **Transport Security**: SSL/TLS encryption
4. **Data Security**: Database encryption and secure file storage

## Monitoring and Observability

### Health Monitoring
- **Application Health**: `/health` endpoint monitoring
- **Database Health**: Connection pool monitoring
- **File Storage Health**: MinIO cluster status
- **Nginx Health**: Upstream server monitoring

### Logging Strategy
- **Application Logs**: Structured JSON logging
- **Access Logs**: Nginx access and error logs
- **Database Logs**: MySQL slow query and error logs
- **System Logs**: Container and host system logs

### Metrics Collection
- **Performance Metrics**: Response times and throughput
- **Resource Metrics**: CPU, memory, and disk usage
- **Business Metrics**: Task creation rates and resolution times
- **Error Metrics**: Error rates and failure patterns

## Backup and Recovery

### Database Backup
```bash
# Automated daily backups
docker exec mysql mysqldump -u root -p reports_db > backup_$(date +%Y%m%d).sql

# Point-in-time recovery setup
mysql> SET GLOBAL binlog_format = 'ROW';
mysql> SET GLOBAL log_bin = ON;
```

### File Storage Backup
```bash
# MinIO data backup
mc mirror minio/reports-bucket /backup/minio/reports-bucket

# Automated backup script
#!/bin/bash
DATE=$(date +%Y%m%d)
mc mirror minio/reports-bucket /backup/minio/$DATE/
```

### Disaster Recovery
- **RTO (Recovery Time Objective)**: 4 hours
- **RPO (Recovery Point Objective)**: 1 hour
- **Backup Retention**: 30 days daily, 12 months monthly
- **Geographic Redundancy**: Multi-region backup storage

## Current Status
- **Latest Version**: v1.15.0
- **Environment**: Production-ready with full infrastructure
- **Database**: MySQL 8.0 with connection pooling and backups
- **File Storage**: MinIO cluster with redundancy
- **Notifications**: Telegram bot with image support
- **Documentation**: Complete API and infrastructure documentation
- **CI/CD**: GitLab pipeline with automated testing and deployment
- **Monitoring**: Health checks and logging infrastructure
- **Security**: SSL/TLS, authentication, and security headers
- **Scalability**: Load balancing and horizontal scaling ready