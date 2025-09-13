package utils

import (
	"time"
)

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
