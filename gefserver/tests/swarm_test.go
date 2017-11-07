package tests

import (
	"log"
	"os"
	"testing"

	"github.com/EUDAT-GEF/GEF/gefserver/db"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"github.com/EUDAT-GEF/GEF/gefserver/pier"
)

func TestMain(m *testing.M) {
	setSwarmMode(false)
	code := m.Run()
	os.Exit(code)
}

func TestAgainInSwarm(t *testing.T) {
	log.Println("************************************")
	log.Println("* Running tests for the swarm mode *")
	log.Println("************************************")

	setSwarmMode(true) // switching to a swarm
	TestClient(t)
	TestExecution(t)
	TestJobTimeOut(t)
	TestMultipleInputsAndOutputs(t)
	TestServer(t)
	setSwarmMode(false) // leaving a swarm
}

func setSwarmMode(activate bool) {
	config, err := def.ReadConfigFile(configFilePath)
	if err != nil {
		log.Fatal(def.Err(err, "reading config files failed"))
		os.Exit(1)
	}

	// overwrite this because when testing we're in a different working directory
	config.Pier.InternalServicesFolder = internalServicesFolder

	db, file, err := db.InitDbForTesting()
	if err != nil {
		log.Fatal(def.Err(err, "creating test db failed"))
		os.Exit(1)
	}
	defer db.Close()
	defer os.Remove(file)

	pier, err := pier.NewPier(&db, config.Pier, config.TmpDir, config.Timeouts)
	if err != nil {
		log.Fatal(def.Err(err, "creating new pier failed"))
		os.Exit(1)
	}

	connID, err := pier.AddDockerConnection(0, config.Docker)
	if err != nil {
		log.Fatal(def.Err(err, "setting docker connection failed"))
		os.Exit(1)
	}

	if activate {
		_, err = pier.InitiateSwarmMode(connID, "127.0.0.1", "127.0.0.1")
		if err != nil {
			log.Fatal(def.Err(err, "switching to the swarm mode or leaving swarm failed"))
			os.Exit(1)
		}
	} else {
		err = pier.LeaveIfInSwarmMode(connID)
		if err != nil {
			log.Fatal(def.Err(err, "leaving swarm failed"))
			os.Exit(1)
		}
	}
}
