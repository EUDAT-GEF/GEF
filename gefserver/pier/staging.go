package pier

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"path/filepath"

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
func (p *Pier) DownStreamContainerFile(volumeID string, fileLocation string, w http.ResponseWriter) error {
	// Copy the file from the volume to a new container
	binds := []dckr.VolBind{
		dckr.NewVolBind(dckr.VolumeID(volumeID), "/root/volume", false),
	}
	containerID, _, err := p.docker.client.StartImage(
		string(p.docker.copyFromVolume.id),
		p.docker.copyFromVolume.repoTag,
		[]string{
			p.docker.copyFromVolume.cmd[0],
			filepath.Join("/root/volume/", fileLocation),
			"/root",
		},
		binds,
		p.docker.limits,
		p.docker.timeouts.Preparation,
		p.docker.timeouts.FileDownload)

	if err != nil {
		return def.Err(err, "copying files from the volume to the container failed")
	}

	// Stream the file from the container
	tarStream, err := p.docker.client.GetTarStream(string(containerID), fileLocation)
	if err != nil {
		return def.Err(err, "GetTarStream failed")
	}

	tarBallReader := tar.NewReader(tarStream)
	header, err := tarBallReader.Next()
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

// ListFiles exported
func (p *Pier) ListFiles(volumeID db.VolumeID, filePath string) ([]VolumeItem, error) {
	var volumeFileList []VolumeItem
	if string(volumeID) == "" {
		return volumeFileList, def.Err(nil, "volume name has not been specified")
	}

	// Bind the container with the volume
	volumesToMount := []dckr.VolBind{
		dckr.NewVolBind(dckr.VolumeID(volumeID), "/root/volume", false),
	}

	// Execute our image (it should produce a JSON file with the list of files)
	containerID, _, err := p.docker.client.StartImage(
		string(p.docker.fileList.id),
		p.docker.fileList.repoTag,
		[]string{
			p.docker.fileList.cmd[0], filePath, "r",
		},
		volumesToMount,
		p.docker.limits,
		p.docker.timeouts.Preparation,
		p.docker.timeouts.VolumeInspection)

	if err != nil {
		return volumeFileList, def.Err(err, "running image failed")
	}


	// Stop but do not remove the container
	_, err = p.docker.client.WaitContainerOrSwarmService(string(containerID), false)
	if err != nil {
		return volumeFileList, def.Err(err, "waiting for container to end failed")
	}

	// Reading the JSON file
	volumeFileList, err = p.readJSON(string(containerID), "/root/_filelist.json")
	if err != nil {
		return volumeFileList, def.Err(err, "readJson failed")
	}

	// Remove a container/swarm service (it was stopped earlier)
	_, err = p.docker.client.WaitContainerOrSwarmService(string(containerID), true)
	if err != nil {
		return volumeFileList, def.Err(err, "waiting for container to end failed")
	}

	return volumeFileList, err
}

// readJSON reads a JSON file with the list of files (in a volume) from a container
func (p *Pier) readJSON(containerID string, filePath string) ([]VolumeItem, error) {
	var volumeFileList []VolumeItem
	tarStream, err := p.docker.client.GetTarStream(containerID, filePath)
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
