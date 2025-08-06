# Pagination Guide

## Overview
ระบบ pagination ได้ถูกเพิ่มเข้ามาในทุก list endpoints เพื่อปรับปรุงประสิทธิภาพและการจัดการข้อมูลจำนวนมาก

## Query Parameters

### Pagination Parameters
- `page`: หมายเลขหน้า (เริ่มต้นที่ 1)
- `limit`: จำนวนรายการต่อหน้า (เริ่มต้นที่ 10, สูงสุด 100)

## API Endpoints with Pagination

### Department List
```
GET /api/v1/department/list?page=1&limit=10
```

### Branch List
```
GET /api/v1/branch/list?page=1&limit=10
```

### Program List
```
GET /api/v1/program/list?page=1&limit=10
```

### IP Phone List
```
GET /api/v1/ipphone/list?page=1&limit=10
```

### Task List
```
GET /api/v1/task/list?page=1&limit=10
```

## Response Format

```json
{
  "success": true,
  "data": [
    // array of items
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 150,
    "total_pages": 15
  }
}
```

## Response Fields

### Pagination Object
- `page`: หน้าปัจจุบัน
- `limit`: จำนวนรายการต่อหน้า
- `total`: จำนวนรายการทั้งหมด
- `total_pages`: จำนวนหน้าทั้งหมด

## Examples

### ขอข้อมูลหน้าแรก (10 รายการ)
```bash
curl "http://localhost:5000/api/v1/department/list"
```

### ขอข้อมูลหน้าที่ 2 (20 รายการต่อหน้า)
```bash
curl "http://localhost:5000/api/v1/department/list?page=2&limit=20"
```

### ขอข้อมูลหน้าที่ 3 (5 รายการต่อหน้า)
```bash
curl "http://localhost:5000/api/v1/department/list?page=3&limit=5"
```

## Default Values
- หากไม่ระบุ `page`: จะใช้ค่า 1
- หากไม่ระบุ `limit`: จะใช้ค่า 10
- หาก `limit` มากกว่า 100: จะใช้ค่า 10
- หาก `page` น้อยกว่า 1: จะใช้ค่า 1

## Implementation Notes
- ข้อมูลจะถูกเรียงลำดับตาม ID จากมากไปน้อย (DESC)
- การนับจำนวนรายการทั้งหมดจะไม่รวมรายการที่ถูก soft delete
- ประสิทธิภาพจะดีขึ้นเมื่อมีข้อมูลจำนวนมาก