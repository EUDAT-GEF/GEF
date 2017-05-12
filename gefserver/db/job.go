package db

import "time"

// Job stores the information about a service execution (used to serialize JSON)
type Job struct {
	ID           JobID
	ServiceID    ServiceID
	Input        string
	Created      time.Time
	State        *JobState
	InputVolume  VolumeID
	OutputVolume VolumeID
	Tasks        []Task
}

// JobState keeps information about a job state
type JobState struct {
	Status string
	Error  string
	Code   int
}

// NewJobStateOk creates a new JobState with no error
func NewJobStateOk(status string, code int) JobState {
	return JobState{
		Status: status,
		Error:  "",
		Code:   code,
	}
}

// NewJobStateError creates a new JobState with specified error
func NewJobStateError(err string, code int) JobState {
	return JobState{
		Error:  err,
		Status: "Error",
		Code:   code,
	}
}

// JobID exported
type JobID string

// VolumeID contains a docker volume ID
type VolumeID string

// ContainerID exported
type ContainerID string

// Task contains tasks related to a specific job (used to serialize JSON)
type Task struct {
	ID            string
	Name          string
	ContainerID   ContainerID
	Error         string
	ExitCode      int
	ConsoleOutput string
}

// LatestOutput used to serialize console output to JSON
type LatestOutput struct {
	Name          string
	ConsoleOutput string
}
