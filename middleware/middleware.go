package middleware

import (
	"net/http"
)

// HeaderMiddleware เพิ่ม headers ที่ต้องการให้กับทุก response
func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// กำหนด headers ที่ต้องการใช้ร่วมกัน
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Applications", "API")
		w.Header().Set("Version", "1.0")
		
		// ส่งต่อไปยัง handler ถัดไป
		next.ServeHTTP(w, r)
	})
}