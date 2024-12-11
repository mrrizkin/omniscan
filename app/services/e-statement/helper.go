package estatement

import (
	"time"
)

func calculateEStatementExpiry(timeBomb *string) (*time.Time, error) {
	if timeBomb == nil {
		t := time.Now().Add(24 * time.Hour)
		return &t, nil
	}

	timeStr := *timeBomb
	if timeStr == "" {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	return &t, err
}

func inArray(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
