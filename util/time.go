package util

import (
	"fmt"
	"strings"
	"time"
)

func FormatDuration(d time.Duration) string {
	var parts []string

	if hours := d / time.Hour; hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
		d %= time.Hour
	}
	if minutes := d / time.Minute; minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dmin", minutes))
		d %= time.Minute
	}
	if seconds := d / time.Second; seconds > 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
		d %= time.Second
	}
	if ms := d.Milliseconds(); ms > 0 {
		parts = append(parts, fmt.Sprintf("%dms", ms))
	}

	return strings.Join(parts, " ")
}
