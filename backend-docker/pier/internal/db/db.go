package db

import (
	"database/sql"
	"gopkg.in/gorp.v1"
	_ "github.com/mattn/go-sqlite3"

	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
	"bytes"
	"time"
	"sort"
)

type VolumeID dckr.VolumeID

type Db struct {m gorp.DbMap}

// Job stores the information about a service execution
type Job struct {
	ID           JobID
	ServiceID    ServiceID
	Input        string
	Created      time.Time
	State        *JobState
	InputVolume  VolumeID
	OutputVolume VolumeID
	Tasks        []TaskInfo
}

// JobState exported
type JobState struct {
	Error  error
	Status string
	Code   int
}

// JobID exported
type JobID string

type jobArray []Job

// TaskInfo exported
type TaskInfo struct {
	Name          string
	ContainerID   dckr.ContainerID
	Error         error
	ExitCode      int
	ConsoleOutput *bytes.Buffer
}

// LatestOutput used to serialize consoleoutput to json
type LatestOutput struct {
	Name          string
	ConsoleOutput string
}

// Bind describes the binding between an IOPort and a docker volume
type Bind struct {
	IOPort   IOPort
	VolumeID dckr.VolumeID
}

// GefSrvLabelPrefix is the prefix identifying GEF related labels
const GefSrvLabelPrefix = "eudat.gef.service."

// Service describes metadata for a GEF service
type Service struct {
	ID          ServiceID
	imageID     dckr.ImageID
	Name        string
	RepoTag     string
	Description string
	Version     string
	Created     time.Time
	Size        int64
	Input       []IOPort
	Output      []IOPort
}

// ServiceID exported
type ServiceID string

// IOPort is an i/o specification for a service
// The service can only read data from volumes and write to a single volume
// Path specifies where the volumes are mounted
type IOPort struct {
	ID   string
	Name string
	Path string
}

// InitDb exported
func InitDb() (Db, error) {
	dataBase, err := sql.Open("sqlite3", "/Users/achernov/job_db.bin")

	/*if err != nil {
		return nil, err
	}*/

	dataBaseMap := &gorp.DbMap{Db: dataBase, Dialect: gorp.SqliteDialect{}}
	dataBaseMap.AddTableWithName(Job{}, "jobs").SetKeys(true, "ID")
	dataBaseMap.AddTableWithName(Job{}, "services").SetKeys(true, "ID")
	err = dataBaseMap.CreateTablesIfNotExists()

	dbm := Db {m: *dataBaseMap}

	return dbm, err
}


// Jobs

func (d *Db) AddJob(job Job) error {
	return d.m.Insert(&job)
}

func (d *Db) RemoveJob(jobID JobID) error {
	_, err := d.m.Exec("delete from jobs where ID=?", jobID)
	return err
}

// ListJobs exported
func (d *Db) ListJobs() ([]Job, error) {
	var jobs []Job
	_, err := d.m.Select(&jobs, "select * from jobs order by ID")
	return jobs, err
}

func (d *Db) GetJob(jobID JobID) (Job, error) {
	var job Job
	err := d.m.SelectOne(&job, "select * from jobs where ID=?", jobID)
	return job, err
}

func (d *Db) SetJobState(jobID JobID, state JobState) error {
	var job Job
	err := d.m.SelectOne(&job, "select * from jobs where ID=?", jobID)
	if err != nil {
		job.State = &state
		_, err = d.m.Update(&job)
	}
	return err
}

func (d *Db) SetJobInputVolume(jobID JobID, inputVolume VolumeID) error {
	var job Job
	err := d.m.SelectOne(&job, "select * from jobs where ID=?", jobID)
	if err != nil {
		job.InputVolume = inputVolume
		_, err = d.m.Update(&job)
	}
	return err
}

func (d *Db) SetJobOutputVolume(jobID JobID, outputVolume VolumeID) error {
	var job Job
	err := d.m.SelectOne(&job, "select * from jobs where ID=?", jobID)
	if err != nil {
		job.OutputVolume = outputVolume
		_, err = d.m.Update(&job)
	}
	return err
}

func (d *Db) AddJobTask(jobID JobID, taskName string, taskContainer dckr.ContainerID, taskError error, taskExitCode int, taskConsoleOutput *bytes.Buffer) error {
	var job Job
	err := d.m.SelectOne(&job, "select * from jobs where ID=?", jobID)
	if err != nil {
		var newTaskInfo TaskInfo
		newTaskInfo.Name = taskName
		newTaskInfo.ContainerID = taskContainer
		newTaskInfo.Error = taskError
		newTaskInfo.ExitCode = taskExitCode
		newTaskInfo.ConsoleOutput = taskConsoleOutput
		job.Tasks = append(job.Tasks, newTaskInfo)
		_, err = d.m.Update(&job)
	}
	return err
}

// Services

func (d *Db) AddService(service Service) error {
	return d.m.Insert(&service)
}

func (d *Db) ListServices() ([]Service, error) {
	var services []Service
	_, err := d.m.Select(&services, "select * from services order by ID")
	return services, err
}

func (d *Db) GetService(serviceID ServiceID) (Service, error) {
	var service Service
	err := d.m.SelectOne(&service, "select * from service where ID=?", serviceID)
	return service, err
}