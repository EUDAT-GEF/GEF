package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/EUDAT-GEF/GEF/backend-docker/dckr"

	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"time"
	"archive/tar"
	"github.com/fsouza/go-dockerclient"
)

const (
	// HarcodedIrodsMountPoint must be removed
	HarcodedIrodsMountPoint = "/data/GEF/datasets/"
	// HarcodedB2DropMountPoint to be removed
	HarcodedB2DropMountPoint = "/webdav"
)

const (
	// Version defines the api version
	Version = "0.1.1"
)

const (
	apiRootPath          = "/api"
	imagesAPIPath        = "/images"
	jobsAPIPath          = "/jobs"
	volumesAPIPath       = "/volumes"
	buildImagesAPIPath   = "/buildImages"
	buildVolumesAPIPath  = "/buildVolumes"
	inspectVolumeAPIPath = "/inspectVolume"
	retrieveFileAPIPath = "/retrieveFile"

	tmpDirDefault = "gefdocker"
	tmpDirPerm    = 0700
	buildsTmpDir  = "builds"

	sock = "/var/run/docker.sock"
	dockerAPIAddr = "http://localhost:4243"
)

// Config keeps the configuration options needed to make a Server
type Config struct {
	Address          string
	ReadTimeoutSecs  int
	WriteTimeoutSecs int

	// TmpDir is the directory to keep session files in
	// If the path is relative, it will be used as a subfolder of the system temporary directory
	TmpDir string
}

// Server is a master struct for serving HTTP API requests
type Server struct {
	server http.Server
	tmpDir string
	docker *dckr.Client
}

// Volume folder content
type VolumeItem struct {
	Name     string
	Size     int64
	IsFolder bool
	Modified time.Time
}

// NewServer creates a new Server
func NewServer(cfg Config, docker dckr.Client) *Server {
	tmpDir := cfg.TmpDir
	if tmpDir == "" {
		tmpDir = tmpDirDefault
	}
	if !filepath.IsAbs(tmpDir) {
		tmpDir = filepath.Join(os.TempDir(), tmpDir)
	}
	if err := os.MkdirAll(tmpDir, os.FileMode(tmpDirPerm)); err != nil {
		log.Println("ERROR: cannot create temporary directory: ", err)
		return nil
	}

	server := &Server{
		server: http.Server{
			Addr: cfg.Address,
			// timeouts seem to trigger even after a correct read
			// ReadTimeout: 	cfg.ReadTimeoutSecs * time.Second,
			// WriteTimeout: 	cfg.WriteTimeoutSecs * time.Second,
		},
		tmpDir: tmpDir,
		docker: &docker,
	}

	router := mux.NewRouter()
	apirouter := router.PathPrefix(apiRootPath).Subrouter()

	apirouter.HandleFunc("/", server.infoHandler).Methods("GET")
	apirouter.HandleFunc("/info", server.infoHandler).Methods("GET")

	apirouter.HandleFunc(buildImagesAPIPath, server.newBuildImageHandler).Methods("POST")
	apirouter.HandleFunc(buildImagesAPIPath+"/{buildID}", server.buildImageHandler).Methods("POST")

	apirouter.HandleFunc(buildVolumesAPIPath, server.newBuildVolumeHandler).Methods("POST")
	apirouter.HandleFunc(buildVolumesAPIPath+"/{buildID}", server.buildVolumeHandler).Methods("POST")

	apirouter.HandleFunc(imagesAPIPath, server.listServicesHandler).Methods("GET")
	apirouter.HandleFunc(imagesAPIPath+"/{imageID}", server.inspectServiceHandler).Methods("GET")

	apirouter.HandleFunc(volumesAPIPath, server.listVolumesHandler).Methods("GET")
	//apirouter.HandleFunc(volumesAPIPath+"/{volumeID}", server.inspectVolumeHandler).Methods("GET")
	apirouter.HandleFunc(inspectVolumeAPIPath+"/{volumeID}", server.inspectVolumeHandler).Methods("GET")
	apirouter.HandleFunc(retrieveFileAPIPath+"/{containerID}/{filePath}", server.retrieveFileHandler).Methods("GET")

	apirouter.HandleFunc(jobsAPIPath, server.executeServiceHandler).Methods("POST")
	apirouter.HandleFunc(jobsAPIPath, server.listJobsHandler).Methods("GET")
	apirouter.HandleFunc(jobsAPIPath+"/{jobID}", server.inspectJobHandler).Methods("GET")

	server.server.Handler = router
	return server
}



// Start starts a new http listener
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	Response{w}.Ok(jmap("version", Version))
}

func (s *Server) retrieveFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerID := string(dckr.ImageID(vars["containerID"]))
	filePath := string(dckr.ImageID(vars["filePath"]))

	tarStream, err := s.docker.GetTarStream(containerID, filePath)

	if err == nil {
		tarBallReader := tar.NewReader(tarStream)
		header, err := tarBallReader.Next()
		if err != nil {
			Response{w}.ServerError("cannot read the reqested file from the archive", err)
			return
		}
		filename := header.Name

		if header.Typeflag == tar.TypeReg {
			if err == nil {
				w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
				w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
				io.Copy(w, tarBallReader)
			} else {
				http.Error(w, "Cannot read the content of the requested file: " + err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

	} else {
		http.Error(w, "Cannot access Docker remote API: " + err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Server) readJSON(containerID string, filePath string) ([]VolumeItem, error) {
	var volumeFileList []VolumeItem
	tarStream, err := s.docker.GetTarStream(containerID, filePath)

	if err == nil {
		tarBallReader := tar.NewReader(tarStream)
		_, err = tarBallReader.Next()
		if err == nil {
			jsonParser := json.NewDecoder(tarBallReader)
			err = jsonParser.Decode(volumeFileList)
		}
	}

	return volumeFileList, err
}

func (s *Server) inspectVolumeHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	vars := mux.Vars(r)
	volId := string(dckr.ImageID(vars["volumeID"]))

	imageID := "d755f01f457926e662451aa753a627c84e6e1f76d517d0982000795630673182"

	// Bind the container with the volume
	volumesToMount := []string{
		volId + ":/root/volume"}

	// Execute our image (it should produce a JSON file with the list of files)
	containerID, err := s.docker.ExecuteImage(dckr.ImageID(imageID), volumesToMount)
	if err != nil {
		Response{w}.ServerError("execute docker image: ", err)
		return
	}

	// Reading the JSON file
	volumeFileList, err := s.readJSON(string(containerID), "/root/_filelist.json")
	if err != nil {
		Response{w}.ServerError("reading the list of files in a volume: ", err)
		return
	}

	// Killing the container
	_, err = s.docker.WaitContainer(containerID, true)
	if err != nil {
		Response{w}.ServerError("removing the container: ", err)
		return
	}

	json.NewEncoder(w).Encode(volumeFileList)

}

func (s *Server) newBuildImageHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	buildID := uuid.NewRandom().String()
	buildDir := filepath.Join(s.tmpDir, buildsTmpDir, buildID)
	if err := os.MkdirAll(buildDir, os.FileMode(tmpDirPerm)); err != nil {
		Response{w}.ServerError("cannot create temporary directory", err)
		return
	}
	loc := apiRootPath + buildImagesAPIPath + "/" + buildID

	Response{w}.Location(loc).Created(jmap("Location", loc, "buildID", buildID))
}

func (s *Server) newBuildVolumeHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	buildID := uuid.NewRandom().String()
	buildDir := filepath.Join(s.tmpDir, buildsTmpDir, buildID)
	if err := os.MkdirAll(buildDir, os.FileMode(tmpDirPerm)); err != nil {
		Response{w}.ServerError("cannot create temporary directory", err)
		return
	}
	loc := apiRootPath + buildVolumesAPIPath + "/" + buildID
	Response{w}.Location(loc).Created(jmap("Location", loc, "buildID", buildID))
}

func (s *Server) buildImageHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	vars := mux.Vars(r)
	buildID := vars["buildID"]
	buildDir := filepath.Join(s.tmpDir, buildsTmpDir, buildID)

	mr, err := r.MultipartReader()
	if err != nil {
		Response{w}.ServerError("while getting multipart reader ", err)
		return
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if part.FileName() == "" {
			continue
		}
		dst, err := os.Create(filepath.Join(buildDir, part.FileName()))
		if err != nil {
			Response{w}.ServerError("while creating file to save file part ", err)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, part); err != nil {
			Response{w}.ServerError("while dumping file part ", err)
			return
		}
	}

	if _, err := os.Stat(filepath.Join(buildDir, "Dockerfile")); os.IsNotExist(err) {
		Response{w}.ServerError("no Dockerfile to build new image ", err)
		return
	}

	image, err := s.docker.BuildImage(buildDir)
	if err != nil {
		Response{w}.ServerError("build docker image: ", err)
		return
	}
	Response{w}.Ok(jmap("Image", image, "Service", extractServiceInfo(image)))
}

func (s *Server) buildVolumeHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	vars := mux.Vars(r)
	buildID := vars["buildID"]
	buildDir := filepath.Join(s.tmpDir, buildsTmpDir, buildID)

	mr, err := r.MultipartReader()
	if err != nil {
		Response{w}.ServerError("while getting multipart reader ", err)
		return
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if part.FileName() == "" {
			continue
		}
		dst, err := os.Create(filepath.Join(buildDir, part.FileName()))
		if err != nil {
			Response{w}.ServerError("while creating file to save file part ", err)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, part); err != nil {
			Response{w}.ServerError("while dumping file part ", err)
			return
		}
	}

	volume, err := s.docker.BuildVolume(buildDir)
	if err != nil {
		Response{w}.ServerError("build docker volume:", err)
		return
	}
	Response{w}.Ok(jmap("Volume", volume))
}

func (s *Server) listServicesHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	images, err := s.docker.ListImages()
	if err != nil {
		Response{w}.ServerError("list of docker images: ", err)
		return
	}
	services := make([]Service, len(images), len(images))
	for i, img := range images {
		// fmt.Println("list serv handler: ", img)
		services[i] = extractServiceInfo(img)
	}
	Response{w}.Ok(jmap("Images", images, "Services", services))
}

func (s *Server) inspectServiceHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	vars := mux.Vars(r)
	imageID := dckr.ImageID(vars["imageID"])
	image, err := s.docker.InspectImage(imageID)
	if err != nil {
		Response{w}.ServerError("inspect docker image: ", err)
		return
	}
	Response{w}.Ok(jmap("Image", image, "Service", extractServiceInfo(image)))
}

func (s *Server) executeServiceHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	imageID := r.FormValue("imageID")
	if imageID == "" {
		vars := mux.Vars(r)
		imageID = vars["imageID"]
	}
	if imageID == "" {
		Response{w}.ServerNewError("execute docker image: imageID required")
		return
	}
	logParam("imageID", imageID)

	image, err := s.docker.InspectImage(dckr.ImageID(imageID))
	if err != nil {
		Response{w}.ServerError("execute docker image: inspectImage: ", err)
	}
	binds, err := makeBinds(r, image)
	if err != nil {
		Response{w}.ServerError("execute docker image: binds: ", err)
		return
	}
	logParam("binds", strings.Join(binds, " : "))

	containerID, err := s.docker.ExecuteImage(dckr.ImageID(imageID), binds)
	fmt.Println("EXECUTE")
	fmt.Println((dckr.ImageID(imageID)))
	if err != nil {
		Response{w}.ServerError("execute docker image: ", err)
		return
	}
	loc := apiRootPath + jobsAPIPath + "/" + string(containerID)
	Response{w}.Location(loc).Created(jmap("Location", loc, "jobID", containerID))
}

// makeBinds construct volume:path binds
func makeBinds(r *http.Request, image dckr.Image) ([]string, error) {
	svc := extractServiceInfo(image)
	var binds []string
	for _, in := range svc.Input {
		volumeID := r.FormValue(in.ID)
		if volumeID == "" {
			return nil, fmt.Errorf("no bind volume for input port: %s", in.Name)
		}
		binds = append(binds, fmt.Sprintf("%s:%s:ro", volumeID, in.Path))
	}
	for _, out := range svc.Output {
		volumeID := r.FormValue(out.ID)
		if volumeID == "" {
			return nil, fmt.Errorf("no bind volume for output port: %s", out.Name)
		}
		binds = append(binds, fmt.Sprintf("%s:%s", volumeID, out.Path))
	}
	return binds, nil
}

func (s *Server) listJobsHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	containers, err := s.docker.ListContainers()
	if err != nil {
		Response{w}.ServerError("list all docker containers: ", err)
		return
	}
	jobs := make([]Job, len(containers), len(containers))
	for i, c := range containers {
		jobs[i] = makeJob(c)
	}
	Response{w}.Ok(jmap("Jobs", jobs))
}

func (s *Server) inspectJobHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	vars := mux.Vars(r)
	contID := dckr.ContainerID(vars["jobID"])
	cont, err := s.docker.InspectContainer(contID)
	if err != nil {
		Response{w}.ServerError("inspect docker container: ", err)
		return
	}
	Response{w}.Ok(jmap("Job", makeJob(cont)))
}

func (s *Server) listVolumesHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	vols, err := s.docker.ListVolumes()
	if err != nil {
		Response{w}.ServerError("list all docker named volumes: ", err)
		return
	}
	Response{w}.Ok(jmap("Volumes", vols))
}
