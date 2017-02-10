package pier

import (
	"fmt"
	"log"
	"time"
	"github.com/pborman/uuid"
	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
)

const stagingVolumeName = "volume-stage-in"
const copyToVolumeName = "volume-to-volume"


// Pier is a master struct for gef-docker abstractions
type Pier struct {
	docker              dckr.Client
	services            *ServiceList
	jobs                *JobList
	tmpDir         	    string
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
		State:     &JobState{nil, "Created"},
	}
	p.jobs.add(job)

	go p.runJob(&job, service, inputPID)

	return job, nil
}

func (p *Pier) runJob(job *Job, service Service, inputPID string) {
	inputVolume, err := p.docker.NewVolume()
	if err != nil {
		p.jobs.setState(job.ID, JobState{def.Err(err, "Error while creating new input volume"), "Error"})
		return
	}
	log.Println("new input volume created: ", inputVolume)
	p.jobs.setInputVolume(job.ID, VolumeID(inputVolume.ID))
	{
		binds := []dckr.VolBind{
			dckr.VolBind{inputVolume.ID, "/volume", false},
		}
		exitCode, consoleOutput, err := p.docker.ExecuteImage(dckr.ImageID(stagingVolumeName), []string{inputPID}, binds, true)
		p.jobs.addTask(job.ID, "Data staging", err, exitCode, consoleOutput)

		log.Println("  staging ended: ", exitCode, ", error: ", err)
		if err != nil {
			p.jobs.setState(job.ID, JobState{def.Err(err, "Data staging failed"), "Error"})
			return
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("Data staging failed (exitCode = %v)", exitCode)
			p.jobs.setState(job.ID, JobState{nil, msg})
			return
		}
	}

	outputVolume, err := p.docker.NewVolume()
	if err != nil {
		p.jobs.setState(job.ID, JobState{def.Err(err, "Error while creating new output volume"), "Error"})
		return
	}
	log.Println("new output volume created: ", outputVolume)
	p.jobs.setOutputVolume(job.ID, VolumeID(outputVolume.ID))
	{
		binds := []dckr.VolBind{
			dckr.VolBind{inputVolume.ID, service.Input[0].Path, true},
			dckr.VolBind{outputVolume.ID, service.Output[0].Path, false},
		}
		exitCode, consoleOutput, err := p.docker.ExecuteImage(dckr.ImageID(service.imageID), nil, binds, true)
		p.jobs.addTask(job.ID, "Service execution", err, exitCode, consoleOutput)

		log.Println("  job ended: ", exitCode, ", error: ", err)
		if err != nil {
			p.jobs.setState(job.ID, JobState{def.Err(err, "Service failed"), "Error"})
			return
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("Service failed (exitCode = %v)", exitCode)
			p.jobs.setState(job.ID, JobState{nil, msg})
			return
		}
	}
	p.jobs.setState(job.ID, JobState{nil, "Ended successfully"})
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

// CopyToVolume exported
func (p *Pier) CopyToVolume(volumeID string, filePath string) (dckr.ContainerID, error){
	binds := []dckr.VolBind{
		dckr.VolBind{dckr.VolumeID(volumeID), "/root/volume", false},
	}
	//exitCode, consoleOutput, err := p.docker.ExecuteImage(dckr.ImageID(copyToVolumeName), []string{filePath, "/root"}, binds, true)
	cont, consoleOutput, err := p.docker.StartImage(dckr.ImageID(copyToVolumeName), []string{filePath, "/root"}, binds, true)

	fmt.Println(cont)
	fmt.Println(consoleOutput)

	if err != nil {
		log.Println("error: ", err)
		return
	}


	return cont

}