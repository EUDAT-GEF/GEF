package pier

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
	"path/filepath"
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
func (p *Pier) ListFiles(volumeID VolumeID, filePath string) ([]VolumeItem, error) {
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
	_, _, _, err = p.docker.WaitContainer(containerID, consoleOutput, true)
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

////////

// BuildVolume exported
func (p *Pier) BuildVolume(pid string) {
	// TODO
	log.Println("BuildVolume: not implemented")
	// buildDir := filepath.Join(s.tmpDir, buildsTmpDir, buildID)

	// // STEP 1: Get a list of files from PID
	// // Temporary solution for the list of files
	// var pidList []string
	// pidList = append(pidList, "#!/bin/ash")
	// pidList = append(pidList, "wget https://b2share.eudat.eu/record/154/files/ISGC2014_022.pdf?version=1 -P /root/volume")
	// pidList = append(pidList, "wget https://b2share.eudat.eu/record/157/files/TenReasonsToSwitchFromMauiToMoab2012-01-05.pdf?version=1 -P /root/volume")
	// pidList = append(pidList, "ls -l /root/volume/")

	// // STEP 2: create a bash script that downloads those files
	// dScriptPath := filepath.Join(buildDir, "downloader.sh")
	// dScriptFile, err := os.Create(dScriptPath)
	// if err != nil {
	// 	Response{w}.ServerError("create script file:", err)
	// 	return
	// }
	// log.Println("Script was created")
	// _, err = dScriptFile.WriteString(strings.Join(pidList, "\n"))
	// if err != nil {
	// 	Response{w}.ServerError("write data into the script file", err)
	// 	return
	// }
	// dScriptFile.Sync()
	// log.Println("Wrote file list")

	// err = dScriptFile.Chmod(0777)
	// if err != nil {
	// 	Response{w}.ServerError("make downloading script executable:", err)
	// 	return
	// }
	// log.Println("Changed permissions")

	// // STEP 3: create an image that includes the script
	// var dockerFileContent []string
	// dockerFileContent = append(dockerFileContent, "FROM alpine:latest")
	// dockerFileContent = append(dockerFileContent, "RUN apk add --update --no-cache openssl openssl-dev ca-certificates")
	// dockerFileContent = append(dockerFileContent, "RUN mkdir /root/volume")
	// dockerFileContent = append(dockerFileContent, "ADD downloader.sh /root")
	// dockerFileContent = append(dockerFileContent, "CMD [\"/root/downloader.sh\"]")

	// dockerFilePath := filepath.Join(buildDir, "Dockerfile")
	// dockerFile, err := os.Create(dockerFilePath)
	// if err != nil {
	// 	Response{w}.ServerError("create script file:", err)
	// 	return
	// }
	// log.Println("Dockerfile was created")
	// _, err = dockerFile.WriteString(strings.Join(dockerFileContent, "\n"))
	// if err != nil {
	// 	Response{w}.ServerError("write data into the  Dockerfile", err)
	// 	return
	// }
	// dockerFile.Sync()
	// log.Println("Wrote Dockerfile content")

	// // STEP 4: create a new empty volume
	// volume, err := s.docker.BuildVolume(buildDir)
	// if err != nil {
	// 	Response{w}.ServerError("build docker volume:", err)
	// 	return
	// }
	// log.Println(volume.ID)
	// log.Println(buildDir)
	// log.Println("Volume was created")

	// image, err := s.docker.BuildImage(buildDir)
	// if err != nil {
	// 	Response{w}.ServerError("build docker image: ", err)
	// 	return
	// }
	// log.Println("Docker image was created")
	// log.Println(image.ID)

	// imageID := string(image.ID)

	// // STEP 5: run the image, as a result we get a volume with our files
	// volumesToMount := []string{
	// 	string(volume.ID) + ":/root/volume"}

	// containerID, err := s.docker.ExecuteImage(dckr.ImageID(imageID), volumesToMount)
	// if err != nil {
	// 	Response{w}.ServerError("execute docker image: ", err)
	// 	return
	// }
	// log.Println("Executed the image")

	// log.Println(containerID)

	// _, err = s.docker.WaitContainer(containerID, true)
	// if err != nil {
	// 	Response{w}.ServerError("removing the container: ", err)
	// 	return
	// }
	// log.Println("Container was removed")

	// err = s.docker.DeleteImage(imageID)
	// if err != nil {
	// 	Response{w}.ServerError("removing the image "+imageID+": ", err)
	// 	return
	// }
	// log.Println("Image was removed")

	// Response{w}.Ok(jmap("Volume", volume))
}
