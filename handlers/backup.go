package handlers

import (
	"reports-api/backup"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func BackupHandler(c *fiber.Ctx) error {
	bs := backup.NewBackupService()
	
	if err := bs.CreateBackup(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Backup failed: " + err.Error()})
	}
	
	return c.JSON(fiber.Map{"message": "Backup created successfully"})
}

func CleanBackupsHandler(c *fiber.Ctx) error {
	daysStr := c.Query("days")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 7 // default 7 days
	}
	
	bs := backup.NewBackupService()
	if err := bs.CleanOldBackups(days); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Clean failed: " + err.Error()})
	}
	
	return c.JSON(fiber.Map{"message": "Old backups cleaned"})
}