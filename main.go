package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"reports-api/db"
	"time"

	"reports-api/routes"
	"runtime"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
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

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Reports API",
		BodyLimit:    100 * 1024 * 1024, // 50MB limit for file uploads
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	})

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://helpdesk.nopadol.com,http://helpdesk-dev.nopadol.com,http://10.0.2.94:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Content-Type,Authorization,X-Requested-With",
		AllowCredentials: true,
	}))
	logger.Info.Println("🌐 CORS enabled using Fiber built-in middleware")

	// Add custom header middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		return c.Next()
	})
	logger.Info.Println("✅ HeaderMiddleware added for common response headers")

	// Serve static files
	app.Static("/static", "./fontend")
	logger.Info.Println("📁 Static file server configured")

	// Health Check endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		version := "unknown"
		if data, err := os.ReadFile("package.json"); err == nil {
			var pkg map[string]interface{}
			if json.Unmarshal(data, &pkg) == nil {
				if v, ok := pkg["version"].(string); ok {
					version = v
				}
			}
		}
		return c.JSON(fiber.Map{
			"status":  "OK",
			"version": version,
		})
	})

	// Register Authentication routes
	logger.Info.Println("🔐 Registering routes...")
	routes.RegisterRoutes(app)
	logger.Info.Println("✅ Routes registered successfully")

	// Test route for RecoveryMiddleware
	app.Get("/test-panic", func(c *fiber.Ctx) error {
		logger.Info.Println("🧪 Testing RecoveryMiddleware with a deliberate panic")
		panic("This is a test panic to verify RecoveryMiddleware is working")
	})
	logger.Info.Println("🧪 Test route for RecoveryMiddleware added at /test-panic")

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
		logger.Info.Println("🔧 Using default port 5001")
	} else {
		logger.Info.Printf("🔧 Using port from environment: %s", port)
	}

	// Swagger UI route
	app.Get("/api/v1/swagger/*", fiberSwagger.WrapHandler)
	logger.Info.Println("📚 Swagger UI available at /api/v1/swagger/index.html")

	// Log system information
	logger.Info.Printf("💻 System Info - CPU Cores: %d, Memory Limit: 384MB", runtime.NumCPU())
	logger.Info.Printf("📊 Go Version: %s", runtime.Version())

	// Start server
	logger.Info.Printf("🚀 Server starting on http://localhost:%s", port)
	logger.Info.Println("🎯 Server is ready to handle requests!")

	if err := app.Listen(":" + port); err != nil {
		logger.Error.Printf("❌ Server failed to start: %v", err)
		log.Fatal(err)
	}
}
