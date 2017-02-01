package tests

import (
	"log"
	"testing"

	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier"
)

const testPID = "11304/a3d012ca-4e23-425e-9e2a-1e6a195b966f"

var configFilePath = "../config.json"
var config []def.DockerConfig

func TestClient(t *testing.T) {
	config, err := def.ReadConfigFile(configFilePath)
	checkMsg(t, err, "reading config files")

	pier, err := pier.NewPier(config.Docker, config.TmpDir)
	checkMsg(t, err, "creating new pier")

	before := pier.ListServices()
	checkMsg(t, err, "listing services failed")

	service, err := pier.BuildService("./docker_test")
	checkMsg(t, err, "build service failed")
	log.Println("built service:", service)

	after := pier.ListServices()
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
	log.Println("job: ", job)

	jobList := pier.ListJobs()
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

	j, err := pier.GetJob(job.ID)
	checkMsg(t, err, "getting job failed")
	if j.ID != job.ID {
		t.Error("job retrieval based on id failure")
		t.FailNow()
	}
}

func TestExecution(t *testing.T) {
	config, err := def.ReadConfigFile(configFilePath)
	checkMsg(t, err, "reading config files")

	pier, err := pier.NewPier(config.Docker, config.TmpDir)
	checkMsg(t, err, "creating new pier")

	service, err := pier.BuildService("./docker_clone")
	checkMsg(t, err, "build service failed")
	log.Println("built service:", service)

	job, err := pier.RunService(service, testPID)
	checkMsg(t, err, "running service failed")

	log.Println("job: ", job)
	jobid := job.ID

	for job.State.Status == "Created" {
		job, err = pier.GetJob(jobid)
		checkMsg(t, err, "getting job failed")
	}

	expect(t, job.State.Error != nil, "job error")

	files, err := pier.ListFiles(job.OutputVolume)
	checkMsg(t, err, "getting volume failed")

	expect(t, len(files) == 1, "bad returned files")
}