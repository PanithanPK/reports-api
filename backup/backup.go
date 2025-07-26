package backup

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type BackupService struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	BackupDir  string
}

func NewBackupService() *BackupService {
	return &BackupService{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASS"),
		DBName:     os.Getenv("DB_NAME"),
		BackupDir:  "./backups",
	}
}

func (bs *BackupService) CreateBackup() error {
	// Create backup directory if not exists
	if err := os.MkdirAll(bs.BackupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_backup_%s.sql", bs.DBName, timestamp)
	filepath := filepath.Join(bs.BackupDir, filename)

	// Get container name from environment or use default
	containerName := os.Getenv("MYSQL_CONTAINER_NAME")
	if containerName == "" {
		containerName = "mysql" // default container name
	}

	// Execute docker exec mysqldump command
	cmd := exec.Command("docker", "exec", containerName, "mysqldump",
		"-u", bs.DBUser,
		fmt.Sprintf("-p%s", bs.DBPassword),
		bs.DBName,
	)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("docker mysqldump failed: %v", err)
	}

	// Write backup to file
	if err := os.WriteFile(filepath, output, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %v", err)
	}

	log.Printf("Backup created successfully: %s", filepath)
	return nil
}

func (bs *BackupService) StartScheduledBackup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if err := bs.CreateBackup(); err != nil {
				log.Printf("Scheduled backup failed: %v", err)
			}
		}
	}()
}

func (bs *BackupService) CleanOldBackups(days int) error {
	files, err := filepath.Glob(filepath.Join(bs.BackupDir, "*.sql"))
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			os.Remove(file)
			log.Printf("Removed old backup: %s", file)
		}
	}
	return nil
}