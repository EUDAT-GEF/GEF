package pier

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier/db"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
	"github.com/pborman/uuid"
)

const stagingVolumeName = "volume-stage-in"
const servicesFolder = "../services/"

// Pier is a master struct for gef-docker abstractions
type Pier struct {
	docker   dckr.Client
	db       *db.Db

	tmpDir   string
}

// NewPier exported
func NewPier(cfgList []def.DockerConfig, tmpDir string, dataBase *db.Db) (*Pier, error) {
	docker, err := dckr.NewClientFirstOf(cfgList)

	if err != nil {
		return nil, def.Err(err, "Cannot create docker client")
	}

	pier := Pier{
		docker:   docker,
		db:       dataBase,
		tmpDir:   tmpDir,
	}
	return &pier, nil
}

// BuildService builds a services based on the content of the provided folder
func (p *Pier) BuildService(buildDir string) (db.Service, error) {
	image, err := p.docker.BuildImage(buildDir)
	if err != nil {
		return db.Service{}, def.Err(err, "docker BuildImage failed")
	}

	service := p.db.NewServiceFromImage(image)
	err = p.db.AddService(service)
	if err != nil {
		return db.Service{}, def.Err(err, "could not add a new service to the database")
	}

	return service, nil
}

// ListServices lists all existing services
func (p *Pier) ListServices() ([]db.Service, error) {
	return p.db.ListServices()
}

// GetService returns a service by ID
func (p *Pier) GetService(serviceID db.ServiceID) (db.Service, error) {
	service, err := p.db.GetService(serviceID)
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
		State:     &db.JobState{"", "Created", -1},
	}

	err := p.db.AddJob(job)

	go p.runJob(&job, service, inputPID)

	return job, err
}

// runJob runs a job
func (p *Pier) runJob(job *db.Job, service db.Service, inputPID string) {
	p.db.SetJobState(job.ID, db.JobState{"", "Creating a new input volume", -1})
	inputVolume, err := p.docker.NewVolume()
	if err != nil {
		p.db.SetJobState(job.ID, db.JobState{"Error while creating new input volume", "Error", 1})
		return
	}
	log.Println("new input volume created: ", inputVolume)
	p.db.SetJobInputVolume(job.ID, db.VolumeID(inputVolume.ID))
	{
		p.db.SetJobState(job.ID, db.JobState{"", "Performing data staging", -1})
		binds := []dckr.VolBind{
			dckr.NewVolBind(inputVolume.ID, "/volume", false),
		}
		containerID, exitCode, consoleOutput, err := p.docker.ExecuteImage(dckr.ImageID(stagingVolumeName), []string{inputPID}, binds, true)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		p.db.AddJobTask(job.ID, "Data staging", containerID, errMsg, exitCode, consoleOutput)
		//fmt.Println(containerID, consoleOutput)

		log.Println("  staging ended: ", exitCode, ", error: ", err)
		if err != nil {
			p.db.SetJobState(job.ID, db.JobState{"Data staging failed", "Error", 1})
			return
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("Data staging failed (exitCode = %v)", exitCode)
			p.db.SetJobState(job.ID, db.JobState{"", msg, 1})
			return
		}
	}
	p.db.SetJobState(job.ID, db.JobState{"", "Creating a new output volume", -1})
	outputVolume, err := p.docker.NewVolume()
	if err != nil {
		p.db.SetJobState(job.ID, db.JobState{"Error while creating new output volume", "Error", 1})
		return
	}
	log.Println("new output volume created: ", outputVolume)
	p.db.SetJobOutputVolume(job.ID, db.VolumeID(outputVolume.ID))
	{
		p.db.SetJobState(job.ID, db.JobState{"", "Executing the service", -1})
		binds := []dckr.VolBind{
			dckr.NewVolBind(inputVolume.ID, service.Input[0].Path, true),
			dckr.NewVolBind(outputVolume.ID, service.Output[0].Path, false),
		}
		containerID, exitCode, consoleOutput, err := p.docker.ExecuteImage(dckr.ImageID(service.ImageID), nil, binds, true)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		p.db.AddJobTask(job.ID, "Service execution", containerID, errMsg, exitCode, consoleOutput)
		//fmt.Println(containerID, consoleOutput)

		log.Println("  job ended: ", exitCode, ", error: ", err)
		if err != nil {
			p.db.SetJobState(job.ID, db.JobState{"Service failed", "Error", 1})
			return
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("Service failed (exitCode = %v)", exitCode)
			p.db.SetJobState(job.ID, db.JobState{"", msg, 1})
			return
		}
	}
	p.db.SetJobState(job.ID, db.JobState{"", "Ended successfully", 0})
}

// RemoveJob removes a job by ID
func (p *Pier) RemoveJob(jobID db.JobID) (db.JobID, error) {
	job, err := p.db.GetJob(jobID)
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
	p.db.RemoveJob(jobID)
	return jobID, nil
}

// ListJobs lists all existing jobs
func (p *Pier) ListJobs() ([]db.Job, error) {
	return p.db.ListJobs()
}

// GetJob returns a job by ID
func (p *Pier) GetJob(jobID db.JobID) (db.Job, error) {
	job, err := p.db.GetJob(jobID)
	if err != nil {
		return job, def.Err(nil, "not found")
	}
	return job, nil
}

// PopulateServiceTable reads the "services" folder, builds images, and adds all the necessary information
// to the database
func (p *Pier) PopulateServiceTable() error {
	log.Println("Reading folder with Dockerfiles for serices: " + servicesFolder)
	doesExist := true
	_, err := os.Stat(servicesFolder)
	if os.IsNotExist(err) {
		doesExist = false
	}
	if doesExist {
		files, _ := ioutil.ReadDir(servicesFolder)
		for _, f := range files {
			if f.IsDir() {
				log.Print("Opening folder: " + f.Name())
				img, err := p.docker.BuildImage(filepath.Join(servicesFolder, f.Name()))

				if err != nil {
					log.Print("failed to create a service")
				} else {
					log.Print("service has been created")
					error := p.db.AddService(p.db.NewServiceFromImage(img))
					if error != nil {
						log.Print(error)
					}
				}
			}
		}
	}

	return nil
}
