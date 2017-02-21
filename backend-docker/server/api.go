package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier"
)

const (
	// Version defines the api version
	Version = "0.2.0"
)

const apiRootPath = "/api"

const (
	buildsTmpDir = "builds"
)

// Server is a master struct for serving HTTP API requests
type Server struct {
	Server http.Server
	pier   *pier.Pier
	tmpDir string
}

// NewServer creates a new Server
func NewServer(cfg def.ServerConfig, pier *pier.Pier, tmpDir string) (*Server, error) {
	tmpDir, err := def.MakeTmpDir(tmpDir)
	if err != nil {
		return nil, def.Err(err, "creating temporary directory failed")
	}

	server := &Server{
		Server: http.Server{
			Addr: cfg.Address,
			// timeouts seem to trigger even after a correct read
			// ReadTimeout: 	cfg.ReadTimeoutSecs * time.Second,
			// WriteTimeout: 	cfg.WriteTimeoutSecs * time.Second,
		},
		pier:   pier,
		tmpDir: tmpDir,
	}

	routes := map[string]func(http.ResponseWriter, *http.Request){
		"GET /":     server.infoHandler,
		"GET /info": server.infoHandler,

		"POST /builds":           server.newBuildImageHandler,
		"POST /builds/{buildID}": server.buildImageHandler,

		"GET /services":             server.listServicesHandler,
		"GET /services/{serviceID}": server.inspectServiceHandler,

		"POST /jobs":        server.executeServiceHandler,
		"GET /jobs":         server.listJobsHandler,
		"GET /jobs/{jobID}": server.inspectJobHandler,

		"GET /volumes/{volumeID}/{path:.*}": server.volumeContentHandler,
	}

	router := mux.NewRouter()
	apirouter := router.PathPrefix(apiRootPath).Subrouter()

	for mp, handler := range routes {
		methodPath := strings.SplitN(mp, " ", 2)
		apirouter.HandleFunc(methodPath[1], handler).Methods(methodPath[0])
	}

	server.Server.Handler = router
	return server, nil
}

// Start starts a new http listener
func (s *Server) Start() error {
	return s.Server.ListenAndServe()
}

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	Response{w}.Ok(jmap("version", Version))
}

func (s *Server) newBuildImageHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	_, buildID, err := def.NewRandomTmpDir(s.tmpDir, buildsTmpDir)
	if err != nil {
		Response{w}.ServerError("cannot create tmp subdir", err)
		return
	}
	loc := urljoin(r, buildID)
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

	service, err := s.pier.BuildService(buildDir)
	if err != nil {
		Response{w}.ServerError("build service failed: ", err)
		return
	}

	Response{w}.Ok(jmap("Service", service))
}

func (s *Server) listServicesHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	services := s.pier.ListServices()
	Response{w}.Ok(jmap("Services", services))
}

func (s *Server) inspectServiceHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	vars := mux.Vars(r)
	service, err := s.pier.GetService(pier.ServiceID(vars["serviceID"]))
	if err != nil {
		Response{w}.ClientError("cannot get service", err)
		return
	}
	Response{w}.Ok(jmap("Service", service))
}

func (s *Server) executeServiceHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	serviceID := r.FormValue("serviceID")
	if serviceID == "" {
		vars := mux.Vars(r)
		serviceID = vars["serviceID"]
	}
	logParam("serviceID", serviceID)

	input := r.FormValue("pid")
	if input == "" {
		vars := mux.Vars(r)
		input = vars["pid"]
	}
	logParam("pid", input)

	if serviceID == "" {
		Response{w}.ServerNewError("execute docker image: serviceID required")
		return
	}
	if input == "" {
		Response{w}.ServerNewError("execute docker image: pid required")
		return
	}

	service, err := s.pier.GetService(pier.ServiceID(serviceID))
	if err != nil {
		Response{w}.ClientError("cannot get service", err)
		return
	}

	job, err := s.pier.RunService(service, input)
	if err != nil {
		Response{w}.ServerError("cannot read the reqested file from the archive", err)
		return
	}

	loc := urljoin(r, string(job.ID))
	Response{w}.Location(loc).Created(jmap("Location", loc, "jobID", job.ID))
}

func (s *Server) listJobsHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	jobs := s.pier.ListJobs()
	Response{w}.Ok(jmap("Jobs", jobs))
}

func (s *Server) inspectJobHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	vars := mux.Vars(r)
	job, err := s.pier.GetJob(pier.JobID(vars["jobID"]))
	if err != nil {
		Response{w}.ClientError("cannot get job", err)
		return
	}
	Response{w}.Ok(jmap("Job", job))
}

func (s *Server) volumeContentHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	vars := mux.Vars(r)
	fileLocation := vars["path"]
	_, hasContent := r.URL.Query()["content"]
	fileName := filepath.Base(fileLocation)

	if hasContent { // Download a file from a volume
		err := s.pier.DownStreamContainerFile(vars["volumeID"], filepath.Join("/root/volume/", fileLocation), w)
		if err != nil {
			Response{w}.ServerError("downloading volume files failed", err)
		}

		Response{w}.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		Response{w}.Header().Set("Content-Disposition", "attachment; filename="+fileName)

	} else { // Return of list of files in a specific location in a volume
		volumeFiles, err := s.pier.ListFiles(pier.VolumeID(vars["volumeID"]), fileLocation)
		if err != nil {
			Response{w}.ServerError("streaming container files failed", err)
		}
		json.NewEncoder(w).Encode(volumeFiles)
	}
}

func (s *Server) buildVolumeHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	vars := mux.Vars(r)
	buildID := vars["buildID"]
	buildDir := filepath.Join(s.tmpDir, buildsTmpDir, buildID)
	s.pier.BuildVolume(buildDir)
}
