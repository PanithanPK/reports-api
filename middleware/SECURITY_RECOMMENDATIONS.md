# Security Recommendations for Reports API

## Security Issues Found and Recommendations for Fixes

### 1. Password Management
**Issue**: Passwords are stored in plaintext format in the database and are compared directly in the LoginHandler function

**Recommendation**: 
1. Use password hashing with bcrypt or Argon2id
2. Do not store passwords in plaintext format
3. Example fix:
```go
// When registering a user
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
if err != nil {
    http.Error(w, "Failed to hash password", http.StatusInternalServerError)
    return
}
// Store hashedPassword in database

// When logging in
err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordFromDB), []byte(credentials.Password))
if err != nil {
    http.Error(w, "Invalid username or password", http.StatusUnauthorized)
    return
}
```

### 2. JWT Authentication
**Issue**: The current system does not use JWT for authentication even though the structure is prepared in LoginResponse

**Recommendation**:
1. Create and validate JWT tokens for authentication
2. Set appropriate token expiration time
3. Example JWT creation:
```go
// Create JWT token
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

### 3. SQL Injection Prevention
**Issue**: Direct queries are used without prepared statements in some parts

**Recommendation**:
1. Use prepared statements for all database queries
2. Validate and sanitize user input before using in SQL commands

### 4. Authentication Middleware
**Issue**: No middleware for checking access permissions to API endpoints

**Recommendation**:
1. Create middleware to validate JWT tokens and access permissions
2. Example middleware:
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Remove "Bearer " from token
        tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
        
        // Validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        
        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Extract data from token
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Add user data to context
        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 5. CSRF (Cross-Site Request Forgery) Protection
**Issue**: No CSRF protection implemented

**Recommendation**:
1. Use CSRF protection middleware such as gorilla/csrf
2. Example usage:
```go
CSRF := csrf.Protect(
    []byte(os.Getenv("CSRF_KEY")),
    csrf.Secure(true),
    csrf.HttpOnly(true),
)
http.ListenAndServe(":5000", CSRF(r))
```

### 6. Session Management
**Issue**: Session management is incomplete, only cookie deletion on logout

**Recommendation**:
1. Use session management library such as gorilla/sessions
2. Store session IDs in database or Redis for session revocation capability

### 7. Sensitive Data Storage in Environment Variables
**Issue**: Sensitive data such as database passwords are stored in .env but without security verification

**Recommendation**:
1. Do not store default passwords in code
2. Use secrets management service in production environment
3. Ensure .env files are not added to git repository

### 8. API Rate Limiting
**Issue**: No API rate limiting implemented

**Recommendation**:
1. Use rate limiting middleware such as tollbooth
2. Example usage:
```go
limiter := tollbooth.NewLimiter(1, nil) // 1 request per second
r.Handle("/api/v1/sensitive-endpoint", tollbooth.LimitHandler(limiter, sensitiveHandler))
```

### 9. Data Logging
**Issue**: Good logging is implemented but may expose sensitive data

**Recommendation**:
1. Do not log sensitive data such as passwords, tokens, or personal information
2. Use appropriate logging levels for each environment

### 10. CORS Configuration
**Issue**: CORS configuration allows all domains to access (`AllowedOrigins: []string{"*"}`)

**Recommendation**:
1. Restrict domains that can access the API
2. Example configuration:
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

## Additional Improvement Guidelines

1. **Sensitive Data Encryption**: Encrypt sensitive data in the database
2. **HTTPS Usage**: Ensure HTTPS is used in production environment
3. **Data Validation**: Use libraries like go-playground/validator for input validation
4. **Security Testing**: Conduct regular security testing
5. **Library Updates**: Regularly check and update libraries to prevent vulnerabilities
6. **Content Security Policy**: Add CSP headers to prevent XSS
7. **Data Access Control**: Verify user permissions for requested data access
8. **Prepared Statements**: Use prepared statements for all database queries
9. **Error Handling**: Do not expose sensitive data in error messages
10. **Data Backup**: Implement secure backup systems and regularly test data recovery