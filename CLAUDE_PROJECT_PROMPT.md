# Reports API - Complete Project Context for Claude

## ภาพรวมโปรเจค
Reports API เป็นระบบ Backend API ที่พัฒนาด้วย **Go (Golang) 1.23** และ **Fiber Framework v2** สำหรับจัดการระบบรายงานปัญหา (Problem Reporting System) ในองค์กร โดยมีการเชื่อมต่อกับฐานข้อมูล **MySQL** และระบบแจ้งเตือนผ่าง **Telegram Bot**

## เทคโนโลยีและ Dependencies หลัก
```go
// go.mod - Key Dependencies
- Go 1.23.0 (toolchain go1.24.4)
- github.com/gofiber/fiber/v2 v2.52.5 (Web Framework)
- github.com/go-sql-driver/mysql v1.9.3 (MySQL Driver)
- github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 (Telegram Bot)
- github.com/joho/godotenv v1.5.1 (Environment Variables)
- github.com/minio/minio-go/v7 v7.0.95 (File Storage)
- github.com/swaggo/fiber-swagger v1.3.0 (API Documentation)
```

## โครงสร้างโปรเจคและไฟล์สำคัญ

### 1. Main Application (`main.go`)
- รองรับ Environment: dev/prod/default
- Memory optimization (384MB limit, 2 CPU cores)
- CORS configuration สำหรับ multiple origins
- Static file serving
- Custom logging system
- Swagger UI integration

### 2. Database Layer (`db/db.go`)
- MySQL connection with connection pooling
- MaxOpenConns: 10, MaxIdleConns: 5
- Connection lifetime: 1 hour, Idle timeout: 30 minutes
- Environment-based configuration

### 3. Routes (`routes/routes.go`)
**Authentication Routes:**
```
POST /api/authEntry/login
POST /api/authEntry/registerUser
POST /api/authEntry/registerAdmin
PUT  /api/authEntry/updateUser
DELETE /api/authEntry/deleteUser
POST /api/authEntry/logout
```

**Problem/Task Management Routes:**
```
GET    /api/v1/problem/list
GET    /api/v1/problem/list/:query
GET    /api/v1/problem/list/:column/:query
POST   /api/v1/problem/create
GET    /api/v1/problem/:id
PUT    /api/v1/problem/update/:id
DELETE /api/v1/problem/delete/:id
PUT    /api/v1/problem/update/assignto/:id
```

**Resolution Routes:**
```
GET    /api/v1/resolution/:id
POST   /api/v1/resolution/create/:id
PUT    /api/v1/resolution/update/:id
DELETE /api/v1/resolution/delete/:id
```

**IP Phone Management Routes:**
```
GET    /api/v1/ipphone/list
GET    /api/v1/ipphone/list/:query
GET    /api/v1/ipphone/:id
GET    /api/v1/ipphone/listall
POST   /api/v1/ipphone/create
PUT    /api/v1/ipphone/update/:id
DELETE /api/v1/ipphone/delete/:id
```

**Program/System Management Routes:**
```
GET    /api/v1/program/list
GET    /api/v1/program/list/:query
POST   /api/v1/program/create
GET    /api/v1/program/type/list
GET    /api/v1/program/type/list/:query
POST   /api/v1/program/type/create
POST   /api/v1/program/type/update/:id
DELETE /api/v1/program/type/delete/:id
GET    /api/v1/program/:id
PUT    /api/v1/program/update/:id
DELETE /api/v1/program/delete/:id
```

**Department Management Routes:**
```
GET    /api/v1/department/list
GET    /api/v1/department/list/:query
GET    /api/v1/department/listall
POST   /api/v1/department/create
GET    /api/v1/department/:id
PUT    /api/v1/department/update/:id
DELETE /api/v1/department/delete/:id
```

**Branch Management Routes:**
```
GET    /api/v1/branch/list
GET    /api/v1/branch/list/:query
POST   /api/v1/branch/create
GET    /api/v1/branch/:id
PUT    /api/v1/branch/update/:id
DELETE /api/v1/branch/delete/:id
```

**Dashboard & Analytics Routes:**
```
GET    /api/v1/dashboard/data
GET    /api/v1/scores/list
GET    /api/v1/scores/:id
PUT    /api/v1/scores/update/:id
DELETE /api/v1/scores/delete/:id
```

**Responsibility Management Routes:**
```
GET    /api/v1/respons/list
GET    /api/v1/respons/:id
POST   /api/v1/respons/create
PUT    /api/v1/respons/update/:id
DELETE /api/v1/respons/delete/:id
```

## Database Schema และ Models

### Core Tables Structure:
1. **tasks** - หลักของระบบ
   - id, ticket_no, phone_id, system_id, department_id
   - issue_type, issue_else, text, status, assignto, assignto_id
   - reported_by, created_by, updated_by, telegram_id, solution_id
   - file_paths (JSON), resolved_at, created_at, updated_at

2. **ip_phones** - จัดการเบอร์โทรภายใน
   - id, number, name, department_id
   - created_at, updated_at, deleted_at, created_by, updated_by, deleted_by

3. **departments** - แผนกต่างๆ
   - id, name, branch_id
   - created_at, updated_at, deleted_at

4. **branches** - สาขาต่างๆ
   - id, name
   - created_at, updated_at, deleted_at, created_by, updated_by, deleted_by

5. **systems_program** - ระบบ/โปรแกรมต่างๆ
   - id, name, type, priority
   - created_at, updated_at, deleted_at, created_by, updated_by, deleted_by

6. **issue_types** - ประเภทปัญหา
   - id, name, created_at

7. **users** - ผู้ใช้งานระบบ
   - id, username, password, role
   - created_at, updated_at, deleted_at, created_by, updated_by, deleted_by

8. **scores** - คะแนนประจำเดือนของแผนก
   - department_id, year, month, score

9. **resolutions** - วิธีการแก้ไขปัญหา
   - id, tasks_id, text, telegram_id, file_paths, resolved_at

10. **telegram_chat** - ข้อมูล Telegram
    - id, chat_id, chat_name, report_id, assignto_id, solution_id

11. **responsibilities** - ผู้รับผิดชอบ
    - id, name, telegram_username

## ฟีเจอร์หลักของระบบ

### 1. Task/Problem Management System
- **Auto Ticket Generation**: รูปแบบ `TK-DDMMYYYY-XXXXX`
- **Status Tracking**: 0=รอดำเนินการ, 1=เสร็จสิ้น
- **File Upload Support**: รองรับรูปภาพผ่าน MinIO
- **Full-text Search**: ค้นหาข้อมูล Task ได้
- **Pagination**: รองรับการแบ่งหน้า
- **Overdue Calculation**: คำนวณเวลาที่เกินกำหนด

### 2. Telegram Integration System
- **Auto Notification**: แจ้งเตือนอัตโนมัติเมื่อสร้าง Task ใหม่
- **Status Updates**: อัปเดตสถานะใน Telegram เมื่อมีการเปลี่ยนแปลง
- **Assignment Notification**: แจ้งเตือนเมื่อมีการมอบหมายงาน
- **Photo Support**: ส่งรูปภาพพร้อมข้อความ
- **Reply System**: ระบบตอบกลับสำหรับ Resolution
- **Message Management**: ลบ/แก้ไขข้อความใน Telegram

### 3. File Management System
- **MinIO Integration**: จัดเก็บไฟล์ผ่าน MinIO
- **Multiple File Upload**: รองรับการอัปโหลดหลายไฟล์
- **File Naming Convention**: `{ticket_no}-{index}-{date}-{filename}`
- **Auto File Cleanup**: ลบไฟล์เก่าเมื่อมีการอัปเดต

### 4. Department Scoring System
- **Monthly Scoring**: คะแนนรายเดือนต่อแผนก
- **Auto Deduction**: ลดคะแนนเมื่อมี Task เกิน 3 ต่อเดือน
- **Score Tracking**: ติดตามคะแนนย้อนหลัง

### 5. User Authentication & Authorization
- **Session-based Auth**: ใช้ Session Cookie
- **Role-based Access**: user/admin roles
- **Soft Delete**: ไม่ลบข้อมูลจริง

### 6. Advanced Search & Filtering
- **Multi-column Search**: ค้นหาได้หลายคอลัมน์
- **Thai Language Support**: รองรับภาษาไทย
- **URL Decode**: รองรับ URL encoding
- **Query Sanitization**: ป้องกัน SQL injection

## Business Logic และ Workflow

### Task Creation Flow:
1. รับข้อมูลจาก Frontend (รองรับ multipart/form-data)
2. Generate Ticket Number อัตโนมัติ
3. Upload ไฟล์ไป MinIO (ถ้ามี)
4. บันทึกข้อมูลลง Database
5. Update Department Score
6. ส่ง Telegram Notification (ถ้าเปิดใช้)
7. บันทึก Message ID สำหรับการอัปเดตภายหลัง

### Task Update Flow:
1. ตรวจสอบการเปลี่ยนแปลงข้อมูล
2. จัดการไฟล์ใหม่ (ลบเก่า อัปโหลดใหม่)
3. อัปเดตข้อมูลใน Database
4. ตรวจสอบการเปลี่ยนผู้รับผิดชอบ
5. อัปเดต Telegram Message
6. ส่ง Notification ใหม่ (ถ้าเปลี่ยนผู้รับผิดชอบ)

### Resolution System:
1. สร้าง Resolution พร้อมไฟล์
2. อัปเดตสถานะ Task เป็น "เสร็จสิ้น"
3. ส่ง Reply Message ใน Telegram
4. อัปเดต Main Message ให้แสดงสถานะเสร็จสิ้น

## Environment Configuration

### Development (.env.dev):
```
PORT=5001
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=password
DB_NAME=report_db
BOT_TOKEN=telegram_bot_token
CHAT_ID=telegram_chat_id
End_POINT=minio_endpoint
ACCESS_KEY=minio_access_key
SECRET_ACCESSKEY=minio_secret_key
BUCKET_NAME=minio_bucket
```

### Production (.env.prod):
```
PORT=5000
DB_HOST=production_host
# ... similar structure
```

## Key Models และ Data Structures

### TaskRequest (สำหรับสร้าง Task):
```go
type TaskRequest struct {
    PhoneID      *int   `json:"phone_id"`
    SystemID     int    `json:"system_id"`
    IssueTypeID  int    `json:"issue_type"`
    IssueElse    string `json:"issue_else"`
    DepartmentID int    `json:"department_id"`
    Text         string `json:"text"`
    ReportedBy   string `json:"reported_by"`
    Telegram     bool   `json:"telegram"`
    CreatedBy    int    `json:"created_by"`
}
```

### TaskWithDetails (สำหรับแสดงผล):
```go
type TaskWithDetails struct {
    ID             int               `json:"id"`
    Ticket         string            `json:"ticket_no"`
    PhoneID        *int              `json:"phone_id"`
    Number         *int              `json:"number"`
    PhoneName      *string           `json:"phone_name"`
    SystemID       int               `json:"system_id"`
    SystemName     string            `json:"system_name"`
    DepartmentName string            `json:"department_name"`
    BranchName     string            `json:"branch_name"`
    Text           string            `json:"text"`
    Status         int               `json:"status"`
    FilePaths      map[string]string `json:"file_paths"`
    Overdue        interface{}       `json:"overdue"`
    CreatedAt      string            `json:"created_at"`
    UpdatedAt      string            `json:"updated_at"`
}
```

## Security Features
- **Input Validation**: ตรวจสอบข้อมูลนำเข้า
- **SQL Injection Prevention**: ใช้ Prepared Statements
- **CORS Configuration**: จำกัด Origins ที่อนุญาต
- **File Upload Security**: จำกัดประเภทไฟล์
- **Session Management**: จัดการ Session อย่างปลอดภัย
- **Soft Delete**: ไม่ลบข้อมูลจริงจากฐานข้อมูล

## Performance Optimizations
- **Connection Pooling**: จำกัดการเชื่อมต่อฐานข้อมูล
- **Memory Limit**: จำกัด Memory ที่ 384MB
- **CPU Limit**: จำกัดใช้ 2 CPU cores
- **Garbage Collection**: ปรับแต่ง GC ให้เหมาะสม
- **Query Optimization**: ใช้ Index และ JOIN อย่างเหมาะสม

## การใช้งานกับ Claude

เมื่อต้องการให้ Claude ช่วยพัฒนาหรือแก้ไขโค้ด ให้อ้างอิงข้อมูลนี้เพื่อ:

1. **เข้าใจ Architecture**: Go + Fiber + MySQL + Telegram + MinIO
2. **ทราบ API Endpoints**: ครบทุก endpoints และ HTTP methods
3. **เข้าใจ Database Schema**: ความสัมพันธ์ระหว่างตาราง
4. **ทราบ Business Logic**: Workflow การทำงานของระบบ
5. **เข้าใจ Models**: โครงสร้างข้อมูลทั้งหมด
6. **ทราบ Dependencies**: Libraries และ packages ที่ใช้
7. **เข้าใจ Environment**: การตั้งค่าสำหรับ dev/prod

## สิ่งสำคัญที่ต้องจำ:
- ระบบใช้ **Soft Delete** (deleted_at) ไม่ลบข้อมูลจริง
- **Ticket Number** สร้างอัตโนมัติตามรูปแบบ TK-DDMMYYYY-XXXXX
- **File Paths** เก็บเป็น JSON ในฐานข้อมูล
- **Telegram Integration** เป็นฟีเจอร์หลักของระบบ
- **Pagination** ใช้ page/limit parameters
- **Status**: 0=รอดำเนินการ, 1=เสร็จสิ้น
- **Timezone**: ใช้ +7 hours (Thailand timezone)
- **File Storage**: ใช้ MinIO แทน local storage
- **Search**: รองรับภาษาไทยและ URL encoding

## Version History:
- v1.7.0: Telegram update และ Assign_to feature
- v1.6.0: Assignto feature  
- v1.5.0: Ticket number system
- v1.4.0: Unit testing
- v1.3.0: Pagination system
- v1.2.0: Telegram integration
- v1.1.0: Initial release