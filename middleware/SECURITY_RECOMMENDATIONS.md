# คำแนะนำด้านความปลอดภัยสำหรับ Reports API

## ประเด็นด้านความปลอดภัยที่พบและคำแนะนำในการแก้ไข

### 1. การจัดการรหัสผ่าน
**ปัญหา**: รหัสผ่านถูกเก็บในรูปแบบข้อความธรรมดา (plaintext) ในฐานข้อมูล และมีการเปรียบเทียบโดยตรงในฟังก์ชัน LoginHandler
```go
if credentials.Password != password {
    http.Error(w, "Invalid username or password", http.StatusUnauthorized)
    return
}
```

**คำแนะนำ**: 
1. ใช้การเข้ารหัสรหัสผ่านด้วย bcrypt หรือ Argon2id
2. ไม่เก็บรหัสผ่านในรูปแบบข้อความธรรมดา
3. ตัวอย่างการแก้ไข:
```go
// เมื่อลงทะเบียนผู้ใช้
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
if err != nil {
    http.Error(w, "Failed to hash password", http.StatusInternalServerError)
    return
}
// เก็บ hashedPassword ลงในฐานข้อมูล

// เมื่อเข้าสู่ระบบ
err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordFromDB), []byte(credentials.Password))
if err != nil {
    http.Error(w, "Invalid username or password", http.StatusUnauthorized)
    return
}
```

### 2. การใช้ JWT สำหรับการยืนยันตัวตน
**ปัญหา**: ระบบปัจจุบันไม่มีการใช้ JWT สำหรับการยืนยันตัวตน แม้ว่าจะมีการเตรียมโครงสร้างไว้แล้วใน LoginResponse

**คำแนะนำ**:
1. สร้างและตรวจสอบ JWT token สำหรับการยืนยันตัวตน
2. กำหนดอายุของ token ที่เหมาะสม
3. ตัวอย่างการสร้าง JWT:
```go
// สร้าง JWT token
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "id":       user.ID,
    "username": user.Username,
    "role":     user.Role,
    "exp":      time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix(),
})

tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
if err != nil {
    return "", err
}
```

### 3. การป้องกัน SQL Injection
**ปัญหา**: มีการใช้ query ตรงๆ โดยไม่มีการใช้ prepared statements ในบางส่วน

**คำแนะนำ**:
1. ใช้ prepared statements ทุกครั้งที่มีการ query ข้อมูล
2. ตรวจสอบและทำความสะอาดข้อมูลที่รับมาจากผู้ใช้ก่อนนำไปใช้ในคำสั่ง SQL

### 4. Middleware สำหรับการตรวจสอบสิทธิ์
**ปัญหา**: ไม่มี middleware สำหรับตรวจสอบสิทธิ์การเข้าถึง API endpoints

**คำแนะนำ**:
1. สร้าง middleware สำหรับตรวจสอบ JWT token และสิทธิ์การเข้าถึง
2. ตัวอย่าง middleware:
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // ตัดคำว่า "Bearer " ออกจาก token
        tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
        
        // ตรวจสอบ token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        
        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // ดึงข้อมูลจาก token
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // เพิ่มข้อมูลผู้ใช้ลงใน context
        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 5. การป้องกัน CSRF (Cross-Site Request Forgery)
**ปัญหา**: ไม่มีการป้องกัน CSRF

**คำแนะนำ**:
1. ใช้ middleware สำหรับป้องกัน CSRF เช่น gorilla/csrf
2. ตัวอย่างการใช้งาน:
```go
CSRF := csrf.Protect(
    []byte(os.Getenv("CSRF_KEY")),
    csrf.Secure(true),
    csrf.HttpOnly(true),
)
http.ListenAndServe(":5000", CSRF(r))
```

### 6. การจัดการ Session
**ปัญหา**: การจัดการ session ยังไม่สมบูรณ์ มีเพียงการลบ cookie เมื่อ logout

**คำแนะนำ**:
1. ใช้ library สำหรับจัดการ session เช่น gorilla/sessions
2. เก็บ session ID ในฐานข้อมูลหรือ Redis เพื่อให้สามารถยกเลิก session ได้

### 7. การเก็บข้อมูลสำคัญใน Environment Variables
**ปัญหา**: มีการเก็บข้อมูลสำคัญเช่นรหัสผ่านฐานข้อมูลใน .env แต่ไม่มีการตรวจสอบความปลอดภัย

**คำแนะนำ**:
1. ไม่ควรเก็บรหัสผ่านเริ่มต้น (default) ในโค้ด
2. ใช้ secrets management service ในสภาพแวดล้อมการผลิต
3. ตรวจสอบว่าไฟล์ .env ไม่ถูกเพิ่มใน git repository

### 8. การจำกัดอัตราการเรียกใช้ API (Rate Limiting)
**ปัญหา**: ไม่มีการจำกัดอัตราการเรียกใช้ API

**คำแนะนำ**:
1. ใช้ middleware สำหรับจำกัดอัตราการเรียกใช้ API เช่น tollbooth
2. ตัวอย่างการใช้งาน:
```go
limiter := tollbooth.NewLimiter(1, nil) // 1 request per second
r.Handle("/api/v1/sensitive-endpoint", tollbooth.LimitHandler(limiter, sensitiveHandler))
```

### 9. การบันทึกข้อมูล (Logging)
**ปัญหา**: มีการบันทึกข้อมูลที่ดีแล้ว แต่อาจเปิดเผยข้อมูลสำคัญ

**คำแนะนำ**:
1. ไม่ควรบันทึกข้อมูลสำคัญเช่นรหัสผ่าน token หรือข้อมูลส่วนบุคคล
2. ใช้ระดับการบันทึกที่เหมาะสมในแต่ละสภาพแวดล้อม

### 10. การตั้งค่า CORS
**ปัญหา**: การตั้งค่า CORS อนุญาตให้ทุกโดเมนเข้าถึงได้ (`AllowedOrigins: []string{"*"}`)

**คำแนะนำ**:
1. จำกัดโดเมนที่สามารถเข้าถึง API ได้
2. ตัวอย่างการตั้งค่า:
```go
c := cors.New(cors.Options{
    AllowedOrigins:   []string{"https://yourdomain.com", "https://app.yourdomain.com"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: true,
    MaxAge:           300,
})
```

## แนวทางการปรับปรุงเพิ่มเติม

1. **การเข้ารหัสข้อมูลสำคัญ**: เข้ารหัสข้อมูลสำคัญในฐานข้อมูล
2. **การใช้ HTTPS**: ตรวจสอบว่าใช้ HTTPS ในสภาพแวดล้อมการผลิต
3. **การตรวจสอบความถูกต้องของข้อมูล**: ใช้ library เช่น go-playground/validator สำหรับตรวจสอบข้อมูลที่รับมา
4. **การทดสอบความปลอดภัย**: ทำการทดสอบความปลอดภัยอย่างสม่ำเสมอ
5. **การอัปเดตไลบรารี**: ตรวจสอบและอัปเดตไลบรารีที่ใช้อย่างสม่ำเสมอเพื่อป้องกันช่องโหว่
6. **การใช้ Content Security Policy**: เพิ่ม header CSP เพื่อป้องกัน XSS
7. **การตรวจสอบการเข้าถึงข้อมูล**: ตรวจสอบว่าผู้ใช้มีสิทธิ์เข้าถึงข้อมูลที่ร้องขอ
8. **การใช้ Prepared Statements**: ใช้ prepared statements ทุกครั้งที่มีการ query ข้อมูล
9. **การจัดการข้อผิดพลาด**: ไม่เปิดเผยข้อมูลสำคัญในข้อความแสดงข้อผิดพลาด
10. **การสำรองข้อมูล**: มีระบบสำรองข้อมูลที่ปลอดภัยและทดสอบการกู้คืนข้อมูลอย่างสม่ำเสมอ