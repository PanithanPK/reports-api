package main

import (
	"log"
	"net/http"
	"os"
	"reports-api/db"
	"reports-api/routes"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Custom logger with levels
type Logger struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}

var logger *Logger

func init() {
	// Set memory limit to 384MB
	debug.SetMemoryLimit(384 * 1024 * 1024) // 384MB in bytes

	// Optimize garbage collector
	debug.SetGCPercent(50)

	// Limit number of processors
	runtime.GOMAXPROCS(2)

	// Initialize custom logger
	logger = &Logger{
		Info:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		Warn:  log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime),
		Error: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime),
	}
}

// Middleware for logging HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		logger.Info.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Add logging for incoming requests
		logger.Info.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

		// Call next handler
		next.ServeHTTP(w, r)

		// Log response time
		duration := time.Since(start)
		logger.Info.Printf("Response: %s %s completed in %v", r.Method, r.URL.Path, duration)
	})
}

func main() {
	logger.Info.Println("🚀 Starting Problem Report System...")

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		logger.Warn.Println("⚠️ .env file not found, using environment variables")
	} else {
		logger.Info.Println("✅ Environment variables loaded from .env file")
	}

	// Initialize database connection
	logger.Info.Println("🔌 Initializing database connection...")
	if err := db.InitDB(); err != nil {
		logger.Error.Printf("❌ Database initialization failed: %v", err)
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer func() {
		if err := db.DB.Close(); err != nil {
			logger.Error.Printf("❌ Error closing database: %v", err)
		} else {
			logger.Info.Println("✅ Database connection closed")
		}
	}()

	// Create router
	r := mux.NewRouter()

	// Add logging middleware
	r.Use(loggingMiddleware)

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend"))))
	logger.Info.Println("📁 Static file server configured")

	// Serve index.html at root
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Printf("📄 Serving index.html to %s", r.RemoteAddr)
		http.ServeFile(w, r, "./frontend/index.html")
	})

	// Register API routes
	logger.Info.Println("🔗 Registering API routes...")
	routes.RegisterRoutes(r)
	logger.Info.Println("✅ API routes registered successfully")

	// Register Authentication routes
	logger.Info.Println("🔐 Registering Authentication routes...")
	routes.RegisterAuthRoutes(r)
	logger.Info.Println("✅ Authentication routes registered successfully")

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	logger.Info.Println("🌐 CORS configured")

	// Catch-all: serve index.html for all other GET requests (for SPA)
	// ต้องวางไว้หลัง API routes เพื่อไม่ให้จับ API requests
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && !strings.HasPrefix(r.URL.Path, "/static/") && !strings.HasPrefix(r.URL.Path, "/problemEntry/") {
			logger.Info.Printf("📄 Serving index.html (catch-all) to %s", r.RemoteAddr)
			http.ServeFile(w, r, "./frontend/index.html")
		} else {
			logger.Warn.Printf("⚠️ 404 Not Found: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			http.NotFound(w, r)
		}
	})

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
		logger.Info.Println("🔧 Using default port 5000")
	} else {
		logger.Info.Printf("🔧 Using port from environment: %s", port)
	}

	// Log system information
	logger.Info.Printf("💻 System Info - CPU Cores: %d, Memory Limit: 384MB", runtime.NumCPU())
	logger.Info.Printf("📊 Go Version: %s", runtime.Version())

	// Start server
	handler := c.Handler(r)
	logger.Info.Printf("🚀 Server starting on http://localhost:%s", port)
	logger.Info.Println("🎯 Server is ready to handle requests!")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Error.Printf("❌ Server failed to start: %v", err)
		log.Fatal(err)
	}
}
