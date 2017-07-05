package tests

import (
	"log"
	"testing"

	"github.com/EUDAT-GEF/GEF/gefserver/db"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"github.com/EUDAT-GEF/GEF/gefserver/pier"
)

const testPID = "11304/a3d012ca-4e23-425e-9e2a-1e6a195b966f"

var configFilePath = "../config.json"
var internalServicesFolder = "../../services/_internal"

func TestClient(t *testing.T) {
	config, err := def.ReadConfigFile(configFilePath)
	checkMsg(t, err, "reading config files")

	db, err := db.InitDb()
	checkMsg(t, err, "creating db")
	defer db.Db.Close()

	pier, err := pier.NewPier(&db, config.TmpDir)
	checkMsg(t, err, "creating new pier")

	err = pier.SetDockerConnection(config.Docker, config.Limits, config.Timeouts, internalServicesFolder)
	checkMsg(t, err, "setting docker connection")

	before, err := db.ListServices()
	checkMsg(t, err, "listing services failed")

	service, err := pier.BuildService("./docker_test")
	checkMsg(t, err, "build service failed")
	log.Println("test service built:", service)

	after, err := db.ListServices()
	checkMsg(t, err, "listing services failed")

	errstr := "Cannot find new service in list"
	for _, x := range after {
		if x.ID == service.ID {
			errstr = ""
			break
		}
	}

	if errstr != "" {
		t.Error("before is: ", len(before), before)
		t.Error("service is: ", service)
		t.Error("after is: ", len(after), after)
		t.Error("")
		t.Error(errstr)
		t.Fail()
		return
	}

	job, err := pier.RunService(service, "")
	checkMsg(t, err, "running service failed")
	log.Println("test job: ", job)

	jobList, err := db.ListJobs()
	if len(jobList) == 0 {
		t.Error("cannot find any job")
		t.FailNow()
	}

	found := false
	for _, j := range jobList {
		if j.ID == job.ID {
			found = true
		}
	}
	if !found {
		t.Error("cannot find the new job in the list of all jobs")
		t.FailNow()
	}

	j, err := db.GetJob(job.ID)
	checkMsg(t, err, "getting job failed")
	if j.ID != job.ID {
		t.Error("job retrieval based on id failure")
		t.FailNow()
	}
}

func TestExecution(t *testing.T) {
	config, err := def.ReadConfigFile(configFilePath)
	checkMsg(t, err, "reading config files")

	db, err := db.InitDb()
	checkMsg(t, err, "creating db")
	defer db.Db.Close()

	pier, err := pier.NewPier(&db, config.TmpDir)
	checkMsg(t, err, "creating new pier")

	err = pier.SetDockerConnection(config.Docker, config.Limits, config.Timeouts, internalServicesFolder)
	checkMsg(t, err, "setting docker connection")

	service, err := pier.BuildService("./clone_test")
	checkMsg(t, err, "build service failed")
	log.Println("test service built:", service)

	job, err := pier.RunService(service, testPID)
	checkMsg(t, err, "running service failed")

	log.Println("test job: ", job)
	jobid := job.ID

	for job.State.Code == -1 {
		job, err = db.GetJob(jobid)
		checkMsg(t, err, "getting job failed")
	}

	if job.State.Error != "" {
		log.Println("test job error:")
		log.Println("state: ", job.State)
		for i, t := range job.Tasks {
			log.Println("task ", i, ":", t)
		}
	}
	expect(t, job.State.Error == "", "job error")

	files, err := pier.ListFiles(job.OutputVolume, "")
	checkMsg(t, err, "getting volume failed")

	expect(t, len(files) == 1, "bad returned files")

	_, err = pier.RemoveJob(jobid)
	checkMsg(t, err, "removing job failed")
}
