package common

import (
	"database/sql"
	"fmt"
	"log"
	"reports-api/db"
	"time"
)

// Alternative: Use string scanning then convert
func GetResolvedAtSafely(db *sql.DB, resolutionID int) (time.Time, error) {
	var resolvedAtStr string
	err := db.QueryRow(`SELECT DATE_FORMAT(resolved_at, '%Y-%m-%d %H:%i:%s') FROM resolutions WHERE id = ?`, resolutionID).Scan(&resolvedAtStr)
	if err != nil {
		return time.Time{}, err
	}

	if resolvedAtStr == "" {
		return time.Time{}, nil
	}

	resolvedAt, err := time.Parse("2006-01-02 15:04:05", resolvedAtStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time %s: %v", resolvedAtStr, err)
	}

	return resolvedAt, nil
}

func Fixtimefeature(timeStr string) string {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		log.Printf("Failed to parse created_at time '%s': %v", timeStr, err)
		parsedTime = time.Now() // ใช้เวลาปัจจุบันถ้า parse ไม่ได้
	}
	return parsedTime.Add(7 * time.Hour).Format("2006/01/02 15:04:05")
}

func Generateticketno() string {
	// create ticket as TK-DDMMYYYY-no using the latest number of that month/year + 1
	now := time.Now().Add(7 * time.Hour)
	dateStr := now.Format("20060102") // วันเดือนปี
	year := now.Year()
	month := int(now.Month())

	// get last ticket number for this month/year
	var lastNo sql.NullInt64
	err := db.DB.QueryRow(`
		SELECT MAX(CAST(RIGHT(ticket_no, 4) AS UNSIGNED)) 
		FROM tasks 
		WHERE YEAR(created_at) = ? AND MONTH(created_at) = ? 
		AND ticket_no REGEXP '^TK-[0-9]{8}-[0-9]{4}$'`, year, month).Scan(&lastNo)

	ticketNo := 1 // default to 1 if no records found
	if err != nil {
		log.Printf("Error getting last ticket no for month/year: %v", err)
	} else if lastNo.Valid {
		ticketNo = int(lastNo.Int64) + 1
	}
	ticket := fmt.Sprintf("TK-%s-%04d", dateStr, ticketNo)
	return ticket
}

// FormatResolvedAt รับ resolved_at และบวก 7 ชั่วโมง แล้ว format เป็น string
func FormatResolvedAt(resolvedAt time.Time) string {
	return resolvedAt.Add(7 * time.Hour).Format("02/01/2006 15:04:05")
}

// FormatResolvedAtFromString รับ resolved_at string แปลงเป็น time แล้วบวก 7 ชั่วโมง
func TimeFromString(TimeStr string) string {
	if TimeStr == "" {
		return ""
	}
	TimeAt, err := time.Parse("2006-01-02 15:04:05", TimeStr)
	if err != nil {
		return TimeStr
	}
	return TimeAt.Add(7 * time.Hour).Format("02/01/2006 15:04:05")
}
