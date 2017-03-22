package db

import (
	"bytes"
	"database/sql"
	"time"

	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pborman/uuid"
	"gopkg.in/gorp.v1"
)

const GefSrvLabelPrefix = "eudat.gef.service." // GefSrvLabelPrefix is the prefix identifying GEF related labels
const GorpVersionColumn = "Revision"           // Column in a table used to keep an internal version number of GORP
const SQLiteDataBasePath = "gef_db.bin"

// Db is used to keep DbMap
type Db struct{ gorp.DbMap }

// JobTable stores the information about a service execution (used to store data in a database)
type JobTable struct {
	ID           string
	ServiceID    string
	Input        string
	Created      time.Time
	Error        string
	Status       string
	Code         int
	InputVolume  string
	OutputVolume string
	Revision     int
}

// TaskTable contains tasks related to a specific job (used to store data in a database)
type TaskTable struct {
	ID            string
	Name          string
	ContainerID   string
	Error         string
	ExitCode      int
	ConsoleOutput string
	Revision      int
	JobID         string
}

// ServiceTable describes metadata for a GEF service (used to store data in a database)
type ServiceTable struct {
	ID          string
	ImageID     string
	Name        string
	RepoTag     string
	Description string
	Version     string
	Revision    int
	Created     time.Time
	Size        int64
}

// IOPortTable is used to store data in a database
type IOPortTable struct {
	ID        string
	Name      string
	Path      string
	IsInput   bool
	Revision  int
	ServiceID string
}

// InitDb initializes the database engine
func InitDb() (Db, error) {
	dataBase, err := sql.Open("sqlite3", SQLiteDataBasePath)
	if err != nil {
		return Db{}, err
	}
	// For each table GORP has a special field to keep information about data version
	dataBaseMap := &gorp.DbMap{Db: dataBase, Dialect: gorp.SqliteDialect{}}
	dataBaseMap.AddTableWithName(JobTable{}, "Jobs").SetKeys(false, "ID").SetVersionCol(GorpVersionColumn)
	dataBaseMap.AddTableWithName(TaskTable{}, "Tasks").SetKeys(false, "ID").SetVersionCol(GorpVersionColumn)
	dataBaseMap.AddTableWithName(ServiceTable{}, "Services").SetKeys(false, "ID").SetVersionCol(GorpVersionColumn)
	dataBaseMap.AddTableWithName(IOPortTable{}, "IOPorts").SetVersionCol(GorpVersionColumn)
	err = dataBaseMap.CreateTablesIfNotExists()
	return Db{*dataBaseMap}, err
}

// AddJob adds a job to the database
func (d *Db) AddJob(job Job) error {
	storedJob := d.job2JobTable(job)
	return d.Insert(&storedJob)
}

// RemoveJob removes a job and all corresponding tasks from the database
func (d *Db) RemoveJob(id JobID) error {
	_, err := d.Exec("DELETE FROM Tasks WHERE jobID=?", string(id))
	if err != nil {
		return err
	}

	_, err = d.Exec("DELETE FROM Jobs WHERE ID=?", string(id))
	return err
}

// RemoveJobTask removes a task from the database
func (d *Db) RemoveJobTask(taskID string) error {
	_, err := d.Exec("DELETE FROM Tasks WHERE ID=?", taskID)
	return err
}

// jobTable2Job performs mapping of the database job table to its JSON representation
func (d *Db) jobTable2Job(storedJob JobTable) (Job, error) {
	var job Job
	var jobState JobState
	var linkedTasks []Task
	var storedTasks []TaskTable
	_, err := d.Select(&storedTasks, "SELECT * FROM Tasks WHERE JobID=?", storedJob.ID)
	if err != nil {
		return job, err
	}

	for _, t := range storedTasks {
		var curTask Task
		curTask.Error = t.Error
		curTask.ConsoleOutput = t.ConsoleOutput
		curTask.ContainerID = dckr.ContainerID(t.ContainerID)
		curTask.ExitCode = t.ExitCode
		curTask.ID = t.ID
		curTask.Name = t.Name

		linkedTasks = append(linkedTasks, curTask)
	}

	jobState.Error = storedJob.Error
	jobState.Status = storedJob.Status
	jobState.Code = storedJob.Code

	job.ID = JobID(storedJob.ID)
	job.ServiceID = ServiceID(storedJob.ServiceID)
	job.Input = storedJob.Input
	job.Created = storedJob.Created
	job.State = &jobState
	job.InputVolume = VolumeID(storedJob.InputVolume)
	job.OutputVolume = VolumeID(storedJob.OutputVolume)
	job.Tasks = linkedTasks
	return job, err
}

// job2JobTable performs mapping of the job JSON representation to its database representation
func (d *Db) job2JobTable(job Job) JobTable {
	var storedJob JobTable
	storedJob.ID = string(job.ID)
	storedJob.ServiceID = string(job.ServiceID)
	storedJob.Input = job.Input
	storedJob.Created = job.Created
	storedJob.Error = job.State.Error
	storedJob.Status = job.State.Status
	storedJob.Code = job.State.Code
	storedJob.InputVolume = string(job.InputVolume)
	storedJob.OutputVolume = string(job.OutputVolume)
	return storedJob
}

// ListJobs returns a list of all jobs ready to be converted into JSON
func (d *Db) ListJobs() ([]Job, error) {
	var jobs []Job
	var jobsFromTable []JobTable
	_, err := d.Select(&jobsFromTable, "SELECT * FROM Jobs ORDER BY ID")
	if err != nil {
		return jobs, err
	}

	for _, j := range jobsFromTable {
		var curJob Job
		curJob, err = d.jobTable2Job(j)
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, curJob)
	}

	return jobs, err
}

// GetJob returns a JSON ready representation of a job
func (d *Db) GetJob(id JobID) (Job, error) {
	var job Job
	var jobFromTable JobTable
	err := d.SelectOne(&jobFromTable, "SELECT * FROM jobs WHERE ID=?", string(id))
	if err != nil {
		return job, err
	}

	job, err = d.jobTable2Job(jobFromTable)
	return job, err
}

// SetJobState sets a job state
func (d *Db) SetJobState(id JobID, state JobState) error {
	var storedJob JobTable
	err := d.SelectOne(&storedJob, "SELECT * FROM jobs WHERE ID=?", string(id))
	if err != nil {
		return err
	}

	storedJob.Error = state.Error
	storedJob.Status = state.Status
	storedJob.Code = state.Code
	_, err = d.Update(&storedJob)
	return err
}

// SetJobInputVolume sets a job input volume
func (d *Db) SetJobInputVolume(id JobID, inputVolume VolumeID) error {
	var storedJob JobTable
	err := d.SelectOne(&storedJob, "SELECT * FROM jobs WHERE ID=?", string(id))
	if err != nil {
		return err
	}

	storedJob.InputVolume = string(inputVolume)
	_, err = d.Update(&storedJob)
	return err
}

// SetJobOutputVolume sets a job output volume
func (d *Db) SetJobOutputVolume(id JobID, outputVolume VolumeID) error {
	var storedJob JobTable
	err := d.SelectOne(&storedJob, "SELECT * from jobs WHERE ID=?", string(id))
	if err != nil {
		return err
	}

	storedJob.OutputVolume = string(outputVolume)
	_, err = d.Update(&storedJob)
	return err
}

// AddJobTask adds a task to a job
func (d *Db) AddJobTask(id JobID, taskName string, taskContainer dckr.ContainerID, taskError string, taskExitCode int, taskConsoleOutput *bytes.Buffer) error {
	var newTask TaskTable
	newTask.ID = uuid.New()
	newTask.Name = taskName
	newTask.ContainerID = string(taskContainer)
	newTask.Error = taskError
	newTask.ExitCode = taskExitCode
	newTask.ConsoleOutput = taskConsoleOutput.String()
	newTask.JobID = string(id)
	return d.Insert(&newTask)
}

// serviceTable2Service performs mapping of the database service table to its JSON representation
func (d *Db) serviceTable2Service(storedService ServiceTable) (Service, error) {
	var service Service
	var storedInputPorts []IOPortTable
	var storedOutputPorts []IOPortTable
	var inputPorts []IOPort
	var outputPorts []IOPort

	_, err := d.Select(&storedInputPorts, "SELECT * FROM IOPorts WHERE IsInput=1 AND ServiceID=?", storedService.ID)
	if err != nil {
		return service, err
	}

	for _, i := range storedInputPorts {
		var curInput IOPort
		curInput.ID = i.ID
		curInput.Name = i.Name
		curInput.Path = i.Path

		inputPorts = append(inputPorts, curInput)
	}

	_, err = d.Select(&storedOutputPorts, "SELECT * FROM IOPorts WHERE IsInput=0 AND ServiceID=?", storedService.ID)
	if err != nil {
		return service, err
	}

	for _, o := range storedOutputPorts {
		var curOutput IOPort
		curOutput.ID = o.ID
		curOutput.Name = o.Name
		curOutput.Path = o.Path

		outputPorts = append(outputPorts, curOutput)
	}

	service.ID = ServiceID(storedService.ID)
	service.ImageID = dckr.ImageID(storedService.ImageID)
	service.Name = storedService.Name
	service.RepoTag = storedService.RepoTag
	service.Description = storedService.Description
	service.Version = storedService.Version
	service.Created = storedService.Created
	service.Size = storedService.Size
	service.Input = inputPorts
	service.Input = inputPorts
	service.Output = outputPorts

	return service, err
}

// service2ServiceTable performs mapping of the service JSON representation to its database representation
func (d *Db) service2ServiceTable(service Service) ServiceTable {
	var storedService ServiceTable
	storedService.ID = string(service.ID)
	storedService.ImageID = string(service.ImageID)
	storedService.Name = service.Name
	storedService.RepoTag = service.RepoTag
	storedService.Description = service.Description
	storedService.Version = service.Version
	storedService.Created = service.Created
	storedService.Size = service.Size
	return storedService
}

// AddIOPort adds input and output ports to the database
func (d *Db) AddIOPort(service Service) error {
	var err error
	for _, p := range service.Input {
		var curInputPort IOPortTable
		curInputPort.Path = p.Path
		curInputPort.Name = p.Name
		curInputPort.ID = p.ID
		curInputPort.IsInput = true
		curInputPort.ServiceID = string(service.ID)
		err = d.Insert(&curInputPort)
		if err != nil {
			return err
		}
	}

	for _, p := range service.Output {
		var curOutputPort IOPortTable
		curOutputPort.Path = p.Path
		curOutputPort.Name = p.Name
		curOutputPort.ID = p.ID
		curOutputPort.IsInput = false
		curOutputPort.ServiceID = string(service.ID)
		err = d.Insert(&curOutputPort)
		if err != nil {
			return err
		}
	}
	return err
}

// AddService creates a new service in the database
func (d *Db) AddService(service Service) error {
	// Before adding a service we need to check if the service with the same name already exists.
	// If it does, we remove it and add a new one
	var servicesFromTable []ServiceTable
	_, err := d.Select(&servicesFromTable, "SELECT * FROM services WHERE Name=?", service.Name)
	if err != nil {
		return err
	}

	if len(servicesFromTable) > 0 {
		for _, s := range servicesFromTable {
			err = d.RemoveService(ServiceID(s.ID))
			if err != nil {
				return err
			}
		}
	}

	err = d.AddIOPort(service)
	if err != nil {
		return err
	}

	storedService := d.service2ServiceTable(service)
	err = d.Insert(&storedService)
	return err
}

// RemoveService removes a service and the corresponding IOPorts from the database
func (d *Db) RemoveService(id ServiceID) error {
	_, err := d.Exec("DELETE FROM IOPorts WHERE ServiceID=?", string(id))
	if err != nil {
		return err
	}

	_, err = d.Exec("DELETE FROM services WHERE ID=?", string(id))
	return err
}

// ListServices produces a list of all services ready to be converted into JSON
func (d *Db) ListServices() ([]Service, error) {
	var services []Service
	var servicesFromTable []ServiceTable
	_, err := d.Select(&servicesFromTable, "SELECT * FROM services ORDER BY ID")
	if err != nil {
		return services, err
	}

	for _, s := range servicesFromTable {
		var curService Service
		curService, err = d.serviceTable2Service(s)
		if err != nil {
			return services, err
		}
		services = append(services, curService)
	}
	return services, err
}

// GetService returns a service ready to be converted into JSON
func (d *Db) GetService(id ServiceID) (Service, error) {
	var service Service
	var serviceFromTable ServiceTable
	err := d.SelectOne(&serviceFromTable, "SELECT * FROM services WHERE ID=?", string(id))
	if err != nil {
		return service, err
	}

	service, err = d.serviceTable2Service(serviceFromTable)
	return service, err
}