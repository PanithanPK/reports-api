# ChatGPT Context - Reports API Development Guide

## สำหรับ ChatGPT: วิธีการทำงานกับโปรเจคนี้

### 1. โครงสร้างโค้ดและหลักการ

#### Architecture Pattern
- **MVC Pattern**: Models, Handlers (Controllers), Routes
- **Repository Pattern**: Database operations ใน handlers
- **Middleware Pattern**: Authentication, CORS, Recovery
- **Clean Architecture**: แยก Business Logic และ Data Layer

#### Code Organization
```
reports-api/
├── main.go              # Application entry point
├── db/                  # Database connection layer
├── models/              # Data structures และ DTOs
├── handlers/            # Business logic (Controllers)
├── routes/              # Route definitions
├── utils/               # Utility functions
├── middleware/          # Middleware functions
└── tests/               # Test files
```

### 2. การพัฒนาฟีเจอร์ใหม่

#### ขั้นตอนการเพิ่มฟีเจอร์ใหม่:
1. **สร้าง Model** ใน `models/` directory
2. **สร้าง Handler** ใน `handlers/` directory  
3. **เพิ่ม Route** ใน `routes/routes.go`
4. **เขียน Test** ใน `tests/` directory

#### ตัวอย่างการสร้าง CRUD:
```go
// 1. Model (models/exampleModel.go)
type Example struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// 2. Handler (handlers/example.go)
func ListExamplesHandler(c *fiber.Ctx) error {
    // Database query
    // Return JSON response
}

// 3. Route (routes/routes.go)
r.Get("/api/v1/example/list", handlers.ListExamplesHandler)
```

### 3. Database Operations

#### การเชื่อมต่อ Database
- ใช้ `db.DB` global variable
- Connection pooling configured
- Support prepared statements

#### Query Patterns
```go
// SELECT
rows, err := db.DB.Query("SELECT * FROM table WHERE condition = ?", value)

// INSERT
result, err := db.DB.Exec("INSERT INTO table (col1, col2) VALUES (?, ?)", val1, val2)

// UPDATE
_, err := db.DB.Exec("UPDATE table SET col1 = ? WHERE id = ?", newVal, id)

// DELETE (Soft Delete)
_, err := db.DB.Exec("UPDATE table SET deleted_at = ? WHERE id = ?", time.Now(), id)
```

### 4. Response Patterns

#### Success Response
```go
return c.JSON(fiber.Map{
    "success": true,
    "data": data,
})
```

#### Error Response
```go
return c.Status(400).JSON(fiber.Map{
    "error": "Error message",
})
```

#### Paginated Response
```go
return c.JSON(models.PaginatedResponse{
    Success: true,
    Data: items,
    Pagination: models.PaginationResponse{
        Page: page,
        Limit: limit,
        Total: total,
        TotalPages: totalPages,
    },
})
```

### 5. Business Logic Patterns

#### Task Management Flow
1. **Create Task**: Generate ticket → Save to DB → Update score → Send Telegram
2. **Update Task**: Validate → Update DB → Update Telegram message
3. **Delete Task**: Soft delete → Delete Telegram message

#### Scoring System Logic
- เริ่มต้น 100 คะแนนต่อเดือน
- ลดคะแนนเมื่อ Task > 3 ต่อเดือน
- คำนวณแยกตาม Department และ Month/Year

#### Telegram Integration
- ส่งข้อความเมื่อสร้าง Task ใหม่
- Update ข้อความเมื่อแก้ไข Task
- ลบข้อความเมื่อลบ Task
- เก็บ message_id ใน Database

### 6. Error Handling Patterns

#### Database Errors
```go
if err != nil {
    log.Printf("Database error: %v", err)
    return c.Status(500).JSON(fiber.Map{"error": "Database error"})
}
```

#### Validation Errors
```go
if err := c.BodyParser(&req); err != nil {
    return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
}
```

#### Not Found Errors
```go
if err == sql.ErrNoRows {
    return c.Status(404).JSON(fiber.Map{"error": "Resource not found"})
}
```

### 7. Testing Patterns

#### Handler Testing
```go
func TestHandler(t *testing.T) {
    app := fiber.New()
    app.Get("/test", handler)
    
    req := httptest.NewRequest("GET", "/test", nil)
    resp, _ := app.Test(req)
    
    assert.Equal(t, 200, resp.StatusCode)
}
```

#### Database Testing
```go
func TestDatabase(t *testing.T) {
    // Setup test database
    // Run test queries
    // Assert results
    // Cleanup
}
```

### 8. Common Development Tasks

#### เพิ่ม API Endpoint ใหม่
1. สร้าง Model struct
2. สร้าง Handler function
3. เพิ่ม Route
4. เขียน Test
5. Update documentation

#### แก้ไข Database Schema
1. Update Model struct
2. แก้ไข Query ใน Handler
3. Update Test cases
4. Migration script (ถ้าจำเป็น)

#### เพิ่ม Validation
1. เพิ่ม validation ใน Handler
2. Return appropriate error message
3. Update Test cases

### 9. Performance Considerations

#### Database Optimization
- ใช้ Prepared Statements
- Limit query results
- Use appropriate indexes
- Connection pooling

#### Memory Management
- Defer close resources
- Avoid memory leaks
- Use appropriate data types

#### Response Optimization
- Pagination for large datasets
- Compress responses
- Cache frequently accessed data

### 10. Security Best Practices

#### Input Validation
- Validate all user inputs
- Use parameterized queries
- Sanitize data

#### Authentication
- Session-based authentication
- Role-based access control
- Secure password handling

#### CORS Configuration
- Allow specific origins
- Restrict methods and headers
- Handle preflight requests

### 11. Debugging และ Logging

#### Logging Patterns
```go
log.Printf("Operation successful: %s", operation)
log.Printf("Error occurred: %v", err)
```

#### Debug Information
- Log request parameters
- Log database queries
- Log response data
- Track execution time

### 12. Environment Management

#### Development vs Production
- Use environment-specific configs
- Different database connections
- Different logging levels
- Different security settings

#### Environment Variables
```go
port := os.Getenv("PORT")
dbHost := os.Getenv("DB_HOST")
```

## สำหรับ ChatGPT: คำแนะนำในการช่วยพัฒนา

### เมื่อได้รับคำขอช่วยเหลือ:
1. **อ่าน Context**: ดูโครงสร้างโปรเจคและ Business Logic
2. **ตรวจสอบ Pattern**: ใช้ Pattern ที่มีอยู่ในโปรเจค
3. **Follow Convention**: ตั้งชื่อไฟล์และ Function ตาม Convention
4. **Include Error Handling**: เพิ่ม Error Handling ที่เหมาะสม
5. **Write Tests**: แนะนำการเขียน Test cases
6. **Consider Performance**: คำนึงถึง Performance และ Security

### ข้อมูลสำคัญที่ต้องจำ:
- ใช้ Fiber Framework (ไม่ใช่ Gin หรือ Echo)
- Database คือ MySQL (ไม่ใช่ PostgreSQL)
- ใช้ Session-based Authentication
- มี Telegram Integration
- มี Scoring System สำหรับ Department
- ใช้ Soft Delete Pattern
- มี Pagination System

### การตอบคำถาม:
- ให้ตัวอย่างโค้ดที่สมบูรณ์
- อธิบาย Business Logic
- แนะนำ Best Practices
- ระบุไฟล์ที่ต้องแก้ไข
- แนะนำการ Test