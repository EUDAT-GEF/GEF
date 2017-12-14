package db

import (
	"time"
)

// Build describes status information about an image being built for a GEF service
type Build struct {
	ID           string
	ConnectionID ConnectionID
	Started      time.Time
	Duration     int64
	State        *BuildState
}

// BuildState keeps information about a build state
type BuildState struct {
	Status string
	Error  string
	Code   int
}
