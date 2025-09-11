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

func Generateticketno() string {
	// create ticket as TK-DDMMYYYY-no using the latest number of that month/year + 1
	now := time.Now().Add(7 * time.Hour)
	dateStr := now.Format("02012006") // วันเดือนปี
	year := now.Year()
	month := int(now.Month())

	// get last ticket number for this month/year
	var lastNo int
	err := db.DB.QueryRow(`SELECT COALESCE(MAX(CAST(SUBSTRING(ticket_no, LENGTH(ticket_no)-4, 5) AS UNSIGNED)), 0) FROM tasks WHERE YEAR(created_at) = ? AND MONTH(created_at) = ?`, year, month).Scan(&lastNo)
	if err != nil {
		log.Printf("Error getting last ticket no for month/year: %v", err)
		lastNo = 0
	}
	// increment by 1
	ticketNo := lastNo + 1
	ticket := fmt.Sprintf("TK-%s-%05d", dateStr, ticketNo)
	return ticket
}
