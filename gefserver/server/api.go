package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"encoding/json"

	"github.com/EUDAT-GEF/GEF/gefserver/db"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"github.com/EUDAT-GEF/GEF/gefserver/pier"
	"github.com/gorilla/mux"
)

const (
	// ServiceName is used for HTTP API
	ServiceName = "GEF"
	// Version defines the api version
	Version = "0.3.0"
)

const apiRootPath = "/api"
const wuiRootPath = "/wui"

const (
	buildsTmpDir = "builds"
)

// Server is a master struct for serving HTTP API requests
type Server struct {
	Server                 http.Server
	TLSCertificateFilePath string
	TLSKeyFilePath         string
	pier                   *pier.Pier
	db                     *db.Db
	tmpDir                 string
	administration         def.AdminConfig
}

// NewServer creates a new Server
func NewServer(cfg def.ServerConfig, pier *pier.Pier, tmpDir string, database *db.Db) (*Server, error) {
	tmpDir, err := def.MakeTmpDir(tmpDir)
	if err != nil {
		return nil, def.Err(err, "creating temporary directory failed")
	}

	server := &Server{
		Server: http.Server{
			Addr:         cfg.Address,
			ReadTimeout:  time.Duration(cfg.ReadTimeoutSecs) * time.Second,
			WriteTimeout: time.Duration(cfg.WriteTimeoutSecs) * time.Second,
		},
		TLSCertificateFilePath: cfg.TLSCertificateFilePath,
		TLSKeyFilePath:         cfg.TLSKeyFilePath,
		pier:                   pier,
		db:                     database,
		tmpDir:                 tmpDir,
		administration:         cfg.Administration,
	}

	routes := []struct {
		route       string
		handler     func(http.ResponseWriter, *http.Request)
		description string
	}{
		{"GET /", server.infoHandler, "misc"},
		{"GET /info", server.infoHandler, "misc"},

		{"GET /user", server.userHandler, "user discovery"},

		{"POST /user/tokens", server.newTokenHandler, "access management"},
		{"GET /user/tokens", server.listTokenHandler, "access discovery"},
		{"DELETE /user/tokens/{tokenID}", server.removeTokenHandler, "access management"},

		{"GET /roles", server.listRolesHandler, "access discovery"},
		{"GET /roles/{roleID}", server.listRoleUsersHandler, "access discovery"},
		{"POST /roles/{roleID}", server.newRoleUserHandler, "access management"},
		{"DELETE /roles/{roleID}/{userID}", server.removeRoleUserHandler, "access management"},

		{"POST /builds", server.newBuildImageHandler, "service deployment"},
		{"POST /builds/{buildID}", server.buildImageHandler, "service deployment"},

		{"GET /services", server.listServicesHandler, "service discovery"},
		{"GET /services/{serviceID}", server.inspectServiceHandler, "service discovery"},
		{"PUT /services/{serviceID}", server.editServiceHandler, "service modification"},
		{"DELETE /services/{serviceID}", server.removeServiceHandler, "service removal"},

		{"POST /jobs", server.executeServiceHandler, "data analysis"},
		{"GET /jobs", server.listJobsHandler, "data discovery"},
		{"GET /jobs/{jobID}", server.inspectJobHandler, "data discovery"},
		{"DELETE /jobs/{jobID}", server.removeJobHandler, "data cleanup"},

		{"GET /volumes/{volumeID}/{path:.*}", server.volumeContentHandler, "data retrieval"},
	}

	router := mux.NewRouter()
	apirouter := router.PathPrefix(apiRootPath).Subrouter()
	for _, hdl := range routes {
		methodPath := strings.SplitN(hdl.route, " ", 2)
		apirouter.HandleFunc(methodPath[1], server.decorate(hdl.handler, hdl.description)).Methods(methodPath[0])
	}
	wuirouter := router.PathPrefix(wuiRootPath).Subrouter()
	{
		wuirouter.HandleFunc("/login", server.decorate(server.oauthLoginHandler, "user login")).Methods("GET")
		wuirouter.HandleFunc("/b2access", server.oauthCallbackHandler).Methods("GET")
		wuirouter.HandleFunc("/logout", server.decorate(server.logoutHandler, "user logout")).Methods("GET")
	}
	router.PathPrefix("/").Handler(http.FileServer(singlePageAppDir("../webui/app/")))

	initB2Access(cfg.B2Access)

	server.Server.Handler = router

	return server, nil
}

type singlePageAppDir string

func (spad singlePageAppDir) Open(name string) (http.File, error) {
	f, err := http.Dir(spad).Open(name)
	if err != nil {
		log.Printf("serve file error: %#v\n", err)
		if _, isPathError := err.(*os.PathError); isPathError {
			log.Printf("    serving index.html instead")
			return http.Dir(spad).Open("/index.html")
		}
	}
	return f, err
}

// Start starts a new http listener
func (s *Server) Start() error {
	return s.Server.ListenAndServeTLS(s.TLSCertificateFilePath, s.TLSKeyFilePath)
}

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	Response{w}.Ok(jmap("ServiceName", ServiceName, "Version", Version,
		"ContactLink", s.administration.ContactLink))
}

func (s *Server) newBuildImageHandler(w http.ResponseWriter, r *http.Request) {
	allow, user := Authorization{s, w, r}.allowCreateBuild()
	if user == nil || !allow {
		return
	}

	_, buildID, err := def.NewRandomTmpDir(s.tmpDir, buildsTmpDir)
	if err != nil {
		Response{w}.ServerError("cannot create tmp subdir", err)
		return
	}
	loc, err := urljoin(r, buildID)
	if err != nil {
		Response{w}.ServerError("urljoin error", err)
		return
	}
	Response{w}.Location(loc).Created(jmap("Location", loc, "buildID", buildID))
}

func (s *Server) buildImageHandler(w http.ResponseWriter, r *http.Request) {
	allow, user := Authorization{s, w, r}.allowUploadIntoBuild()
	if user == nil || !allow {
		return
	}

	vars := mux.Vars(r)
	buildID := vars["buildID"]
	buildDir := filepath.Join(s.tmpDir, buildsTmpDir, buildID)

	mr, err := r.MultipartReader()
	if err != nil {
		Response{w}.ServerError("while getting multipart reader ", err)
		return
	}

	var service db.Service

	foundImageFileName := ""
	tarFileFound := false
	dockerFileFound := false
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if part.FileName() == "" {
			continue
		}

		log.Println("\tupload file " + part.FileName())
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

		if strings.HasSuffix(strings.ToLower(part.FileName()), ".tar") || strings.HasSuffix(strings.ToLower(part.FileName()), ".tar.gz") {
			tarFileFound = true
			foundImageFileName = part.FileName()
		}

		if strings.ToLower(part.FileName()) == "dockerfile" {
			dockerFileFound = true
		}

	}

	// Building an image from a Dockerfile
	if dockerFileFound {
		if _, err := os.Stat(filepath.Join(buildDir, "Dockerfile")); os.IsNotExist(err) {
			Response{w}.ServerError("no Dockerfile to build new image ", err)
			return
		}

		service, err = s.pier.BuildService(user.ID, buildDir)
		if err != nil {
			Response{w}.ServerError("build service failed: ", err)
			return
		}
	} else {
		// Importing an existing image from a tar archive
		if tarFileFound {
			log.Println("Docker image file has been detected, trying to import")
			log.Println(filepath.Join(buildDir, foundImageFileName))
			service, err = s.pier.ImportImage(user.ID, filepath.Join(buildDir, foundImageFileName))
			if err != nil {
				Response{w}.ServerError("while importing a Docker image file ", err)
				return
			}

			log.Println("Docker image has been imported")
		} else {
			Response{w}.ServerNewError("there is neither Dockerfile nor Tar archive")
			return
		}
	}

	Response{w}.Ok(jmap("Service", service))
}

func (s *Server) listServicesHandler(w http.ResponseWriter, r *http.Request) {
	allow, _ := Authorization{s, w, r}.allowListServices()
	if !allow {
		return
	}
	services, err := s.db.ListServices()
	if err != nil {
		Response{w}.ClientError("cannot get services", err)
		return
	}
	Response{w}.Ok(jmap("Services", services))
}

func (s *Server) inspectServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := db.ServiceID(vars["serviceID"])

	allow, _ := Authorization{s, w, r}.allowInspectService(serviceID)
	if !allow {
		return
	}

	service, err := s.db.GetService(serviceID)
	if err != nil {
		Response{w}.ClientError("cannot get service", err)
		return
	}
	Response{w}.Ok(jmap("Service", service))
}

func (s *Server) editServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := db.ServiceID(vars["serviceID"])

	allow, user := Authorization{s, w, r}.allowEditService(serviceID)
	if !allow {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var service db.Service
	err := decoder.Decode(&service)
	if err != nil {
		Response{w}.ClientError("cannot get service from JSON", err)
		return
	}
	defer r.Body.Close()

	if serviceID != service.ID {
		Response{w}.ServerNewError("update service: ID mismatch")
		return
	}

	err = s.db.RemoveService(user.ID, service.ID)
	if err != nil {
		Response{w}.ClientError("cannot remove service", err)
		return
	}

	err = s.db.AddService(user.ID, service)
	if err != nil {
		Response{w}.ClientError("cannot add service", err)
		return
	}

	Response{w}.Ok(jmap("Service", service))
}

func (s *Server) removeServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := db.ServiceID(vars["serviceID"])

	allow, user := Authorization{s, w, r}.allowRemoveService(serviceID)
	if !allow {
		return
	}

	service, err := s.db.GetService(serviceID)
	if err != nil {
		Response{w}.ClientError("cannot find service", err)
		return
	}

	err = s.db.RemoveService(user.ID, serviceID)
	if err != nil {
		Response{w}.ClientError("cannot remove service", err)
		return
	}

	service.Deleted = true

	err = s.db.AddService(user.ID, service)
	if err != nil {
		Response{w}.ClientError("cannot add service", err)
		return
	}

	Response{w}.Ok(jmap("Service", service))
}

func (s *Server) executeServiceHandler(w http.ResponseWriter, r *http.Request) {
	serviceID := r.FormValue("serviceID")
	if serviceID == "" {
		vars := mux.Vars(r)
		serviceID = vars["serviceID"]
	}
	logParam("serviceID", serviceID)

	allow, user := Authorization{s, w, r}.allowCreateJob()
	if !allow {
		return
	}

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

	service, err := s.db.GetService(db.ServiceID(serviceID))
	if err != nil {
		Response{w}.ClientError("cannot get service", err)
		return
	}

	job, err := s.pier.RunService(user.ID, service.ID, input)
	if err != nil {
		Response{w}.ServerError("cannot read the reqested file from the archive", err)
		return
	}

	loc, err := urljoin(r, string(job.ID))
	if err != nil {
		Response{w}.ServerError("urljoin error", err)
		return
	}
	Response{w}.Location(loc).Created(jmap("Location", loc, "jobID", job.ID))
}

func (s *Server) listJobsHandler(w http.ResponseWriter, r *http.Request) {
	allow, _ := Authorization{s, w, r}.allowListJobs()
	if !allow {
		return
	}

	jobs, err := s.db.ListJobs()

	if err != nil {
		Response{w}.ClientError("cannot get jobs", err)
		return
	}
	Response{w}.Ok(jmap("Jobs", jobs))
}

func (s *Server) inspectJobHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := db.JobID(vars["jobID"])
	allow, _ := Authorization{s, w, r}.allowInspectJob(jobID)
	if !allow {
		return
	}

	job, err := s.db.GetJob(jobID)
	if err != nil {
		Response{w}.ClientError("cannot get job", err)
		return
	}
	Response{w}.Ok(jmap("Job", job))
}

func (s *Server) removeJobHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := db.JobID(vars["jobID"])
	allow, user := Authorization{s, w, r}.allowRemoveJob(jobID)
	if !allow {
		return
	}

	job, err := s.db.GetJob(jobID)
	if err != nil {
		Response{w}.ClientError(err.Error(), err)
		return
	}

	err = s.db.RemoveJob(user.ID, jobID)
	if err != nil {
		Response{w}.ClientError(err.Error(), err)
		return
	}
	Response{w}.Ok(jmap("Job", job))
}

func (s *Server) volumeContentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volumeID := vars["volumeID"]
	job, err := s.db.GetJobOwningVolume(volumeID)
	if err != nil {
		Response{w}.ServerError("downloading volume files failed", err)
		return
	}
	allow, _ := Authorization{s, w, r}.allowGetJobData(job.ID)
	if !allow {
		return
	}

	fileLocation := vars["path"]
	_, hasContent := r.URL.Query()["content"]
	fileName := filepath.Base(fileLocation)

	if hasContent { // Download a file from a volume
		err := s.pier.DownStreamContainerFile(vars["volumeID"], filepath.Join("/root/volume/", fileLocation), w)
		if err != nil {
			Response{w}.ServerError("downloading volume files failed", err)
			return
		}

		Response{w}.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		Response{w}.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	} else { // Return of list of files in a specific location in a volume
		volumeFiles, err := s.pier.ListFiles(db.VolumeID(vars["volumeID"]), fileLocation)
		if err != nil {
			Response{w}.ServerError("streaming container files failed", err)
			return
		}
		Response{w}.Ok(jmap("volumeID", vars["volumeID"], "volumeContent", volumeFiles))
	}
}
