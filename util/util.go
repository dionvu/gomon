package util

import "time"

// "2006-01-02T15:04:05Z07:00"
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
