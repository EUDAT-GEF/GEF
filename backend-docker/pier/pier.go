package pier

import (
	"fmt"
	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
	"github.com/pborman/uuid"
	"log"
	"time"
)

const stagingVolumeName = "volume-stage-in"

// Pier is a master struct for gef-docker abstractions
type Pier struct {
	docker   dckr.Client
	services *ServiceList
	jobs     *JobList
	tmpDir   string
}

// VolumeID exported
type VolumeID dckr.VolumeID

// NewPier exported
func NewPier(cfgList []def.DockerConfig, tmpDir string) (*Pier, error) {
	docker, err := dckr.NewClientFirstOf(cfgList)
	if err != nil {
		return nil, def.Err(err, "Cannot create docker client")
	}

	pier := Pier{
		docker:   docker,
		services: NewServiceList(),
		jobs:     NewJobList(),
		tmpDir:   tmpDir,
	}

	// Populate the list of services
	images, err := docker.ListImages()
	if err != nil {
		log.Println(def.Err(err, "Error while initializing services"))
	} else {
		for _, img := range images {
			pier.services.add(newServiceFromImage(img))
		}
	}

	return &pier, nil
}

// BuildService exported
func (p *Pier) BuildService(buildDir string) (Service, error) {
	image, err := p.docker.BuildImage(buildDir)
	if err != nil {
		return Service{}, def.Err(err, "docker BuildImage failed")
	}

	service := newServiceFromImage(image)
	p.services.add(service)

	return service, nil
}

// ListServices exported
func (p *Pier) ListServices() []Service {
	return p.services.list()
}

// GetService exported
func (p *Pier) GetService(serviceID ServiceID) (Service, error) {
	service, ok := p.services.get(serviceID)
	if !ok {
		return service, def.Err(nil, "not found")
	}
	return service, nil
}

// RunService exported
func (p *Pier) RunService(service Service, inputPID string) (Job, error) {
	job := Job{
		ID:        JobID(uuid.New()),
		ServiceID: service.ID,
		Created:   time.Now(),
		Input:     inputPID,
		State:     &JobState{nil, "Created", -1},
	}
	p.jobs.add(job)

	go p.runJob(&job, service, inputPID)

	return job, nil
}

func (p *Pier) runJob(job *Job, service Service, inputPID string) {
	p.jobs.setState(job.ID, JobState{nil, "Creating a new input volume", -1})
	inputVolume, err := p.docker.NewVolume()
	if err != nil {
		p.jobs.setState(job.ID, JobState{def.Err(err, "Error while creating new input volume"), "Error", 1})
		return
	}
	log.Println("new input volume created: ", inputVolume)
	p.jobs.setInputVolume(job.ID, VolumeID(inputVolume.ID))
	{
		p.jobs.setState(job.ID, JobState{nil, "Performing data staging", -1})
		binds := []dckr.VolBind{
			dckr.NewVolBind(inputVolume.ID, "/volume", false),
		}
		containerID, exitCode, consoleOutput, err := p.docker.ExecuteImage(dckr.ImageID(stagingVolumeName), []string{inputPID}, binds, true)
		p.jobs.addTask(job.ID, "Data staging", containerID, err, exitCode, consoleOutput)

		log.Println("  staging ended: ", exitCode, ", error: ", err)
		if err != nil {
			p.jobs.setState(job.ID, JobState{def.Err(err, "Data staging failed"), "Error", 1})
			return
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("Data staging failed (exitCode = %v)", exitCode)
			p.jobs.setState(job.ID, JobState{nil, msg, 1})
			return
		}
	}
	p.jobs.setState(job.ID, JobState{nil, "Creating a new output volume", -1})
	outputVolume, err := p.docker.NewVolume()
	if err != nil {
		p.jobs.setState(job.ID, JobState{def.Err(err, "Error while creating new output volume"), "Error", 1})
		return
	}
	log.Println("new output volume created: ", outputVolume)
	p.jobs.setOutputVolume(job.ID, VolumeID(outputVolume.ID))
	{
		p.jobs.setState(job.ID, JobState{nil, "Executing the service", -1})
		binds := []dckr.VolBind{
			dckr.NewVolBind(inputVolume.ID, service.Input[0].Path, true),
			dckr.NewVolBind(outputVolume.ID, service.Output[0].Path, false),
		}
		containerID, exitCode, consoleOutput, err := p.docker.ExecuteImage(dckr.ImageID(service.imageID), nil, binds, true)
		p.jobs.addTask(job.ID, "Service execution", containerID, err, exitCode, consoleOutput)

		log.Println("  job ended: ", exitCode, ", error: ", err)
		if err != nil {
			p.jobs.setState(job.ID, JobState{def.Err(err, "Service failed"), "Error", 1})
			return
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("Service failed (exitCode = %v)", exitCode)
			p.jobs.setState(job.ID, JobState{nil, msg, 1})
			return
		}
	}
	p.jobs.setState(job.ID, JobState{nil, "Ended successfully", 0})
}

// ListJobs exported
func (p *Pier) ListJobs() []Job {
	return p.jobs.list()
}

// GetJob exported
func (p *Pier) GetJob(jobID JobID) (Job, error) {
	job, ok := p.jobs.get(jobID)

	if !ok {
		return job, def.Err(nil, "not found")
	}
	return job, nil
}

// RemoveJob exported
func (p *Pier) RemoveJob(jobID JobID) (JobID, error) {
	job, ok := p.jobs.get(jobID)
	if !ok {
		return jobID, def.Err(nil, "not found")
	}

	// Removing volumes
	err := p.docker.RemoveVolume(dckr.VolumeID(job.InputVolume))
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
	p.jobs.remove(jobID)
	return jobID, nil
}
