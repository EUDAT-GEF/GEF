package pier

import (
	"fmt"
	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
	"github.com/pborman/uuid"
	"log"
	"time"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/db"
)

const stagingVolumeName = "volume-stage-in"

// Pier is a master struct for gef-docker abstractions
type Pier struct {
	docker   dckr.Client
	dataBase *db.Db
	services []db.Service
	jobs     []db.Job
	tmpDir   string
}

// VolumeID exported
type VolumeID dckr.VolumeID

// NewPier exported
func NewPier(cfgList []def.DockerConfig, tmpDir string, dataBase *db.Db) (*Pier, error) {
	docker, err := dckr.NewClientFirstOf(cfgList)

	var allServices []db.Service
	var allJobs []db.Job
	allServices, err = dataBase.ListServices()
	if err != nil {
		return nil, def.Err(err, "Cannot retrieve a list of services")
	}
	allJobs, err = dataBase.ListJobs()
	if err != nil {
		return nil, def.Err(err, "Cannot retrieve a list of jobs")
	}
	if err != nil {
		return nil, def.Err(err, "Cannot create docker client")
	}

	pier := Pier{
		docker:   docker,
		dataBase: dataBase,
		services: allServices,
		jobs:     allJobs,
		tmpDir:   tmpDir,
	}

	// Populate the list of services
	/*images, err := docker.ListImages()
	if err != nil {
		log.Println(def.Err(err, "Error while initializing services"))
	} else {
		for _, img := range images {
			err = dataBase.AddService(dataBase.NewServiceFromImage(img))
			if err != nil {
				return nil, def.Err(err, "Cannot retrieve a list of services")
			}
		}
	}*/

	return &pier, nil
}

// BuildService exported
func (p *Pier) BuildService(buildDir string) (db.Service, error) {
	image, err := p.docker.BuildImage(buildDir)
	if err != nil {
		return db.Service{}, def.Err(err, "docker BuildImage failed")
	}

	service := p.dataBase.NewServiceFromImage(image)
	err = p.dataBase.AddService(service)
	if err != nil {
		return db.Service{}, def.Err(err, "could not add a new service to the database")
	}

	return service, nil
}

// ListServices exported
func (p *Pier) ListServices() ([]db.Service, error) {
	return p.dataBase.ListServices()
}

// GetService exported
func (p *Pier) GetService(serviceID db.ServiceID) (db.Service, error) {
	service, err := p.dataBase.GetService(serviceID)
	if err != nil {
		return service, def.Err(nil, "not found")
	}
	return service, nil
}

// RunService exported
func (p *Pier) RunService(service db.Service, inputPID string) (db.Job, error) {
	job := db.Job{
		ID:        db.JobID(uuid.New()),
		ServiceID: service.ID,
		Created:   time.Now(),
		Input:     inputPID,
		State:     &db.JobState{nil, "Created", -1},
	}

	err := p.dataBase.AddJob(job)

	go p.runJob(&job, service, inputPID)

	return job, err
}

func (p *Pier) runJob(job *db.Job, service db.Service, inputPID string) {
	p.dataBase.SetJobState(job.ID, db.JobState{nil, "Creating a new input volume", -1})
	inputVolume, err := p.docker.NewVolume()
	if err != nil {
		p.dataBase.SetJobState(job.ID, db.JobState{def.Err(err, "Error while creating new input volume"), "Error", 1})
		return
	}
	log.Println("new input volume created: ", inputVolume)
	p.dataBase.SetJobInputVolume(job.ID, db.VolumeID(inputVolume.ID))
	{
		p.dataBase.SetJobState(job.ID, db.JobState{nil, "Performing data staging", -1})
		binds := []dckr.VolBind{
			dckr.NewVolBind(inputVolume.ID, "/volume", false),
		}
		containerID, exitCode, consoleOutput, err := p.docker.ExecuteImage(dckr.ImageID(stagingVolumeName), []string{inputPID}, binds, true)
		p.dataBase.AddJobTask(job.ID, "Data staging", containerID, err, exitCode, consoleOutput)

		log.Println("  staging ended: ", exitCode, ", error: ", err)
		if err != nil {
			p.dataBase.SetJobState(job.ID, db.JobState{def.Err(err, "Data staging failed"), "Error", 1})
			return
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("Data staging failed (exitCode = %v)", exitCode)
			p.dataBase.SetJobState(job.ID, db.JobState{nil, msg, 1})
			return
		}
	}
	p.dataBase.SetJobState(job.ID, db.JobState{nil, "Creating a new output volume", -1})
	outputVolume, err := p.docker.NewVolume()
	if err != nil {
		p.dataBase.SetJobState(job.ID, db.JobState{def.Err(err, "Error while creating new output volume"), "Error", 1})
		return
	}
	log.Println("new output volume created: ", outputVolume)
	p.dataBase.SetJobOutputVolume(job.ID, db.VolumeID(outputVolume.ID))
	{
		p.dataBase.SetJobState(job.ID, db.JobState{nil, "Executing the service", -1})
		binds := []dckr.VolBind{
			dckr.NewVolBind(inputVolume.ID, service.Input[0].Path, true),
			dckr.NewVolBind(outputVolume.ID, service.Output[0].Path, false),
		}
		containerID, exitCode, consoleOutput, err := p.docker.ExecuteImage(dckr.ImageID(service.ImageID), nil, binds, true)
		p.dataBase.AddJobTask(job.ID, "Service execution", containerID, err, exitCode, consoleOutput)

		log.Println("  job ended: ", exitCode, ", error: ", err)
		if err != nil {
			p.dataBase.SetJobState(job.ID, db.JobState{def.Err(err, "Service failed"), "Error", 1})
			return
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("Service failed (exitCode = %v)", exitCode)
			p.dataBase.SetJobState(job.ID, db.JobState{nil, msg, 1})
			return
		}
	}
	p.dataBase.SetJobState(job.ID, db.JobState{nil, "Ended successfully", 0})
}

// ListJobs exported
func (p *Pier) ListJobs() ([]db.Job, error) {
	return p.dataBase.ListJobs()
}

// GetJob exported
func (p *Pier) GetJob(jobID db.JobID) (db.Job, error) {
	job, err := p.dataBase.GetJob(jobID)

	if err != nil {
		return job, def.Err(nil, "not found")
	}
	return job, nil
}

// RemoveJob exported
func (p *Pier) RemoveJob(jobID db.JobID) (db.JobID, error) {
	job, err := p.dataBase.GetJob(jobID)
	if err != nil {
		return jobID, def.Err(nil, "not found")
	}

	// Removing volumes
	err = p.docker.RemoveVolume(dckr.VolumeID(job.InputVolume))
	if err != nil {
		return jobID, def.Err(err, "Input volume is not set")
	}
	err = p.docker.RemoveVolume(dckr.VolumeID(job.OutputVolume))
	if err != nil {
		return jobID, def.Err(err, "Output volume is not set")
	}

	// Stopping the latest or the current task (if it is running)
	if len(job.Tasks) > 0 {
		p.docker.RemoveContainer(string(job.Tasks[len(job.Tasks)-1].ContainerID))
	}

	// Removing the job from the list
	p.dataBase.RemoveJob(jobID)
	return jobID, nil
}
