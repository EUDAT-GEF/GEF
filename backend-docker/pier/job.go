package pier

import (
	"sort"
	"sync"
	"time"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
	"bytes"
)

// Job stores the information about a service execution
type Job struct {
	ID           JobID
	ServiceID    ServiceID
	Input        string
	Created      time.Time
	State        *JobState
	InputVolume  VolumeID
	OutputVolume VolumeID
	Tasks        []TaskStatus
}

// JobState export
type JobState struct {
	Error  error
	Status string
}

// Task status exported
type TaskStatus struct {
	Name string
	Error error
	ExitCode int
	ConsoleOutput *bytes.Buffer
}

// JobID exported
type JobID string

func (p *Pier) SetState(job Job, state JobState) {
	// TODO: set this atomically

	job.State = &state
}

type jobArray []Job

func (jl jobArray) Len() int {
	return len(jl)
}
func (jl jobArray) Swap(i, j int) {
	jl[i], jl[j] = jl[j], jl[i]
}
func (jl jobArray) Less(i, j int) bool {
	return jl[i].ID < jl[j].ID
}

// Bind describes the binding between an IOPort and a docker volume
type Bind struct {
	IOPort   IOPort
	VolumeID dckr.VolumeID
}

// JobList is a shared structure that stores info about all jobs
type JobList struct {
	sync.Mutex
	cache map[JobID]Job
}

// NewJobList exported
func NewJobList() *JobList {
	return &JobList{
		cache: make(map[JobID]Job),
	}
}

func (jobList *JobList) add(job Job) {
	jobList.Lock()
	defer jobList.Unlock()
	jobList.cache[job.ID] = job
}

func (jobList *JobList) list() []Job {
	jobList.Lock()
	defer jobList.Unlock()
	all := make([]Job, len(jobList.cache), len(jobList.cache))
	i := 0
	for _, job := range jobList.cache {
		all[i] = job
		i++
	}
	sort.Sort(jobArray(all))
	return all
}

func (jobList *JobList) get(key JobID) (Job, bool) {
	jobList.Lock()
	defer jobList.Unlock()
	job, ok := jobList.cache[key]
	return job, ok
}

func (jobList *JobList) setState(jobID JobID, state JobState) {
	jobList.Lock()
	defer jobList.Unlock()
	job := jobList.cache[jobID]
	job.State = &state
	jobList.cache[jobID] = job
}

func (jobList *JobList) setInputVolume(jobID JobID, inputVolume VolumeID) {
	jobList.Lock()
	defer jobList.Unlock()
	job := jobList.cache[jobID]
	job.InputVolume = inputVolume
	jobList.cache[jobID] = job
}

func (jobList *JobList) setOutputVolume(jobID JobID, outputVolume VolumeID) {
	jobList.Lock()
	defer jobList.Unlock()
	job := jobList.cache[jobID]
	job.OutputVolume = outputVolume
	jobList.cache[jobID] = job
}

func (jobList *JobList) addTask(jobID JobID, taskName string, taskError error, taskExitCode int, taskConsoleOutput *bytes.Buffer) {
	jobList.Lock()
	defer jobList.Unlock()
	job := jobList.cache[jobID]

	var newTask TaskStatus
	newTask.Name = taskName
	newTask.Error = taskError
	newTask.ExitCode = taskExitCode
	newTask.ConsoleOutput = taskConsoleOutput
	job.Tasks = append(job.Tasks, newTask)

	jobList.cache[jobID] = job
}