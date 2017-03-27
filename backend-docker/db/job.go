package db

import (
	"time"
)

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
	Error  string
	Status string
	Code   int
}

// JobID exported
type JobID string

// VolumeID contains a docker volume ID
type VolumeID string

// Task contains tasks related to a specific job (used to serialize JSON)
type Task struct {
	ID            string
	Name          string
	ContainerID   string
	Error         string
	ExitCode      int
	ConsoleOutput string
}

// LatestOutput used to serialize console output to JSON
type LatestOutput struct {
	Name          string
	ConsoleOutput string
}
