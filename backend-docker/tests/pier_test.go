package tests

import (
	"log"
	"testing"

	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier"
)

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

	job, err := pier.Run(service)
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

	j := pier.GetJob(job.ID)
	if j.ID != job.ID {
		t.Error("job retrieval based on id failure")
		t.FailNow()
	}
}
