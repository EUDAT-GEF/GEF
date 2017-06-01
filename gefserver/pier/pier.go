package pier

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"strconv"
	"strings"

	"github.com/EUDAT-GEF/GEF/gefserver/db"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"github.com/EUDAT-GEF/GEF/gefserver/pier/internal/dckr"
	"github.com/pborman/uuid"
)

// GefSrvLabelPrefix is the prefix identifying GEF related labels
const GefSrvLabelPrefix = "eudat.gef.service."
const InternalImagePrefix = "internal_"
const GefImageTag = "gef"
const ServiceImagePrefix = "service_"

// Pier is a master struct for gef-docker abstractions
type Pier struct {
	docker *dockerConnection
	db     *db.Db
	tmpDir string
}

type dockerConnection struct {
	client         dckr.Client
	limits         def.LimitConfig
	stageIn        internalImage
	fileList       internalImage
	copyFromVolume internalImage
}

type internalImage struct {
	id      dckr.ImageID
	repoTag string
	cmd     []string
}

// NewPier creates a new pier with all the needed setup
func NewPier(dataBase *db.Db, tmpDir string) (*Pier, error) {
	pier := Pier{
		db:     dataBase,
		docker: nil,
		tmpDir: tmpDir,
	}
	log.Println("Pier created")
	return &pier, nil
}

// SetDockerConnection instantiates the docker client and sets the pier's docker connection
func (p *Pier) SetDockerConnection(config def.DockerConfig, limits def.LimitConfig, internalServicesFolder string) error {
	client, err := dckr.NewClient(config)
	if err != nil {
		return def.Err(err, "Cannot create docker client for config:", config)
	}

	buildInternalImage := func(docker dckr.Client, name string) (internalImage, error) {
		log.Print("building internal service: " + name)
		path := filepath.Join(internalServicesFolder, name)
		abspath, err := filepath.Abs(path)
		var newImage internalImage
		if err != nil {
			return newImage, def.Err(err, "absolute filepath failed: %s", path)
		}
		img, err := docker.BuildImage(abspath)
		if err != nil {
			return newImage, def.Err(err, "internal image build failed: %s", abspath)
		}
		err = docker.TagImage(string(img.ID), InternalImagePrefix+string(img.ID), GefImageTag)
		if err != nil {
			return newImage, def.Err(err, "could not tag an internal service: %s", string(img.ID))
		}
		newImage.id = img.ID
		newImage.cmd = img.Cmd
		newImage.repoTag = img.RepoTag
		return newImage, nil
	}

	stageInImage, err := buildInternalImage(client, "volume-stage-in")
	if err != nil {
		return err
	}
	fileListImage, err := buildInternalImage(client, "volume-filelist")
	if err != nil {
		return err
	}
	copyFromVolumeImage, err := buildInternalImage(client, "copy-from-volume")
	if err != nil {
		return err
	}

	p.docker = &dockerConnection{
		client,
		limits,
		stageInImage,
		fileListImage,
		copyFromVolumeImage,
	}
	return nil
}

// BuildService builds a services based on the content of the provided folder
func (p *Pier) BuildService(buildDir string) (db.Service, error) {
	image, err := p.docker.client.BuildImage(buildDir)
	if err != nil {
		return db.Service{}, def.Err(err, "docker BuildImage failed")
	}
	log.Println("Tagging the image")
	err = p.docker.client.TagImage(string(image.ID), ServiceImagePrefix+string(image.ID), GefImageTag)
	if err != nil {
		return db.Service{}, def.Err(err, "could not tag a service image: %s", string(image.ID))
	}

	service := NewServiceFromImage(image)
	service.RepoTag = ServiceImagePrefix+string(image.ID) + ":" + GefImageTag
	err = p.db.AddService(service)
	if err != nil {
		return db.Service{}, def.Err(err, "could not add a new service to the database")
	}

	return service, nil
}

// RunService exported
func (p *Pier) RunService(service db.Service, inputPID string) (db.Job, error) {
	jobState := db.NewJobStateOk("Created", -1)
	job := db.Job{
		ID:        db.JobID(uuid.New()),
		ServiceID: service.ID,
		Created:   time.Now(),
		Input:     inputPID,
		State:     &jobState,
	}

	err := p.db.AddJob(job)

	go p.runJob(&job, service, inputPID)

	return job, err
}

func (p *Pier) runJob(job *db.Job, service db.Service, inputPID string) {
	err2str := func(err error) string {
		if err == nil {
			return ""
		}
		return err.Error()
	}

	var err error
	var inputVolume dckr.Volume
	{
		p.db.SetJobState(job.ID, db.NewJobStateOk("Creating a new input volume", -1))
		inputVolume, err = p.docker.client.NewVolume()
		if err != nil {
			p.db.SetJobState(job.ID, db.NewJobStateError("Error while creating new input volume", 1))
			return
		}
		p.db.SetJobInputVolume(job.ID, db.VolumeID(inputVolume.ID))
	}

	{
		p.db.SetJobState(job.ID, db.NewJobStateOk("Performing data staging", -1))
		binds := []dckr.VolBind{
			dckr.NewVolBind(inputVolume.ID, "/volume", false),
		}

		containerID, exitCode, output, err := p.docker.client.ExecuteImage(
			string(p.docker.stageIn.id),
			p.docker.stageIn.repoTag,
			append(p.docker.stageIn.cmd, inputPID),
			binds,
			p.docker.limits,
			true)

		p.db.AddJobTask(job.ID, "Data staging", string(containerID), err2str(err), exitCode, output)
		if err != nil {
			p.db.SetJobState(job.ID, db.NewJobStateError("Data staging failed", 1))
			return
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("Data staging failed (exitCode = %v)", exitCode)
			p.db.SetJobState(job.ID, db.NewJobStateOk(msg, 1))
			return
		}
	}

	var outputVolume dckr.Volume
	{
		p.db.SetJobState(job.ID, db.NewJobStateOk("Creating a new output volume", -1))
		outputVolume, err = p.docker.client.NewVolume()
		if err != nil {
			p.db.SetJobState(job.ID, db.NewJobStateError("Error while creating new output volume", 1))
			return
		}
		p.db.SetJobOutputVolume(job.ID, db.VolumeID(outputVolume.ID))
	}

	{
		p.db.SetJobState(job.ID, db.NewJobStateOk("Executing the service", -1))
		binds := []dckr.VolBind{
			dckr.NewVolBind(inputVolume.ID, service.Input[0].Path, true),
			dckr.NewVolBind(outputVolume.ID, service.Output[0].Path, false),
		}
		containerID, exitCode, output, err := p.docker.client.ExecuteImage(
			string(service.ImageID),
			service.RepoTag,
			service.Cmd,
			binds,
			p.docker.limits,
			true)
		p.db.AddJobTask(job.ID, "Service execution", string(containerID), err2str(err), exitCode, output)

		//log.Println("  job ended: ", exitCode, ", error: ", err)
		if err != nil {
			p.db.SetJobState(job.ID, db.NewJobStateError("Service failed", 1))
			return
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("Service failed (exitCode = %v)", exitCode)
			p.db.SetJobState(job.ID, db.NewJobStateOk(msg, 1))
			return
		}
	}

	p.db.SetJobState(job.ID, db.NewJobStateOk("Ended successfully", 0))
}

// RemoveJob removes a job by ID
func (p *Pier) RemoveJob(jobID db.JobID) (db.Job, error) {
	job, err := p.db.GetJob(jobID)
	if err != nil {
		return job, def.Err(nil, "not found")
	}

	// Removing volumes
	err = p.docker.client.RemoveVolume(dckr.VolumeID(job.InputVolume))
	if err != nil {
		return job, def.Err(err, "Input volume is not set")
	}
	err = p.docker.client.RemoveVolume(dckr.VolumeID(job.OutputVolume))
	if err != nil {
		return job, def.Err(err, "Output volume is not set")
	}

	// Stopping the latest or the current task (if it is running)
	if len(job.Tasks) > 0 {
		p.docker.client.RemoveContainer(string(job.Tasks[len(job.Tasks)-1].ContainerID))
	}

	// Removing the job from the list
	p.db.RemoveJob(jobID)
	return job, nil
}

// BuildServicesFromFolder reads a folder with services, builds images for
// these services, and adds all the necessary information to the database
func (p *Pier) BuildServicesFromFolder(servicesFolder string) error {
	log.Println("Building services from: " + servicesFolder)
	files, err := ioutil.ReadDir(servicesFolder)
	if err != nil {
		return def.Err(err, "cannot read folder %s", servicesFolder)
	}
	for _, f := range files {
		if f.IsDir() {
			log.Print("building service: " + f.Name())
			img, err := p.docker.client.BuildImage(filepath.Join(servicesFolder, f.Name()))

			if err != nil {
				log.Print("  failed to create service: ", err)
			} else {
				log.Print("  service has been created")

				err = p.docker.client.TagImage(string(img.ID), f.Name(), "latest")
				if err != nil {
					log.Print("  could not tag service")
				}

				img, err = p.docker.client.InspectImage(img.ID)
				if err != nil {
					log.Print("  failed to inspect image: ", err)
				}

				err = p.db.AddService(NewServiceFromImage(img))
				if err != nil {
					log.Print("  failed to add service to the database: ", err)
				}
			}
		}
	}
	return err
}

// ImportImage installs a docker tar file as a docker image
func (p *Pier) ImportImage(imageFilePath string) (db.Service, error) {
	imageID, err := p.docker.client.ImportImageFromTar(imageFilePath)
	if err != nil {
		return db.Service{}, def.Err(err, "docker ImportImage failed")
	}

	image, err := p.docker.client.InspectImage(imageID)

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
		Cmd:     image.Cmd,
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
