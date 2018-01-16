package db

import (
	"time"
)

// Build describes status information about an image being built for a GEF service
type Build struct {
	ID           string
	ServiceID    ServiceID
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

// NewBuildStateOk creates a new BuildState with no error
func NewBuildStateOk(status string, code int) BuildState {
	return BuildState{
		Status: status,
		Error:  "",
		Code:   code,
	}
}

// NewBuildStateError creates a new BuildState with specified error
func NewBuildStateError(err string, code int) BuildState {
	return BuildState{
		Error:  err,
		Status: "Error",
		Code:   code,
	}
}
