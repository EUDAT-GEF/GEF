package pier

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"path/filepath"

	"github.com/EUDAT-GEF/GEF/backend-docker/db"
	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
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

const volumeFileListName = "volume-filelist"
const copyFromVolumeName = "copy-from-volume"

// DownStreamContainerFile exported
func (p *Pier) DownStreamContainerFile(volumeID string, fileLocation string, w http.ResponseWriter) error {
	// Copy the file from the volume to a new container
	binds := []dckr.VolBind{
		dckr.NewVolBind(dckr.VolumeID(volumeID), "/root/volume", false),
	}
	containerID, _, err := p.docker.StartImage(dckr.ImageID(copyFromVolumeName), []string{filepath.Join("/root/volume/", fileLocation), "/root"}, binds)

	if err != nil {
		return def.Err(err, "copying files from the volume to the container failed")
	}

	// Stream the file from the container
	tarStream, err := p.docker.GetTarStream(string(containerID), fileLocation)
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
	containerID, consoleOutput, err := p.docker.StartImage(dckr.ImageID(volumeFileListName), []string{filePath, "r"}, volumesToMount)

	if err != nil {
		return volumeFileList, def.Err(err, "running image failed")
	}

	// Reading the JSON file
	volumeFileList, err = p.readJSON(string(containerID), "/root/_filelist.json")
	if err != nil {
		return volumeFileList, def.Err(err, "readJson failed")
	}

	// Killing the container
	_, _, err = p.docker.WaitContainer(containerID, consoleOutput, true)
	if err != nil {
		return volumeFileList, def.Err(err, "waiting for container to end failed")
	}

	return volumeFileList, err
}

// readJSON reads a JSON file with the list of files (in a volume) from a container
func (p *Pier) readJSON(containerID string, filePath string) ([]VolumeItem, error) {
	var volumeFileList []VolumeItem
	tarStream, err := p.docker.GetTarStream(containerID, filePath)
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
