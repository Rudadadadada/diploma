package utils

import (
	"time"
)

func TruncateTime(inputTime time.Time) string {
	return inputTime.Format("2006-01-02 15:04:05")
}