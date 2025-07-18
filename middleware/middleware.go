package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

// JWTMiddleware verifies JWT tokens in Authorization header
// Currently disabled - will be implemented in future versions
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Token verification is currently disabled
		// Just pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

// RoleMiddleware checks if user has required role
// Currently disabled - will be implemented in future versions
func RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Role checking is currently disabled
			// Just pass the request to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware implements a simple rate limiter
// Note: For production, consider using a more robust solution like tollbooth or a Redis-based limiter
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	// Map to store client IPs and their request timestamps
	clients := make(map[string][]time.Time)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr

			// Clean up old requests (older than 1 minute)
			now := time.Now()
			if timestamps, exists := clients[clientIP]; exists {
				var recent []time.Time
				for _, t := range timestamps {
					if now.Sub(t) < time.Minute {
						recent = append(recent, t)
					}
				}
				clients[clientIP] = recent
			}

			// Check if client has exceeded rate limit
			if len(clients[clientIP]) >= requestsPerMinute {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Add current request timestamp
			clients[clientIP] = append(clients[clientIP], now)

			next.ServeHTTP(w, r)
		})
	}
}

// BasicSecurityHeadersMiddleware adds basic security headers to all responses
func BasicSecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add basic security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "DENY")

		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get origin from request
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "*") {
				// Allow all origins
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				// Check if origin is in allowed list
				originAllowed := false
				for _, allowedOrigin := range allowedOrigins {
					if origin == allowedOrigin {
						originAllowed = true
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}

				// If origin is not allowed, set a default response
				if !originAllowed && origin != "" {
					// Either set a default allowed origin or just continue without setting the header
					// which will effectively block the cross-origin request
				}
			}

			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "300")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Call next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware จัดการกับ panic ที่เกิดขึ้นระหว่างการประมวลผลคำขอ
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// สร้าง defer function เพื่อจับ panic
		defer func() {
			if err := recover(); err != nil {
				// สร้าง stack trace
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, false)]
				
				// บันทึกข้อผิดพลาดและ stack trace
				log.Printf("PANIC: %v\n%s", err, stack)
				
				// ส่งข้อความแสดงข้อผิดพลาดกลับไปยังไคลเอนต์
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"Internal Server Error", "message":"The server encountered an unexpected condition"}`)
			}
		}()
		
		// เรียกใช้ handler ถัดไป
		next.ServeHTTP(w, r)
	})
}
