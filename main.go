package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"reports-api/db"
	"reports-api/middleware"
	"reports-api/routes"
	"runtime"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// CurrentEnvironment เก็บสภาพแวดล้อมปัจจุบัน (dev, prod, หรือ default)
var CurrentEnvironment string

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

func main() {
	logger.Info.Println("🚀 Starting Problem Report System...")

	// Parse command-line arguments
	env := flag.String("env", "", "Environment (dev or prod)")
	devFlag := flag.Bool("d", false, "Use development environment")
	prodFlag := flag.Bool("p", false, "Use production environment")
	flag.Parse()

	// Check flags and positional arguments
	selectedEnv := *env
	if *devFlag {
		selectedEnv = "dev"
	} else if *prodFlag {
		selectedEnv = "prod"
	} else if len(flag.Args()) > 0 && selectedEnv == "" {
		// Check if environment is passed as a positional argument
		arg := flag.Args()[0]
		if arg == "dev" || arg == "prod" {
			selectedEnv = arg
		}
	}

	// Load environment variables based on environment
	envFile := ".env"
	CurrentEnvironment = "default"
	if selectedEnv == "dev" {
		envFile = ".env.dev"
		CurrentEnvironment = "dev"
		logger.Info.Println("🔧 Running in DEVELOPMENT environment")
	} else if selectedEnv == "prod" {
		envFile = ".env.prod"
		CurrentEnvironment = "prod"
		logger.Info.Println("🔧 Running in PRODUCTION environment")
	} else {
		logger.Info.Println("🔧 Running with default environment")
	}

	// เก็บสภาพแวดล้อมใน environment variable เพื่อให้โค้ดส่วนอื่นเข้าถึงได้
	os.Setenv("APP_ENV", CurrentEnvironment)

	// Load environment variables
	err := godotenv.Load(envFile)
	if err != nil {
		logger.Warn.Printf("⚠️ %s file not found, using environment variables", envFile)
	} else {
		logger.Info.Printf("✅ Environment variables loaded from %s file", envFile)
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

	// mux.CORSMethodMiddleware(r)

	// Add middleware
	// r.Use(middleware.RecoveryMiddleware) // ต้องใส่เป็นตัวแรกเพื่อจับ panic ในทุก middleware อื่นๆ
	// r.Use(loggingMiddleware)
	// r.Use(middleware.RateLimitMiddleware(60)) // จำกัดการเข้าถึงที่ 60 คำขอต่อวินาที
	// r.Use(middleware.BasicSecurityHeadersMiddleware)
	
	// เพิ่ม HeaderMiddleware เพื่อกำหนด headers ให้กับทุก response
	r.Use(middleware.HeaderMiddleware)
	logger.Info.Println("✅ HeaderMiddleware added for common response headers")

	// Setup CORS using rs/cors package
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	})
	
	// Use the CORS handler
	handler := c.Handler(r)
	logger.Info.Println("🌐 CORS enabled using github.com/rs/cors package")

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./fontend"))))
	logger.Info.Println("📁 Static file server configured")

	// Serve index.html at root
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Printf("📄 Serving index.html to %s", r.RemoteAddr)
		http.ServeFile(w, r, "./fontend/index.html")
	})

	// Register API routes
	logger.Info.Println("🔗 Registering API routes...")
	routes.RegisterRoutes(r)
	logger.Info.Println("✅ API routes registered successfully")

	// Register Authentication routes
	logger.Info.Println("🔐 Registering Authentication routes...")
	routes.RegisterAuthRoutes(r)
	logger.Info.Println("✅ Authentication routes registered successfully")

	// Test route for RecoveryMiddleware
	r.HandleFunc("/test-panic", func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("🧪 Testing RecoveryMiddleware with a deliberate panic")
		panic("This is a test panic to verify RecoveryMiddleware is working")
	}).Methods("GET")
	logger.Info.Println("🧪 Test route for RecoveryMiddleware added at /test-panic")

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
	logger.Info.Printf("🚀 Server starting on http://localhost:%s", port)
	logger.Info.Println("🎯 Server is ready to handle requests!")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Error.Printf("❌ Server failed to start: %v", err)
		log.Fatal(err)
	}
}
