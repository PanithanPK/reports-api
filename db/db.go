package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	log.Println("🔌 Initializing database connection...")

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Println("📋 Building DSN from environment variables...")

		// Get database configuration from environment variables
		dbUser := os.Getenv("DB_USER")
		if dbUser == "" {
			dbUser = "root"
			log.Println("⚠️ Using default DB_USER: root")
		} else {
			log.Printf("✅ Using DB_USER: %s", dbUser)
		}

		dbPass := os.Getenv("DB_PASS")
		if dbPass == "" {
			dbPass = "123456"
			log.Println("⚠️ Using default DB_PASS: PASS")
		} else {
			log.Println("✅ Using DB_PASS from environment")
		}

		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "192.168.1.173"
			log.Printf("⚠️ Using default DB_HOST: %s", dbHost)
		} else {
			log.Printf("✅ Using DB_HOST: %s", dbHost)
		}

		dbPort := os.Getenv("DB_PORT")
		if dbPort == "" {
			dbPort = "3306"
			log.Printf("⚠️ Using default DB_PORT: %s", dbPort)
		} else {
			log.Printf("✅ Using DB_PORT: %s", dbPort)
		}

		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "report_db"
			log.Printf("⚠️ Using default DB_NAME: %s", dbName)
		} else {
			log.Printf("✅ Using DB_NAME: %s", dbName)
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
		log.Printf("🔗 DSN built: %s@tcp(%s:%s)/%s", dbUser, dbHost, dbPort, dbName)
	} else {
		log.Println("✅ Using MYSQL_DSN from environment")
	}

	var err error
	log.Println("🔌 Opening database connection...")
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("❌ Error opening database: %v", err)
		return fmt.Errorf("error opening database: %w", err)
	}

	// Set connection pool settings
	log.Println("⚙️ Configuring connection pool...")
	DB.SetMaxOpenConns(10)                  // Limit maximum number of open connections
	DB.SetMaxIdleConns(5)                   // Limit idle connections
	DB.SetConnMaxLifetime(time.Hour)        // Connection max lifetime
	DB.SetConnMaxIdleTime(30 * time.Minute) // Max idle time
	log.Println("✅ Connection pool configured")

	log.Println("🔍 Testing database connection...")
	if err = DB.Ping(); err != nil {
		log.Printf("❌ Error connecting to database: %v", err)
		return fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("✅ Connected to MySQL database!")
	log.Printf("📊 Database stats - MaxOpenConns: %d, MaxIdleConns: %d", 10, 5)
	return nil
}

// GetDBStats returns current database statistics
func GetDBStats() sql.DBStats {
	if DB != nil {
		stats := DB.Stats()
		log.Printf("📊 DB Stats - OpenConnections: %d, InUse: %d, Idle: %d",
			stats.OpenConnections, stats.InUse, stats.Idle)
		return stats
	}
	return sql.DBStats{}
}
