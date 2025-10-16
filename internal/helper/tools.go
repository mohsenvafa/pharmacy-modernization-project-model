package helper

import (
	"context"
	"fmt"
	"time"
)

func CalculateAge(dob time.Time) int {
	if dob.IsZero() {
		return 0
	}
	now := time.Now()
	age := now.Year() - dob.Year()
	anniversary := time.Date(now.Year(), dob.Month(), dob.Day(), dob.Hour(), dob.Minute(), dob.Second(), dob.Nanosecond(), dob.Location())
	if now.Before(anniversary) {
		age--
	}
	if age < 0 {
		age = 0
	}
	return age
}

func FormatShortDate(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("Jan 2, 2006")
}

func WaitOrContext(ctx context.Context, seconds int) bool {
	if seconds <= 0 {
		return true
	}
	select {
	case <-time.After(time.Duration(seconds) * time.Second):
		return true
	case <-ctx.Done():
		return false
	}
}

func FormatDecimal(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func FormatInt(value int) string {
	return fmt.Sprintf("%d", value)
}

func FormatShortDateFromString(dateStr string) string {
	if dateStr == "" {
		return "-"
	}
	// Try to parse ISO 8601 format
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// If parsing fails, return the original string
		return dateStr
	}
	return FormatShortDate(t)
}
