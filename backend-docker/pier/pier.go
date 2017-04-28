package pier

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"strconv"
	"strings"

	"github.com/EUDAT-GEF/GEF/backend-docker/db"
	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
	"github.com/pborman/uuid"
)

const GefSrvLabelPrefix = "eudat.gef.service." // GefSrvLabelPrefix is the prefix identifying GEF related labels
const stagingVolumeName = "volume-stage-in"
const servicesFolder = "../services/"
const internalServicesFolder = "_internal"

// Pier is a master struct for gef-docker abstractions
type Pier struct {
	docker dckr.Client
	db     *db.Db
	tmpDir string
	limits def.LimitConfig
}

// NewPier exported
func NewPier(cfgList []def.DockerConfig, tmpDir string, cntrLimits def.LimitConfig, dataBase *db.Db) (*Pier, error) {
	docker, err := dckr.NewClientFirstOf(cfgList)

	if err != nil {
		return nil, def.Err(err, "Cannot create docker client")
	}

	pier := Pier{
		docker: docker,
		db:     dataBase,
		tmpDir: tmpDir,
		limits: cntrLimits,
	}

	return &pier, nil
}

// BuildService builds a services based on the content of the provided folder
func (p *Pier) BuildService(buildDir string) (db.Service, error) {
	image, err := p.docker.BuildImage(buildDir)
	if err != nil {
		return db.Service{}, def.Err(err, "docker BuildImage failed")
	}

	service := NewServiceFromImage(image)
	err = p.db.AddService(service)
	if err != nil {
		return db.Service{}, def.Err(err, "could not add a new service to the database")
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
		containerID, exitCode, consoleOutput, err := p.docker.ExecuteImage(dckr.ImageID(stagingVolumeName), []string{inputPID}, binds, p.limits, true)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		p.db.AddJobTask(job.ID, "Data staging", string(containerID), errMsg, exitCode, consoleOutput)

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
		containerID, exitCode, consoleOutput, err := p.docker.ExecuteImage(dckr.ImageID(service.ImageID), nil, binds, p.limits, true)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		p.db.AddJobTask(job.ID, "Service execution", string(containerID), errMsg, exitCode, consoleOutput)

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
func (p *Pier) RemoveJob(jobID db.JobID) (db.Job, error) {
	job, err := p.db.GetJob(jobID)
	if err != nil {
		return job, def.Err(nil, "not found")
	}

	// Removing volumes
	err = p.docker.RemoveVolume(dckr.VolumeID(job.InputVolume))
	if err != nil {
		return job, def.Err(err, "Input volume is not set")
	}
	err = p.docker.RemoveVolume(dckr.VolumeID(job.OutputVolume))
	if err != nil {
		return job, def.Err(err, "Output volume is not set")
	}

	// Stopping the latest or the current task (if it is running)
	if len(job.Tasks) > 0 {
		p.docker.RemoveContainer(string(job.Tasks[len(job.Tasks)-1].ContainerID))
	}

	// Removing the job from the list
	p.db.RemoveJob(jobID)
	return job, nil
}

// buildServicesFromFolder builds an image from the specified folder and assigns a tag to it based on the corresponding folder name
func (p *Pier) buildServicesFromFolder(inputFolder string) error {
	_, err := os.Stat(inputFolder)
	if os.IsNotExist(err) {
		return nil
	}

	files, _ := ioutil.ReadDir(inputFolder)
	for _, f := range files {
		if f.IsDir() && f.Name() != internalServicesFolder {
			log.Print("Opening folder: " + f.Name())
			img, err := p.docker.BuildImage(filepath.Join(inputFolder, f.Name()))

			if err != nil {
				log.Print("failed to create a service: ", err)
			} else {
				log.Print("service has been created")

				err = p.docker.TagImage(string(img.ID), f.Name(), "latest")
				if err != nil {
					log.Print("could not tag the service")
				}

				img, err = p.docker.InspectImage(img.ID)
				if err != nil {
					log.Print("failed to inspect the image: ", err)
				}

				err = p.db.AddService(NewServiceFromImage(img))
				if err != nil {
					log.Print("failed to add the service to the database: ", err)
				}
			}
		}
	}
	return nil
}

// PopulateServiceTable reads the "services" folder, builds images, and adds all the necessary information
// to the database
func (p *Pier) PopulateServiceTable() error {
	log.Println("Reading the folder with Dockerfiles for internal services: " + filepath.Join(servicesFolder, internalServicesFolder))
	err := p.buildServicesFromFolder(filepath.Join(servicesFolder, internalServicesFolder))
	if err != nil {
		return err
	}
	log.Println("Reading the folder with Dockerfiles for demo services: " + servicesFolder)
	err = p.buildServicesFromFolder(servicesFolder)
	return err
}

func (p *Pier) ImportImage(imageFilePath string) (db.Service, error) {
	imageID, err := p.docker.ImportImageFromTar(imageFilePath)
	if err != nil {
		return db.Service{}, def.Err(err, "docker ImportImage failed")
	}

	image, err := p.docker.InspectImage(imageID)

	if err != nil {
		return db.Service{}, err
	}

	service := NewServiceFromImage(image)
	err = p.db.AddService(service)
	if err != nil {
		return db.Service{}, def.Err(err, "could not add a new service to the database")
	}

	return service, nil
}

// NewServiceFromImage extracts metadata and creates a valid GEF service
func NewServiceFromImage(image dckr.Image) db.Service {
	srv := db.Service{
		ID:      db.ServiceID(uuid.New()),
		ImageID: db.ImageID(image.ID),
		RepoTag: image.RepoTag,
		Created: image.Created,
		Size:    image.Size,
	}

	for k, v := range image.Labels {
		if !strings.HasPrefix(k, GefSrvLabelPrefix) {
			continue
		}
		k = k[len(GefSrvLabelPrefix):]
		ks := strings.Split(k, ".")
		if len(ks) == 0 {
			continue
		}
		switch ks[0] {
		case "name":
			srv.Name = v
		case "description":
			srv.Description = v
		case "version":
			srv.Version = v
		case "input":
			addVecValue(&srv.Input, ks[1:], v)
		case "output":
			addVecValue(&srv.Output, ks[1:], v)
		default:
			log.Println("Unknown GEF service label: ", k, "=", v)
		}
	}

	{
		in := make([]db.IOPort, 0, len(srv.Input))
		for _, p := range srv.Input {
			if p.Path != "" {
				p.ID = fmt.Sprintf("input%d", len(in))
				in = append(in, p)
			}
		}
		srv.Input = in
	}
	{
		out := make([]db.IOPort, 0, len(srv.Output))
		for _, p := range srv.Output {
			if p.Path != "" {
				p.ID = fmt.Sprintf("output%d", len(out))
				out = append(out, p)
			}
		}
		srv.Output = out
	}

	return srv
}

// addVecValue is used by the NewServiceFromImage
func addVecValue(vec *[]db.IOPort, ks []string, value string) {
	if len(ks) < 2 {
		log.Println("ERROR: GEF service label I/O key error (need 'port number . key name')", ks)
		return
	}
	id, err := strconv.ParseUint(ks[0], 10, 8)
	if err != nil {
		log.Println("ERROR: GEF service label: expecting integer argument for IOPort, instead got: ", ks)
	}
	for len(*vec) < int(id)+1 {
		*vec = append(*vec, db.IOPort{})
	}
	switch ks[1] {
	case "name":
		(*vec)[id].Name = value
	case "path":
		(*vec)[id].Path = value
	}
}
