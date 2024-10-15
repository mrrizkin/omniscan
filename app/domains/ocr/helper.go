package ocr

import (
	"time"
)

func calculateMutasiExpiry(timeBomb string) (time.Time, error) {
	if timeBomb == "" {
		return time.Now().Add(24 * time.Hour), nil
	}
	return time.Parse("2006-01-02 15:04:05", timeBomb)
}
