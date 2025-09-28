package helper

import "time"

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
