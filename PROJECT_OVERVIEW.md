# Reports API - Project Overview

## ภาพรวมโปรเจค
Reports API เป็นระบบ Backend API ที่พัฒนาด้วย Go (Golang) และ Fiber Framework สำหรับจัดการระบบรายงานปัญหา (Problem Reporting System) ในองค์กร โดยมีการเชื่อมต่อกับฐานข้อมูล MySQL และระบบแจ้งเตือนผ่าน Telegram

## เทคโนโลยีที่ใช้
- **Backend**: Go 1.22 + Fiber Framework v2
- **Database**: MySQL
- **Authentication**: Session-based + JWT Token
- **Notification**: Telegram Bot API
- **Deployment**: Docker + Docker Compose

## โครงสร้างโปรเจค

### 1. Main Application (`main.go`)
- รองรับ Environment แบบ dev/prod/default
- การจัดการ Memory และ CPU optimization
- CORS configuration
- Static file serving
- Custom logging system

### 2. Database Layer (`db/`)
- **db.go**: การเชื่อมต่อ MySQL พร้อม Connection Pool
- รองรับ Environment Variables สำหรับ Database Config
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
- แยก Authentication routes และ Business logic routes

### 6. Utils (`utils/`)
- **pagination.go**: Pagination utilities
- **PAGINATION_GUIDE.md**: Pagination documentation

### 7. Tests (`tests/`)
- Unit tests และ Integration tests
- Database testing
- Handler testing
- Pagination testing

## ฟีเจอร์หลัก

### 1. Task/Problem Management
- สร้าง/แก้ไข/ลบ Task
- ระบบ Ticket Number แบบ Auto-generate
- การค้นหา Task แบบ Full-text search
- Pagination support
- Status tracking (0=Open, 1=In Progress, 2=Closed)
- Overdue calculation

### 2. IP Phone Management
- จัดการข้อมูล IP Phone
- เชื่อมโยงกับ Department
- Search และ Filter

### 3. Department & Branch Management
- จัดการ Department และ Branch
- Hierarchical structure (Branch -> Department)
- Department scoring system

### 4. User Authentication
- Session-based authentication
- Role-based access (user/admin)
- User registration/login/logout
- Password management

### 5. Telegram Integration
- แจ้งเตือนเมื่อมี Task ใหม่
- Update status ผ่าน Telegram
- Delete notification

### 6. Dashboard & Reporting
- Dashboard data aggregation
- Department performance scoring
- Statistics และ Analytics

### 7. Scoring System
- คำนวณคะแนน Department ตาม Task count
- Monthly scoring
- Auto-deduction เมื่อ Task > 3 ต่อเดือน

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

## Database Schema (หลัก)

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

## การทำงานของระบบ

### 1. Task Creation Flow
1. รับข้อมูล Task จาก Frontend
2. Generate Ticket Number แบบ Auto
3. บันทึกลง Database
4. Update Department Score
5. ส่ง Notification ผ่าน Telegram (ถ้าเปิดใช้)

### 2. Telegram Integration
- ส่งข้อความเมื่อมี Task ใหม่
- Update ข้อความเมื่อแก้ไข Task
- ลบข้อความเมื่อลบ Task
- เก็บ message_id สำหรับ Update/Delete

### 3. Scoring System
- คำนวณคะแนนรายเดือนต่อ Department
- เริ่มต้นที่ 100 คะแนน
- ลดคะแนน 1 คะแนนเมื่อมี Task > 3 ต่อเดือน

### 4. Pagination System
- รองรับ Query Parameters: page, limit
- Default: page=1, limit=10
- Response รวม pagination metadata

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

## การ Deploy

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

## การใช้งานกับ ChatGPT
เมื่อต้องการให้ ChatGPT ช่วยพัฒนาหรือแก้ไขโค้ด ให้อ้างอิงไฟล์นี้เพื่อ:
1. เข้าใจโครงสร้างโปรเจค
2. ทราบ API endpoints ที่มีอยู่
3. เข้าใจ Database schema
4. ทราบ Business logic และ Flow การทำงาน
5. เข้าใจ Technology stack ที่ใช้

## Version History
- v1.7.0: Telegram update และ Assign_to feature
- v1.6.0: Assignto feature
- v1.5.0: Ticket number system
- v1.4.0: Unit testing
- v1.3.0: Pagination system
- v1.2.0: Telegram integration
- v1.1.0: Initial release