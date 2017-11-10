package db

import (
	"bytes"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	// imported for side-effect only (package init)
	_ "github.com/mattn/go-sqlite3"
	"github.com/pborman/uuid"
	gorp "gopkg.in/gorp.v1"

	"github.com/EUDAT-GEF/GEF/gefserver/def"
)

// Column in a table used to keep an internal version number of GORP
const gorpVersionColumn = "Revision"
const sqliteDataBasePath = "gef_db.bin"

// Db is used to keep DbMap
type Db struct {
	db gorp.DbMap
}

// ConnectionID is the type used to identify a docker connection
type ConnectionID int

// ConnectionTable stores the information about a docker connection (used to store data in a database)
type ConnectionTable struct {
	ID          int
	Endpoint    string // unique key
	Description string
	TLSVerify   bool
	CertPath    string
	KeyPath     string
	CAPath      string
	Revision    int
}

// JobTable stores the information about a service execution (used to store data in a database)
type JobTable struct {
	ID           string
	ConnectionID int
	ServiceID    string
	Created      time.Time
	Duration     int64 // duration time in seconds
	Error        string
	Status       string
	Code         int
	Revision     int
}

// VolumeTable contains information about input and output volumes for jobs
type VolumeTable struct {
	ID         string
	IsInput    bool
	JobID      string
	IOPortName string
	Content    string
	Revision   int
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
	JobID          string
	Revision       int
}

// ServiceTable describes metadata for a GEF service (used to store data in a database)
type ServiceTable struct {
	ID           string
	ConnectionID int
	ImageID      string
	Name         string
	RepoTag      string
	Description  string
	Version      string
	Created      time.Time
	Deleted      bool
	Size         int64
	Revision     int
}

// IOPortTable is used to store info about service inputs and outputs in a database
type IOPortTable struct {
	ID        string
	Name      string
	Path      string
	IsInput   bool
	ServiceID string
	Type      string
	FileName  string
	Revision  int
}

// ServiceCmdTable stores CMD options for services
type ServiceCmdTable struct {
	ID        int
	Cmd       string
	Index     int
	ServiceID string
	Revision  int
}

// UserTable stores the users in the db
type UserTable struct {
	ID       int64
	Name     string
	Email    string
	Created  time.Time
	Updated  time.Time
	Revision int
}

// TokenTable stores user tokens in the db
type TokenTable struct {
	ID       int64
	Name     string // token name, user defined
	Secret   string // token secret, a random string
	UserID   int64
	Expire   time.Time
	Revision int
}

//CommunityTable stores the communities in the db
type CommunityTable struct {
	ID          int64
	Name        string
	Description string
	Revision    int
}

// RoleTable stores user roles in the db
type RoleTable struct {
	ID          int64
	Name        string
	CommunityID int64 // most roles are per community
	Description string
	Revision    int
}

// UserRoleTable stores user mapping to roles in the db
type UserRoleTable struct {
	UserID   int64
	RoleID   int64
	Revision int
}

// OwnerTable stores object ownerships
type OwnerTable struct {
	UserID     int64
	ObjectType string
	ObjectID   string
	Revision   int
}

// InitDb initializes the database engine
func InitDb() (Db, error) {
	dataBase, err := sql.Open("sqlite3", sqliteDataBasePath)
	if err != nil {
		return Db{}, err
	}
	return setupDatabase(dataBase)
}

// InitDbForTesting must only be used for tests
func InitDbForTesting() (Db, string, error) {
	filename := filepath.Join(os.TempDir(), "gef_db_test_"+uuid.New()+".bin")
	log.Println("new testing database file: ", filename)
	dataBase, err := sql.Open("sqlite3", filename)
	if err != nil {
		return Db{}, "", err
	}
	db, err := setupDatabase(dataBase)
	if err != nil {
		err = def.Err(err, "error in setupDatabase")
	}
	return db, filename, err
}

func setupDatabase(dataBase *sql.DB) (Db, error) {
	// For each table GORP has a special field to keep information about data version
	dataBaseMap := &gorp.DbMap{Db: dataBase, Dialect: gorp.SqliteDialect{}}

	connectionTable := dataBaseMap.AddTableWithName(ConnectionTable{}, "Connections").SetKeys(true, "ID")
	{
		connectionTable.SetVersionCol(gorpVersionColumn)
		connectionTable.ColMap("Endpoint").SetUnique(true)
	}

	dataBaseMap.AddTableWithName(JobTable{}, "Jobs").SetKeys(false, "ID").SetVersionCol(gorpVersionColumn)

	dataBaseMap.AddTableWithName(VolumeTable{}, "Volumes").SetKeys(false, "ID").SetVersionCol(gorpVersionColumn)

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
		tokensTable.ColMap("Secret").SetUnique(true)
	}

	communityTable := dataBaseMap.AddTableWithName(CommunityTable{}, "Communities").SetKeys(true, "ID")
	{
		communityTable.SetVersionCol(gorpVersionColumn)
		communityTable.ColMap("Name").SetUnique(true)
	}

	rolesTable := dataBaseMap.AddTableWithName(RoleTable{}, "Roles").SetKeys(true, "ID")
	{
		rolesTable.SetVersionCol(gorpVersionColumn)
	}

	dataBaseMap.AddTableWithName(UserRoleTable{}, "UserRoles").SetVersionCol(gorpVersionColumn)
	dataBaseMap.AddTableWithName(OwnerTable{}, "Owners").SetVersionCol(gorpVersionColumn)

	err := dataBaseMap.CreateTablesIfNotExists()
	if err != nil {
		return Db{}, err
	}

	db := Db{db: *dataBaseMap}
	err = initializeDatabaseValues(db)
	if err != nil {
		err = def.Err(err, "error in initializeDatabaseValues")
	}
	return db, err
}

func initializeDatabaseValues(d Db) error {
	_, err := d.AddRole(SuperAdminRoleName, 0, "Super Administrator of the site, with all privileges.")
	if err != nil {
		return def.Err(err, "error in AddRole:SuperAdmin")
	}

	description := "The EUDAT community. Use this community if no other is suited for you."
	_, err = d.AddCommunity("EUDAT", description, true)
	if err != nil {
		return def.Err(err, "error in AddCommunity:EUDAT")
	}
	return nil
}

func isNoResultsError(e error) bool {
	if e == nil {
		return false
	}
	return strings.HasSuffix(e.Error(), "no rows in result set")
}

// Close closes the db connections
func (d *Db) Close() {
	d.db.Db.Close()
}

// AddConnection adds a connection to the database
func (d *Db) AddConnection(userID int64, connection def.DockerConfig) (ConnectionID, error) {
	var ct ConnectionTable
	err := d.db.SelectOne(&ct,
		"SELECT * FROM connections WHERE Endpoint=?",
		connection.Endpoint)

	if err != nil && !isNoResultsError(err) {
		return 0, def.Err(err, "db inquiry about docker connections failed: %v", err)
	}

	ct.Endpoint = connection.Endpoint
	ct.Description = connection.Description
	ct.TLSVerify = connection.TLSVerify
	ct.CertPath = connection.CertPath
	ct.KeyPath = connection.KeyPath
	ct.CAPath = connection.CAPath

	if isNoResultsError(err) {
		err = d.db.Insert(&ct)
		if err != nil {
			return 0, def.Err(err, "db inserting docker connection failed: %v", err)
		}

		ownership := OwnerTable{
			UserID:     userID,
			ObjectType: "Connection",
			ObjectID:   string(ct.ID),
		}
		err = d.db.Insert(&ownership)
		if err != nil {
			return 0, def.Err(err, "db inserting connection ownership failed")
		}
	} else {
		d.db.Update(&ct)
	}

	return ConnectionID(ct.ID), nil
}

// RemoveConnection removes a connection to the database
func (d *Db) RemoveConnection(connectionID ConnectionID) error {
	_, err := d.db.Exec("DELETE FROM connections WHERE ID=?", connectionID)
	if err != nil {
		return err
	}

	_, err = d.db.Exec("DELETE FROM Owners WHERE ObjectType=? AND ObjectID=?",
		"Connection", string(connectionID))
	if err != nil {
		return err
	}
	return nil
}

// GetConnections returns a map of all connections ready to be converted into JSON
func (d *Db) GetConnections() (map[ConnectionID]def.DockerConfig, error) {
	var connectionTable []ConnectionTable
	_, err := d.db.Select(&connectionTable, "SELECT * FROM connections ORDER BY ID")
	if err != nil {
		return nil, err
	}

	connections := make(map[ConnectionID]def.DockerConfig)
	for _, c := range connectionTable {
		connections[ConnectionID(c.ID)] = def.DockerConfig{
			Endpoint:    c.Endpoint,
			Description: c.Description,
			TLSVerify:   c.TLSVerify,
			CertPath:    c.CertPath,
			KeyPath:     c.KeyPath,
			CAPath:      c.CAPath,
		}
	}

	return connections, err
}

// GetFirstConnectionID returns the default (first) connection id
func (d *Db) GetFirstConnectionID() (ConnectionID, error) {
	var connectionIDs []int
	_, err := d.db.Select(&connectionIDs, "SELECT ID FROM connections ORDER BY ID LIMIT 1")
	if err != nil {
		return 0, err
	}
	if len(connectionIDs) < 1 {
		return 0, def.Err(nil, "Cannot find any connection id in the database")
	}

	return ConnectionID(connectionIDs[0]), nil
}

// GetConnectionOwners does what it says
func (d *Db) GetConnectionOwners(connectionID ConnectionID) ([]int64, error) {
	var ownersTable []OwnerTable
	_, err := d.db.Select(&ownersTable,
		"SELECT * FROM owners WHERE ObjectType=? AND ObjectID=?",
		"Connection", string(connectionID))
	if err != nil && !isNoResultsError(err) {
		log.Printf("ERROR in GetConnectionOwners: %#v", err)
	}
	owners := make([]int64, 0, len(ownersTable))
	for _, o := range ownersTable {
		owners = append(owners, o.UserID)
	}
	return owners, err
}

// IsConnectionOwner checks if a certain user owns a certain connection
func (d *Db) IsConnectionOwner(userID int64, connectionID ConnectionID) bool {
	var x OwnerTable
	err := d.db.SelectOne(&x,
		"SELECT * FROM owners WHERE UserID=? AND ObjectType=? AND ObjectID=?",
		userID, "Connection", string(connectionID))
	if err != nil && !isNoResultsError(err) {
		log.Printf("ERROR in IsConnectionOwner: %#v", err)
	}
	return err == nil
}

// AddJob adds a job to the database
func (d *Db) AddJob(userID int64, job Job) error {
	storedJob := d.job2JobTable(job)
	err := d.db.Insert(&storedJob)
	if err != nil {
		return err
	}
	ownership := OwnerTable{
		UserID:     userID,
		ObjectType: "Job",
		ObjectID:   string(job.ID),
	}
	return d.db.Insert(&ownership)
}

// RemoveJob removes a job and all corresponding tasks from the database
func (d *Db) RemoveJob(id JobID) error {
	_, err := d.db.Exec("DELETE FROM Tasks WHERE jobID=?", string(id))
	if err != nil {
		return err
	}

	_, err = d.db.Exec("DELETE FROM Jobs WHERE ID=?", string(id))
	if err != nil {
		return err
	}

	_, err = d.db.Exec("DELETE FROM Owners WHERE ObjectType=? AND ObjectID=?",
		"Job", string(id))
	if err != nil {
		return err
	}
	return nil
}

// RemoveJobTask removes a task from the database
func (d *Db) RemoveJobTask(taskID string) error {
	_, err := d.db.Exec("DELETE FROM Tasks WHERE ID=?", taskID)
	return err
}

// CountRunningJobs returns a number of jobs currently running
func (d *Db) CountRunningJobs() int64 {
	count, err := d.db.SelectInt(
		"SELECT count(*) FROM jobs WHERE jobs.Code<0")

	if err != nil && !isNoResultsError(err) {
		log.Printf("ERROR in CountRunningJobs: %#v", err)
	}
	return count
}

// CountUserRunningJobs returns a number of running jobs owned by a specific user
func (d *Db) CountUserRunningJobs(userID int64) int64 {
	count, err := d.db.SelectInt(
		"SELECT count(*) FROM owners INNER JOIN jobs on owners.ObjectID = jobs.ID WHERE owners.UserID=? AND jobs.Code<0", userID)

	if err != nil { //&& !isNoResultsError(err) {
		log.Printf("ERROR in CountUserRunningJobs: %#v", err)
	}
	return count
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

	var storedVolumes []VolumeTable
	_, err = d.db.Select(&storedVolumes, "SELECT * FROM Volumes WHERE JobID=?", storedJob.ID)
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

	var inputVolumes []JobVolume
	var outputVolumes []JobVolume
	for s := range storedVolumes {
		var curJobVolume JobVolume
		curJobVolume.Name = storedVolumes[s].IOPortName
		curJobVolume.VolumeID = VolumeID(storedVolumes[s].ID)
		if storedVolumes[s].IsInput {
			inputVolumes = append(inputVolumes, curJobVolume)
		} else {
			outputVolumes = append(outputVolumes, curJobVolume)
		}
	}

	jobState.Error = storedJob.Error
	jobState.Status = storedJob.Status
	jobState.Code = storedJob.Code

	job.ID = JobID(storedJob.ID)
	job.ConnectionID = ConnectionID(storedJob.ConnectionID)
	job.ServiceID = ServiceID(storedJob.ServiceID)
	job.Created = storedJob.Created

	if jobState.Code < 0 {
		job.Duration = time.Now().Unix() - job.Created.Unix()
	} else {
		job.Duration = storedJob.Duration
	}

	job.State = &jobState
	job.InputVolume = inputVolumes
	job.OutputVolume = outputVolumes
	job.Tasks = linkedTasks
	return job, err
}

// job2JobTable performs mapping of the job JSON representation to its database representation
func (d *Db) job2JobTable(job Job) JobTable {
	var storedJob JobTable
	storedJob.ID = string(job.ID)
	storedJob.ConnectionID = int(job.ConnectionID)
	storedJob.ServiceID = string(job.ServiceID)
	storedJob.Created = job.Created
	storedJob.Duration = job.Duration
	storedJob.Error = job.State.Error
	storedJob.Status = job.State.Status
	storedJob.Code = job.State.Code
	return storedJob
}

func (d *Db) updateJobDurationTime(job Job) {
	err := d.SetJobDurationTime(job.ID, time.Now().Unix()-job.Created.Unix())
	if err != nil {
		log.Println(err)
	}
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

// AddJobVolume sets a job input/output volume
func (d *Db) AddJobVolume(id JobID, volume VolumeID, isInput bool, portName string, content string) error {
	var storedVolumes VolumeTable
	storedVolumes.ID = string(volume)
	storedVolumes.JobID = string(id)
	storedVolumes.IsInput = isInput
	storedVolumes.IOPortName = portName
	storedVolumes.Content = content
	return d.db.Insert(&storedVolumes)
}

// SetJobDurationTime sets job finish time
func (d *Db) SetJobDurationTime(id JobID, duration int64) error {
	var storedJob JobTable
	err := d.db.SelectOne(&storedJob, "SELECT * from jobs WHERE ID=?", string(id))
	if err != nil {
		return err
	}

	storedJob.Duration = duration
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
		curInput.Type = i.Type
		curInput.FileName = i.FileName

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
		curOutput.Type = o.Type
		curOutput.FileName = o.FileName

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
	service.ConnectionID = ConnectionID(storedService.ConnectionID)
	service.ImageID = ImageID(storedService.ImageID)
	service.Name = storedService.Name
	service.RepoTag = storedService.RepoTag
	service.Description = storedService.Description
	service.Version = storedService.Version
	service.Cmd = storedCmd
	service.Created = storedService.Created
	service.Deleted = storedService.Deleted
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
	storedService.ConnectionID = int(service.ConnectionID)
	storedService.ImageID = string(service.ImageID)
	storedService.Name = service.Name
	storedService.RepoTag = service.RepoTag
	storedService.Description = service.Description
	storedService.Version = service.Version
	storedService.Created = service.Created
	storedService.Deleted = service.Deleted
	storedService.Size = service.Size
	return storedService
}

// AddIOPort adds input and output ports to the database
func (d *Db) AddIOPort(service Service) error {
	var err error
	for _, p := range service.Input {
		var curInputPort IOPortTable
		curInputPort.FileName = p.FileName
		curInputPort.Type = p.Type
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
		curOutputPort.FileName = p.FileName
		curOutputPort.Type = p.Type
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
func (d *Db) AddService(userID int64, service Service) error {
	// Before adding a service we need to check if the service with the same name already exists.
	// If it does, we remove it and add a new one
	var servicesFromTable []ServiceTable
	_, err := d.db.Select(&servicesFromTable,
		"SELECT * FROM services WHERE Name=? AND ConnectionID=?",
		service.Name, service.ConnectionID)
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
	if err != nil {
		return err
	}

	ownership := OwnerTable{
		UserID:     userID,
		ObjectType: "Service",
		ObjectID:   string(service.ID),
	}
	return d.db.Insert(&ownership)
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
	if err != nil {
		return err
	}

	_, err = d.db.Exec("DELETE FROM Owners WHERE ObjectType=? AND ObjectID=?",
		"Service", string(id))
	if err != nil {
		return err
	}
	return nil
}

// ListServices produces a list of all services ready to be converted into JSON
func (d *Db) ListServices() ([]Service, error) {
	var services []Service
	var servicesFromTable []ServiceTable
	_, err := d.db.Select(&servicesFromTable, "SELECT * FROM services WHERE Deleted=? ORDER BY Name", false)
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

// GetJobOwningVolume returns a service ready to be converted into JSON
func (d *Db) GetJobOwningVolume(volumeID string) (Job, error) {
	var dbjob JobTable
	err := d.db.SelectOne(&dbjob, "SELECT jobs.* FROM jobs INNER JOIN volumes ON jobs.ID = volumes.JobID WHERE volumes.ID=?", volumeID)
	if err != nil {
		return Job{}, err
	}
	return d.jobTable2Job(dbjob)
}
