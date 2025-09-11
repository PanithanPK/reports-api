# Reports API - Project Overview

## Project Overview
Reports API is a Backend API system developed with Go (Golang) and Fiber Framework for managing organizational Problem Reporting System. It features MySQL database integration and Telegram notification system.

## Technologies Used
- **Backend**: Go 1.22 + Fiber Framework v2
- **Database**: MySQL
- **Authentication**: Session-based + JWT Token
- **Notification**: Telegram Bot API
- **Deployment**: Docker + Docker Compose

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
- **reportProblemModel.go**: Task/Problem related models
- **userModels.go**: User authentication models  
- **departmentModels.go**: Department management models
- **branchModels.go**: Branch management models
- **listphoneModel.go**: IP Phone management models
- **programModels.go**: System/Program models
- **scoresModel.go**: Department scoring models
- **dashboardModel.go**: Dashboard data models
- **paginationModels.go**: Pagination utilities

### 4. Handlers (`handlers/`)
- **reportProblem.go**: Task/Problem CRUD operations
- **user.go**: Authentication & User management
- **department.go**: Department management
- **branche.go**: Branch management
- **listphones.go**: IP Phone management
- **program.go**: System/Program management
- **scores.go**: Department scoring system
- **dashboard.go**: Dashboard data aggregation
- **telegram.go**: Telegram notification integration

### 5. Routes (`routes/`)
- **routes.go**: API routes registration
- Separate Authentication routes and Business logic routes

### 6. Utils (`utils/`)
- **pagination.go**: Pagination utilities
- **PAGINATION_GUIDE.md**: Pagination documentation

### 7. Tests (`tests/`)
- Unit tests และ Integration tests
- Database testing
- Handler testing
- Pagination testing

## Main Features

### 1. Task/Problem Management
- Create/Edit/Delete Tasks
- Auto-generate Ticket Number system
- Full-text search for Tasks
- Pagination support
- Status tracking (0=Open, 1=In Progress, 2=Closed)
- Overdue calculation

### 2. IP Phone Management
- Manage IP Phone data
- Link with Department
- Search and Filter

### 3. Department & Branch Management
- Manage Department and Branch
- Hierarchical structure (Branch -> Department)
- Department scoring system

### 4. User Authentication
- Session-based authentication
- Role-based access (user/admin)
- User registration/login/logout
- Password management

### 5. Telegram Integration
- Notification when new Task is created
- Update status via Telegram
- Delete notification

### 6. Dashboard & Reporting
- Dashboard data aggregation
- Department performance scoring
- Statistics and Analytics

### 7. Scoring System
- Calculate Department score based on Task count
- Monthly scoring
- Auto-deduction when Task > 3 per month

## API Endpoints

### Problem/Task Routes
```
GET    /api/v1/problem/list           - List tasks with pagination
GET    /api/v1/problem/search/:query  - Search tasks
POST   /api/v1/problem/create         - Create new task
GET    /api/v1/problem/:id            - Get task detail
PUT    /api/v1/problem/update/:id     - Update task
DELETE /api/v1/problem/delete/:id     - Delete task
PUT    /api/v1/updateTaskStatus       - Update task status only
```

### IP Phone Routes
```
GET    /api/v1/ipphone/list           - List IP phones
GET    /api/v1/ipphone/listall        - List all IP phones
GET    /api/v1/ipphone/search/:query  - Search IP phones
POST   /api/v1/ipphone/create         - Create IP phone
PUT    /api/v1/ipphone/update/:id     - Update IP phone
DELETE /api/v1/ipphone/delete/:id     - Delete IP phone
```

### Department Routes
```
GET    /api/v1/department/list        - List departments
GET    /api/v1/department/listall     - List all departments
GET    /api/v1/department/search/:query - Search departments
POST   /api/v1/department/create      - Create department
GET    /api/v1/department/:id         - Get department detail
PUT    /api/v1/department/update/:id  - Update department
DELETE /api/v1/department/delete/:id  - Delete department
```

### Branch Routes
```
GET    /api/v1/branch/list            - List branches
POST   /api/v1/branch/create          - Create branch
GET    /api/v1/branch/:id             - Get branch detail
PUT    /api/v1/branch/update/:id      - Update branch
DELETE /api/v1/branch/delete/:id      - Delete branch
```

### Program Routes
```
GET    /api/v1/program/list           - List programs
POST   /api/v1/program/create         - Create program
PUT    /api/v1/program/update/:id     - Update program
DELETE /api/v1/program/delete/:id     - Delete program
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

### Dashboard & Scoring Routes
```
GET    /api/v1/dashboard/data         - Get dashboard data
GET    /api/v1/scores/list            - List department scores
GET    /api/v1/scores/:id             - Get score detail
PUT    /api/v1/scores/update/:id      - Update score
DELETE /api/v1/scores/delete/:id      - Delete score
GET    /api/v1/users                  - Get users list
```

## Database Schema (Main)

### tasks table
- id, ticket_no, phone_id, system_id, department_id
- text, status, assignto, message_id
- created_at, updated_at, created_by, updated_by

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
- id, name
- created_at, updated_at, deleted_at

### users table
- id, username, password, role
- created_at, updated_at, deleted_at

### scores table
- id, department_id, year, month, score

## System Workflow

### 1. Task Creation Flow
1. Receive Task data from Frontend
2. Auto-generate Ticket Number
3. Save to Database
4. Update Department Score
5. Send Notification via Telegram (if enabled)

### 2. Telegram Integration
- Send message when new Task is created
- Update message when Task is edited
- Delete message when Task is deleted
- Store message_id for Update/Delete

### 3. Scoring System
- Calculate monthly score per Department
- Start with 100 points
- Deduct 1 point when Task > 3 per month

### 4. Pagination System
- Supports Query Parameters: page, limit
- Default: page=1, limit=10
- Response includes pagination metadata

## Environment Configuration

### Development (.env.dev)
```
PORT=5001
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=password
DB_NAME=report_db
```

### Production (.env.prod)
```
PORT=5000
DB_HOST=production_host
DB_PORT=3306
DB_USER=prod_user
DB_PASS=prod_password
DB_NAME=report_db
```

## Deployment

### Docker
```bash
# Build image
docker build -t reports-api .

# Run container
docker run -p 5000:5000 --env-file .env reports-api
```

### Manual Run
```bash
# Development
go run main.go dev

# Production  
go run main.go prod
```

## Security Features
- CORS configuration
- Input validation
- SQL injection prevention
- Session management
- Role-based access control
- Soft delete implementation

## Performance Optimization
- Database connection pooling
- Memory limit (384MB)
- CPU limit (2 cores)
- Garbage collection optimization
- Query optimization with indexes

## Testing
- Unit tests สำหรับ Handlers
- Integration tests สำหรับ Database
- Pagination testing
- Mock testing
- Performance testing

## Usage with ChatGPT
When you need ChatGPT to help develop or modify code, reference this file to:
1. Understand project structure
2. Know existing API endpoints
3. Understand Database schema
4. Know Business logic and workflow
5. Understand Technology stack used

## Version History
- v1.7.0: Telegram update และ Assign_to feature
- v1.6.0: Assignto feature
- v1.5.0: Ticket number system
- v1.4.0: Unit testing
- v1.3.0: Pagination system
- v1.2.0: Telegram integration
- v1.1.0: Initial release