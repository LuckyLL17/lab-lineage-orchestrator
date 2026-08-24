package platform

import "time"

type Clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now().UTC()
}

func (RealClock) Since(value time.Time) time.Duration {
	return time.Since(value)
}

func BeginningOfDay(value time.Time) time.Time {
	current := value.UTC()
	year, month, day := current.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
