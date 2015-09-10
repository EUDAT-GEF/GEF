package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gefx/gef-docker/dckr"

	"code.google.com/p/go-uuid/uuid"
	"github.com/gorilla/mux"
)

const (
	Version       = "0.1.0"
	apiRootPath   = "/api"
	imagesAPIPath = "/images"
	jobsAPIPath   = "/jobs"
	buildsAPIPath = "/builds"
	buildsTmpDir  = "builds"
	tmpDirDefault = "gefdocker"
	tmpDirPerm    = 0700
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

	apirouter.HandleFunc(buildsAPIPath, server.newBuildHandler).Methods("POST")
	apirouter.HandleFunc(buildsAPIPath+"/{buildID}", server.buildHandler).Methods("POST")

	apirouter.HandleFunc(imagesAPIPath, server.listServicesHandler).Methods("GET")
	apirouter.HandleFunc(imagesAPIPath+"/{imageID}", server.inspectServiceHandler).Methods("GET")

	apirouter.HandleFunc(jobsAPIPath, server.listJobsHandler).Methods("GET")
	apirouter.HandleFunc(jobsAPIPath, server.executeServiceHandler).Methods("POST")
	apirouter.HandleFunc(jobsAPIPath+"/{jobID}/", server.inspectJobHandler).Methods("GET")

	server.server.Handler = router
	return server
}

// Start starts a new http listener
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	Response{w}.Ok(jmap("version", version))
}

func (s Server) newBuildHandler(w http.ResponseWriter, r *http.Request) {
	buildID := string(uuid.NewRandom())
	Response{w}.Location(buildID).Ok(jmap("Location", buildID))
}

func (s Server) buildHandler(w http.ResponseWriter, r *http.Request) {
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
	Response{w}.Ok(jmap("Image", image))
}

func (s Server) listServicesHandler(w http.ResponseWriter, r *http.Request) {
	imgIDs, err := s.docker.ListImages()
	if err != nil {
		Response{w}.ServerError("list of docker images: ", err)
		return
	}
	Response{w}.Ok(jmap("ImageIDs", imgIDs))
}

func (s Server) inspectServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := dckr.ImageID(vars["imageID"])
	image, err := s.docker.InspectImage(imageID)
	if err != nil {
		Response{w}.ServerError("inspect docker image: ", err)
		return
	}
	srv := extractServiceInfo(image.Labels)
	Response{w}.Ok(jmap("Image", image, "Service", srv))
}

func (s Server) listJobsHandler(w http.ResponseWriter, r *http.Request) {
	containers, err := s.docker.ListContainers()
	if err != nil {
		Response{w}.ServerError("list all containers: ", err)
		return
	}
	Response{w}.Ok(jmap("Containers", containers))
}

func (s Server) executeServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := dckr.ImageID(vars["imageID"])
	containerID, err := s.docker.ExecuteImage(imageID)
	if err != nil {
		Response{w}.ServerError("execute docker image: ", err)
		return
	}
	Response{w}.Ok(jmap("ContainerID", containerID))
}

func (s Server) inspectJobHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contID := dckr.ContainerID(vars["jobID"])
	cont, err := s.docker.InspectContainer(contID)
	if err != nil {
		Response{w}.ServerError("inspect container: ", err)
		return
	}
	Response{w}.Ok(jmap("Container", cont))
}
