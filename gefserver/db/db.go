package db

import (
	"bytes"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	// imported for side-effect only (package init)
	_ "github.com/mattn/go-sqlite3"

	"github.com/pborman/uuid"
	"gopkg.in/gorp.v1"
)

// Column in a table used to keep an internal version number of GORP
const gorpVersionColumn = "Revision"
const sqliteDataBasePath = "gef_db.bin"

// Db is used to keep DbMap
type Db struct {
	db gorp.DbMap
}

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
	ID             string
	Name           string
	ContainerID    string
	SwarmServiceID string
	Error          string
	ExitCode       int
	ConsoleOutput  string
	Revision       int
	JobID          string
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

// IOPortTable is used to store info about service inputs and outputs in a database
type IOPortTable struct {
	ID        string
	Name      string
	Path      string
	IsInput   bool
	Revision  int
	ServiceID string
}

// ServiceCmdTable stores CMD options for services
type ServiceCmdTable struct {
	ID        int
	Cmd       string
	Index     int
	Revision  int
	ServiceID string
}

// UserTable is used to store the users in a database
type UserTable struct {
	ID       int64
	Name     string
	Email    string
	Created  time.Time
	Updated  time.Time
	Revision int
}

// TokenTable used to store tokens in the db
type TokenTable struct {
	ID       int64
	Name     string // token name, user defined
	Secret   string // token secret, a random string
	UserID   int64
	Expire   time.Time
	Revision int
}

// InitDb initializes the database engine
func InitDb() (Db, error) {
	dataBase, err := sql.Open("sqlite3", sqliteDataBasePath)
	if err != nil {
		return Db{}, err
	}
	return initializeDatabase(dataBase)
}

// InitDbForTesting must only be used for tests
func InitDbForTesting() (Db, string, error) {
	filename := filepath.Join(os.TempDir(), "gef_db_test_"+uuid.New()+".bin")
	log.Println("new testing database file: ", filename)
	dataBase, err := sql.Open("sqlite3", filename)
	if err != nil {
		return Db{}, "", err
	}
	db, err := initializeDatabase(dataBase)
	return db, filename, err
}

func initializeDatabase(dataBase *sql.DB) (Db, error) {
	// For each table GORP has a special field to keep information about data version
	dataBaseMap := &gorp.DbMap{Db: dataBase, Dialect: gorp.SqliteDialect{}}

	dataBaseMap.AddTableWithName(JobTable{}, "Jobs").SetKeys(false, "ID").SetVersionCol(gorpVersionColumn)

	dataBaseMap.AddTableWithName(TaskTable{}, "Tasks").SetKeys(false, "ID").SetVersionCol(gorpVersionColumn)

	dataBaseMap.AddTableWithName(ServiceTable{}, "Services").SetKeys(false, "ID").SetVersionCol(gorpVersionColumn)

	dataBaseMap.AddTableWithName(IOPortTable{}, "IOPorts").SetVersionCol(gorpVersionColumn)

	dataBaseMap.AddTableWithName(ServiceCmdTable{}, "ServiceCmd").SetKeys(true, "ID").SetVersionCol(gorpVersionColumn)

	userTable := dataBaseMap.AddTableWithName(UserTable{}, "Users").SetKeys(true, "ID")
	{
		userTable.SetVersionCol(gorpVersionColumn)
		userTable.ColMap("Email").SetUnique(true)
	}

	tokensTable := dataBaseMap.AddTableWithName(TokenTable{}, "Tokens").SetKeys(true, "ID")
	{
		tokensTable.SetVersionCol(gorpVersionColumn)
		tokensTable.ColMap("Name").SetUnique(true)
		tokensTable.ColMap("Secret").SetUnique(true)
	}

	err := dataBaseMap.CreateTablesIfNotExists()
	return Db{db: *dataBaseMap}, err
}

// Close closes the db connections
func (d *Db) Close() {
	d.db.Db.Close()
}

// AddJob adds a job to the database
func (d *Db) AddJob(job Job) error {
	storedJob := d.job2JobTable(job)
	return d.db.Insert(&storedJob)
}

// RemoveJob removes a job and all corresponding tasks from the database
func (d *Db) RemoveJob(id JobID) error {
	_, err := d.db.Exec("DELETE FROM Tasks WHERE jobID=?", string(id))
	if err != nil {
		return err
	}

	_, err = d.db.Exec("DELETE FROM Jobs WHERE ID=?", string(id))
	return err
}

// RemoveJobTask removes a task from the database
func (d *Db) RemoveJobTask(taskID string) error {
	_, err := d.db.Exec("DELETE FROM Tasks WHERE ID=?", taskID)
	return err
}

// jobTable2Job performs mapping of the database job table to its JSON representation
func (d *Db) jobTable2Job(storedJob JobTable) (Job, error) {
	var job Job
	var jobState JobState
	var linkedTasks []Task
	var storedTasks []TaskTable
	_, err := d.db.Select(&storedTasks, "SELECT * FROM Tasks WHERE JobID=?", storedJob.ID)
	if err != nil {
		return job, err
	}

	for _, t := range storedTasks {
		var curTask Task
		curTask.Error = t.Error
		curTask.ConsoleOutput = t.ConsoleOutput
		curTask.ContainerID = ContainerID(t.ContainerID)
		curTask.SwarmServiceID = t.SwarmServiceID
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
	_, err := d.db.Select(&jobsFromTable, "SELECT * FROM Jobs ORDER BY ID")
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
	err := d.db.SelectOne(&jobFromTable, "SELECT * FROM jobs WHERE ID=?", string(id))
	if err != nil {
		return job, err
	}

	job, err = d.jobTable2Job(jobFromTable)
	return job, err
}

// SetJobState sets a job state
func (d *Db) SetJobState(id JobID, state JobState) error {
	var storedJob JobTable
	err := d.db.SelectOne(&storedJob, "SELECT * FROM jobs WHERE ID=?", string(id))
	if err != nil {
		return err
	}

	storedJob.Error = state.Error
	storedJob.Status = state.Status
	storedJob.Code = state.Code
	_, err = d.db.Update(&storedJob)
	return err
}

// SetJobInputVolume sets a job input volume
func (d *Db) SetJobInputVolume(id JobID, inputVolume VolumeID) error {
	var storedJob JobTable
	err := d.db.SelectOne(&storedJob, "SELECT * FROM jobs WHERE ID=?", string(id))
	if err != nil {
		return err
	}

	storedJob.InputVolume = string(inputVolume)
	_, err = d.db.Update(&storedJob)
	return err
}

// SetJobOutputVolume sets a job output volume
func (d *Db) SetJobOutputVolume(id JobID, outputVolume VolumeID) error {
	var storedJob JobTable
	err := d.db.SelectOne(&storedJob, "SELECT * from jobs WHERE ID=?", string(id))
	if err != nil {
		return err
	}

	storedJob.OutputVolume = string(outputVolume)
	_, err = d.db.Update(&storedJob)
	return err
}

// AddJobTask adds a task to a job
func (d *Db) AddJobTask(id JobID, taskName string, taskContainer string, taskSwarmService string,
	taskError string, taskExitCode int, taskConsoleOutput *bytes.Buffer) error {
	var newTask TaskTable
	newTask.ID = uuid.New()
	newTask.Name = taskName
	newTask.ContainerID = string(taskContainer)
	newTask.SwarmServiceID = taskSwarmService
	newTask.Error = taskError
	newTask.ExitCode = taskExitCode
	newTask.ConsoleOutput = taskConsoleOutput.String()
	newTask.JobID = string(id)
	return d.db.Insert(&newTask)
}

// serviceTable2Service performs mapping of the database service table to its JSON representation
func (d *Db) serviceTable2Service(storedService ServiceTable) (Service, error) {
	var service Service
	var storedInputPorts []IOPortTable
	var storedOutputPorts []IOPortTable
	var inputPorts []IOPort
	var outputPorts []IOPort
	var selectedCmd []ServiceCmdTable

	_, err := d.db.Select(&storedInputPorts, "SELECT * FROM IOPorts WHERE IsInput=1 AND ServiceID=?",
		storedService.ID)
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

	_, err = d.db.Select(&storedOutputPorts, "SELECT * FROM IOPorts WHERE IsInput=0 AND ServiceID=?",
		storedService.ID)
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

	_, err = d.db.Select(&selectedCmd, "SELECT * FROM ServiceCmd WHERE ServiceID=?", storedService.ID)
	if err != nil {
		return service, err
	}

	var storedCmd []string
	// arguments should have a correct order
	for i, s := range selectedCmd {
		for _, item := range selectedCmd {
			if item.Index == i {
				storedCmd = append(storedCmd, s.Cmd)
				break
			}
		}
	}

	service.ID = ServiceID(storedService.ID)
	service.ImageID = ImageID(storedService.ImageID)
	service.Name = storedService.Name
	service.RepoTag = storedService.RepoTag
	service.Description = storedService.Description
	service.Version = storedService.Version
	service.Cmd = storedCmd
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
		err = d.db.Insert(&curInputPort)
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
		err = d.db.Insert(&curOutputPort)
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
	_, err := d.db.Select(&servicesFromTable, "SELECT * FROM services WHERE Name=?", service.Name)
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

	err = d.AddCmd(service)
	if err != nil {
		return err
	}

	storedService := d.service2ServiceTable(service)
	err = d.db.Insert(&storedService)
	return err
}

// AddCmd adds an array of cmd options to the specified service
func (d *Db) AddCmd(service Service) error {
	var storedCmd ServiceCmdTable
	for i, c := range service.Cmd {
		storedCmd.Cmd = c
		storedCmd.Index = i
		storedCmd.ServiceID = string(service.ID)
		err := d.db.Insert(&storedCmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveService removes a service and the corresponding IOPorts from the database
func (d *Db) RemoveService(id ServiceID) error {
	// remove linked ports
	_, err := d.db.Exec("DELETE FROM IOPorts WHERE ServiceID=?", string(id))
	if err != nil {
		return err
	}

	// remove linked cmd
	_, err = d.db.Exec("DELETE FROM ServiceCmd WHERE ServiceID=?", string(id))
	if err != nil {
		return err
	}

	// remove the service itself
	_, err = d.db.Exec("DELETE FROM services WHERE ID=?", string(id))
	return err
}

// ListServices produces a list of all services ready to be converted into JSON
func (d *Db) ListServices() ([]Service, error) {
	var services []Service
	var servicesFromTable []ServiceTable
	_, err := d.db.Select(&servicesFromTable, "SELECT * FROM services ORDER BY ID")
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
	err := d.db.SelectOne(&serviceFromTable, "SELECT * FROM services WHERE ID=?", string(id))
	if err != nil {
		return service, err
	}

	service, err = d.serviceTable2Service(serviceFromTable)
	return service, err
}
