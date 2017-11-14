package pier

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"path/filepath"

	"log"

	"github.com/EUDAT-GEF/GEF/gefserver/db"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"github.com/EUDAT-GEF/GEF/gefserver/pier/internal/dckr"
)

// VolumeItem describes a folder content
type VolumeItem struct {
	Name       string       `json:"name"`
	Size       int64        `json:"size"`
	Modified   time.Time    `json:"modified"`
	IsFolder   bool         `json:"isFolder"`
	Path       string       `json:"path"`
	FolderTree []VolumeItem `json:"folderTree"`
}

// DownStreamContainerFile exported
func (p *Pier) DownStreamContainerFile(volumeID string, fileLocation string, limits def.LimitConfig, timeouts def.TimeoutConfig, w http.ResponseWriter) error {
	job, err := p.db.GetJobOwningVolume(string(volumeID))
	if err != nil {
		return err
	}
	docker, found := p.docker[job.ConnectionID]
	if !found {
		return def.Err(nil, "Cannot find docker connection")
	}

	// Copy the file from the volume to a new container
	binds := []dckr.VolBind{
		dckr.NewVolBind(dckr.VolumeID(volumeID), "/root/volume", false),
	}
	containerID, swarmServiceID, _, err := docker.client.StartImageOrSwarmService(
		string(docker.copyToAndFromVolume.id),
		docker.copyToAndFromVolume.repoTag,
		[]string{
			docker.copyToAndFromVolume.cmd[0],
			filepath.Join("/root/volume/", fileLocation),
			"/root",
		},
		binds,
		limits,
		timeouts)

	if err != nil {
		return def.Err(err, "copying files from the volume to the container failed")
	}

	// Stream the file from the container
	tarStream, err := docker.client.GetTarStream(string(containerID), fileLocation)

	if err != nil {
		return def.Err(err, "GetTarStream failed")
	}

	tarBallReader := tar.NewReader(tarStream)
	header, err := tarBallReader.Next()
	defer func() {
		err := docker.client.TerminateContainerOrSwarmService(string(containerID), swarmServiceID)
		if err != nil {
			log.Println("error while forcefully removing container in DownStreamContainerFile", err)
		}
	}()
	if err != nil {
		return def.Err(err, "reading tarball failed")
	}
	filename := header.Name

	if header.Typeflag == tar.TypeReg {
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
		w.Header().Set("Content-Type", "application/octet-stream")
		io.Copy(w, tarBallReader)
	} else {
		http.Error(w, "Error", http.StatusInternalServerError)
		return errors.New("internal error while reading tarball")
	}
	return nil
}

// UploadFileIntoVolume exported
func (p *Pier) UploadFileIntoVolume(volumeID string, srcFileLocation string, dstFileName string, limits def.LimitConfig, timeouts def.TimeoutConfig) error {
	job, err := p.db.GetJobOwningVolume(string(volumeID))
	if err != nil {
		return err
	}
	docker, found := p.docker[job.ConnectionID]
	if !found {
		return def.Err(nil, "Cannot find docker connection")
	}

	// Copy the file from the volume to a new container
	binds := []dckr.VolBind{
		dckr.NewVolBind(dckr.VolumeID(volumeID), "/root/volume", false),
	}
	containerID, swarmServiceID, _, err := docker.client.StartImageOrSwarmService(
		string(docker.copyToAndFromVolume.id),
		docker.copyToAndFromVolume.repoTag,
		[]string{
			"ls",
		},
		binds,
		limits,
		timeouts)

	if err != nil {
		return def.Err(err, "data uploading container failed")
	}

	err = docker.client.UploadFile2Container(string(containerID), srcFileLocation, "/root/volume")
	if err != nil {
		return def.Err(err, "data uploading failed")
	}


	err = docker.client.TerminateContainerOrSwarmService(string(containerID), swarmServiceID)
	if err != nil {
		log.Println("error while forcefully removing container in UploadFileIntoVolume", err)
	}


	return nil
}

// ListFiles exported
func (p *Pier) ListFiles(volumeID db.VolumeID, filePath string, limits def.LimitConfig, timeouts def.TimeoutConfig) ([]VolumeItem, error) {
	job, err := p.db.GetJobOwningVolume(string(volumeID))
	if err != nil {
		return nil, err
	}
	docker, found := p.docker[job.ConnectionID]
	if !found {
		return nil, def.Err(nil, "Cannot find docker connection")
	}

	var volumeFileList []VolumeItem
	if string(volumeID) == "" {
		return volumeFileList, def.Err(nil, "volume name has not been specified")
	}

	// Bind the container with the volume
	volumesToMount := []dckr.VolBind{
		dckr.NewVolBind(dckr.VolumeID(volumeID), "/root/volume", false),
	}

	// Execute our image (it should produce a JSON file with the list of files)
	containerID, swarmServiceID, _, err := docker.client.StartImageOrSwarmService(
		string(docker.fileList.id),
		docker.fileList.repoTag,
		[]string{
			docker.fileList.cmd[0], filePath, "r",
		},
		volumesToMount,
		limits,
		timeouts)

	if err != nil {
		return volumeFileList, def.Err(err, "running image failed")
	}

	// Stop but do not remove the container
	_, err = docker.client.WaitContainerOrSwarmService(string(containerID))
	if err != nil {
		return volumeFileList, def.Err(err, "WaitContainerOrSwarmService failed")
	}

	// Reading the JSON file
	volumeFileList, err = p.readJSON(docker.client, string(containerID), "/root/_filelist.json")
	if err != nil {
		return volumeFileList, def.Err(err, "readJson failed")
	}

	// Remove a container/swarm service (it was stopped earlier)
	err = docker.client.TerminateContainerOrSwarmService(string(containerID), swarmServiceID)
	if err != nil {
		return volumeFileList, def.Err(err, "TerminateContainerOrSwarmService failed")
	}

	return volumeFileList, err
}

// readJSON reads a JSON file with the list of files (in a volume) from a container
func (p *Pier) readJSON(dockerClient dckr.Client, containerID string, filePath string) ([]VolumeItem, error) {
	var volumeFileList []VolumeItem
	tarStream, err := dockerClient.GetTarStream(containerID, filePath)
	if err != nil {
		return nil, def.Err(err, "GetTarStream(%s) failed", filePath)
	}

	tarBallReader := tar.NewReader(tarStream)
	_, err = tarBallReader.Next()
	if err != nil {
		return nil, def.Err(err, "tarball reader failed", filePath)
	}

	jsonParser := json.NewDecoder(tarBallReader)
	err = jsonParser.Decode(&volumeFileList)

	return volumeFileList, err
}
