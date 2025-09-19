# API Usage Guide

## Base URL
```
Development: http://localhost:5001
Production: https://your-domain.com
```

## Authentication

### Login
```http
POST /authEntry/login
Content-Type: application/json

{
  "username": "your_username",
  "password": "your_password"
}
```

### Register User
```http
POST /authEntry/registerUser
Content-Type: application/json

{
  "username": "new_user",
  "password": "password123",
  "email": "user@example.com"
}
```

## Problem Management

### List Problems
```http
GET /api/v1/problem/list
Authorization: Bearer <token>
```

### Create Problem
```http
POST /api/v1/problem/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "System Issue",
  "description": "Problem description",
  "priority": "high",
  "status": "open"
}
```

### Update Problem
```http
PUT /api/v1/problem/update/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "System Issue (Updated)",
  "status": 2
}
```

#### Status Values
- `0` - รอดำเนินการ (Pending)
- `1` - เสร็จสิ้น (Resolved)
- `2` - กำลังดำเนินการ (In Progress)

### Delete Problem
```http
DELETE /api/v1/problem/delete/{id}
Authorization: Bearer <token>
```

## Task Management

### List Tasks
```http
GET /api/v1/task/list
Authorization: Bearer <token>
```

### Create Task
```http
POST /api/v1/task/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "New Task",
  "description": "Task description",
  "assigned_to": "user_id",
  "due_date": "2024-12-31"
}
```

## IP Phone Management

### List IP Phones
```http
GET /api/v1/ipphone/list
Authorization: Bearer <token>
```

### Create IP Phone
```http
POST /api/v1/ipphone/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "phone_number": "1001",
  "ip_address": "192.168.1.100",
  "location": "Office Room A",
  "status": "active"
}
```

## Program Management

### List Programs
```http
GET /api/v1/program/list
Authorization: Bearer <token>
```

### Create Program
```http
POST /api/v1/program/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "New Program",
  "version": "1.0.0",
  "description": "Program description"
}
```

## Department Management

### List Departments
```http
GET /api/v1/department/list
Authorization: Bearer <token>
```

### Create Department
```http
POST /api/v1/department/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "New Department",
  "description": "Department description",
  "manager": "Manager Name"
}
```

## Branch Management

### List Branches
```http
GET /api/v1/branch/list
Authorization: Bearer <token>
```

### Create Branch
```http
POST /api/v1/branch/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "New Branch",
  "address": "Branch Address",
  "phone": "02-xxx-xxxx"
}
```

## Error Responses

### Standard Error Format
```json
{
  "error": "Error message description",
  "code": "ERROR_CODE",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### HTTP Status Codes
- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Access denied
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Pagination

For APIs that return lists, pagination can be used:

```http
GET /api/v1/problem/list?page=1&limit=10
```

### Response Format
```json
{
  "data": [...],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 50,
    "items_per_page": 10
  }
}
```

## File Upload

For file uploads:

```http
POST /api/v1/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: [binary data]
```

## Rate Limiting

API has request rate limits:
- 100 requests per minute for regular users
- 1000 requests per minute for admins

When limit exceeded, returns HTTP 429 Too Many Requests