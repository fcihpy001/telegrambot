package services

import (
	"fmt"
	"time"
)

// unix seconds -> yyyymmdd
func ToDay(ts int64) string {
	tm := time.Unix(ts, 0)
	return fmt.Sprintf("%d%02d%02d", tm.Year(), tm.Month(), tm.Day())
}
