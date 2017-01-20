package pier

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/pborman/uuid"

	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
)

// Pier is a master struct for gef-docker abstractions
type Pier struct {
	docker   dckr.Client
	services *ServiceList
	jobs     *JobList
	tmpDir   string
}

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

func (p *Pier) BuildService(buildDir string) (Service, error) {
	image, err := p.docker.BuildImage(buildDir)
	if err != nil {
		return Service{}, def.Err(err, "docker BuildImage failed")
	}

	service := newServiceFromImage(image)
	p.services.add(service)

	return service, nil
}

func (p *Pier) ListServices() []Service {
	return p.services.list()
}

func (p *Pier) GetService(serviceID ServiceID) Service {
	return p.services.get(serviceID)
}

func (p *Pier) Run(service Service) (Job, error) {
	imageID := strings.Replace(string(service.imageID), "sha256:", "", 1)
	containerID, err := p.docker.ExecuteImage(dckr.ImageID(imageID), nil)
	if err != nil {
		return Job{}, def.Err(err, "docker ExecuteImage failed")
	}

	job := Job{
		ID:          JobID(uuid.New()),
		ServiceID:   service.ID,
		containerID: containerID,
		Status:      "Created",
		Created:     time.Now(),
	}
	p.jobs.add(job)

	return job, nil
}

func (p *Pier) ListJobs() []Job {
	jobs := p.jobs.list()
	for i := range jobs {
		p.updateMessage(&jobs[i])
	}
	return jobs
}

func (p *Pier) GetJob(jobID JobID) Job {
	job := p.jobs.get(jobID)
	p.updateMessage(&job)
	return job
}

///////////////////////////////////////////////////////////////////////////////

var statusRegExp = regexp.MustCompile("(([0-9]+)|([a-zA-Z]+)) ([a-zA-Z]+) ([a-zA-Z]+)")

func (p *Pier) updateMessage(job *Job) {
	cont, err := p.docker.InspectContainer(job.containerID)
	if err != nil {
		statusMessage := cont.State.Status
		statusMessage = strings.Replace(statusMessage, "About", "", 1)
		statusMessage = statusRegExp.ReplaceAllString(statusMessage, "")
		statusMessage = strings.Trim(statusMessage, " ")
		statusMessage = strings.Replace(statusMessage, "Exited", "Finished", 1)
		job.Status = statusMessage
	}
}
