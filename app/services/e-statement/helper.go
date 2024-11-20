package estatement

import (
	"time"
)

func calculateEStatementExpiry(timeBomb string) (time.Time, error) {
	if timeBomb == "" {
		return time.Now().Add(24 * time.Hour), nil
	}
	return time.Parse("2006-01-02 15:04:05", timeBomb)
}

func inArray(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
