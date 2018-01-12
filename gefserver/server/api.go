package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	inputTmpDir  = "inputs"
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
	limits                 def.LimitConfig
	timeouts               def.TimeoutConfig
}

// NewServer creates a new Server
func NewServer(cfg def.Configuration, pier *pier.Pier, database *db.Db) (*Server, error) {
	tmpDir, err := def.MakeTmpDir(cfg.TmpDir)
	if err != nil {
		return nil, def.Err(err, "creating temporary directory failed")
	}

	server := &Server{
		Server: http.Server{
			Addr:         cfg.Server.Address,
			ReadTimeout:  time.Duration(cfg.Server.ReadTimeoutSecs) * time.Second,
			WriteTimeout: time.Duration(cfg.Server.WriteTimeoutSecs) * time.Second,
		},
		TLSCertificateFilePath: cfg.Server.TLSCertificateFilePath,
		TLSKeyFilePath:         cfg.Server.TLSKeyFilePath,
		pier:                   pier,
		db:                     database,
		tmpDir:                 tmpDir,
		administration:         cfg.Server.Administration,
		limits:                 cfg.Limits,
		timeouts:               cfg.Timeouts,
	}

	routes := []struct {
		route       string
		handler     func(http.ResponseWriter, *http.Request, environment)
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

		{"POST /builds", server.newBuildImageHandler, "build initialization"},
		{"POST /builds/{buildID}", server.startBuildImageHandler, "build start"},
		{"GET /builds/{buildID}", server.inspectBuildImageHandler, "build discovery"},

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

	initB2Access(cfg.Server.B2Access)

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

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request, e environment) {
	Response{w}.Ok(jmap("ServiceName", ServiceName, "Version", Version,
		"ContactLink", s.administration.ContactLink))
}

func (s *Server) newBuildImageHandler(w http.ResponseWriter, r *http.Request, e environment) {
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

func (s *Server) getConnectionIDParam(r *http.Request) (db.ConnectionID, error) {
	vars := mux.Vars(r)
	connectionIDString := vars["connectionID"]
	if connectionIDString == "" {
		return s.db.GetFirstConnectionID()
	}
	connectionID, err := strconv.Atoi(connectionIDString)
	if err != nil {
		return 0, def.Err(err, "bad connectionID parameter")
	}
	if connectionID < 1 {
		return 0, def.Err(nil, "connectionID parameter should be > 0")
	}
	return db.ConnectionID(connectionID), nil
}

func (s *Server) startBuildImageHandler(w http.ResponseWriter, r *http.Request, e environment) {
	allow, user := Authorization{s, w, r}.allowUploadIntoBuild()
	if user == nil || !allow {
		return
	}

	hasDockerfile := false
	tarArchiveName := ""
	vars := mux.Vars(r)
	buildID := vars["buildID"]
	buildDir := filepath.Join(s.tmpDir, buildsTmpDir, buildID)

	connectionID, err := s.getConnectionIDParam(r)
	if err != nil {
		Response{w}.ClientError("bad connectionID", err)
	}

	var newBuild = db.Build{
		ID:           buildID,
		ConnectionID: connectionID,
		Started:      time.Now(),
		Duration:     0,
		State: &db.BuildState{
			Status: "Image build has been initiated",
			Error:  "",
			Code:   -1,
		},
	}

	err = s.db.AddBuild(newBuild)
	if err != nil {
		err = s.db.SetBuildState(buildID, db.NewBuildStateError("Failed to add the new build to the database", 1))
		if err != nil {
			log.Println(err)
		}
		Response{w}.ServerError("while adding a new build ", err)
		return
	}

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
		err = s.db.SetBuildState(buildID, db.NewBuildStateOk("Uploading file "+part.FileName(), -1))
		if err != nil {
			log.Println(err)
		}

		fmt.Println("BUILD ID = " + buildID)

		log.Println("\tupload file " + part.FileName())
		dst, err := os.Create(filepath.Join(buildDir, part.FileName()))
		if err != nil {
			log.Print("while creating file to save file part ", err)
			err = s.db.SetBuildState(buildID, db.NewBuildStateError("Failed while creating file to save file part "+part.FileName(), 1))
			if err != nil {
				log.Println(err)
			}
			return
		}
		defer dst.Close()

		if _, err = io.Copy(dst, part); err != nil {
			log.Print("while dumping file part ", err)
			err = s.db.SetBuildState(buildID, db.NewBuildStateError("Failed while dumping file part "+part.FileName(), 1))
			if err != nil {
				log.Println(err)
			}
			return
		}

		if strings.HasSuffix(strings.ToLower(part.FileName()), ".tar") || strings.HasSuffix(strings.ToLower(part.FileName()), ".tar.gz") {
			tarArchiveName = part.FileName()
		}

		if strings.ToLower(part.FileName()) == "dockerfile" {
			hasDockerfile = true
		}

	}

	// Building from tar
	if tarArchiveName != "" && !hasDockerfile {
		err = s.db.SetBuildState(buildID, db.NewBuildStateOk("Importing an image from a tar archive", -1))
		if err != nil {
			log.Println(err)
		}
		go s.pier.StartServiceBuildFromTar(buildID, buildDir, connectionID, user.ID, tarArchiveName)
	}

	// Building from a Dockerfile
	if hasDockerfile {
		err := s.db.SetBuildState(buildID, db.NewBuildStateOk("Building an image from a Dockerfile", -1))
		if err != nil {
			log.Println(err)
		}
		go s.pier.StartServiceBuildFromFile(buildID, buildDir, connectionID, user.ID)
	}

	Response{w}.Ok(jmap("buildID", buildID))
}

func (s *Server) inspectBuildImageHandler(w http.ResponseWriter, r *http.Request, e environment) {
	allow, user := Authorization{s, w, r}.allowInspectBuild()
	if user == nil || !allow {
		return
	}

	vars := mux.Vars(r)
	buildID := vars["buildID"]

	build, err := s.db.GetBuild(buildID)
	if err != nil {
		Response{w}.ClientError("cannot get a build", err)
		return
	}
	Response{w}.Ok(jmap("Build", build))
}

func (s *Server) listServicesHandler(w http.ResponseWriter, r *http.Request, e environment) {
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

func (s *Server) inspectServiceHandler(w http.ResponseWriter, r *http.Request, e environment) {
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

func (s *Server) editServiceHandler(w http.ResponseWriter, r *http.Request, e environment) {
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

	oldService, err := s.db.GetService(service.ID)
	if err != nil {
		Response{w}.ClientError("cannot retrieve old version of the service", err)
		return
	}

	err = s.db.RemoveService(service.ID)
	if err != nil {
		Response{w}.ClientError("cannot remove service", err)
		return
	}

	service.ConnectionID = oldService.ConnectionID
	err = s.db.AddService(user.ID, service)
	if err != nil {
		Response{w}.ClientError("cannot add service", err)
		return
	}

	Response{w}.Ok(jmap("Service", service))
}

func (s *Server) removeServiceHandler(w http.ResponseWriter, r *http.Request, e environment) {
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

	err = s.db.RemoveService(serviceID)
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

func (s *Server) executeServiceHandler(w http.ResponseWriter, r *http.Request, e environment) {
	input := r.FormValue("pid")
	if input == "" {
		vars := mux.Vars(r)
		input = vars["pid"]
	}

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

	if serviceID == "" {
		Response{w}.ServerNewError("execute docker image: serviceID required")
		return
	}

	service, err := s.db.GetService(db.ServiceID(serviceID))
	if err != nil {
		Response{w}.ClientError("cannot get service", err)
		return
	}

	// getting multiple inputs
	var allInputs []string
	if input == "" {
		for _, value := range service.Input {
			inputName := "pid_" + value.ID
			currentInput := r.FormValue(inputName)
			if currentInput == "" {
				vars := mux.Vars(r)
				currentInput = vars[inputName]
			}

			// creating a temporary input file
			if (strings.ToLower(value.Type) == "string") && (value.FileName != "") {
				path, _, err := def.NewRandomTmpDir(s.tmpDir, inputTmpDir)
				if err != nil {
					Response{w}.ServerError("cannot create a temporary folder for an input file", err)
					return
				}
				tmpFile, err := os.Create(filepath.Join(path, value.FileName))
				if err != nil {
					Response{w}.ServerError("cannot create a temporary input file", err)
					return
				}
				defer tmpFile.Close()
				_, err = tmpFile.WriteString(currentInput)
				if err != nil {
					Response{w}.ServerError("cannot write string data into a file", err)
					return
				}
				currentInput = filepath.Join(path, value.FileName)
			}
			allInputs = append(allInputs, currentInput)
		}
	} else {
		allInputs = append(allInputs, input)
	}

	if len(allInputs) == 0 {
		Response{w}.ServerNewError("execute docker image: some input is required")
		return
	}

	job, err := s.pier.RunService(user.ID, service.ID, allInputs, s.limits, s.timeouts)
	if err != nil {
		Response{w}.ServerError("cannot read the requested file from the archive", err)
		return
	}

	loc, err := urljoin(r, string(job.ID))
	if err != nil {
		Response{w}.ServerError("urljoin error", err)
		return
	}
	Response{w}.Location(loc).Created(jmap("Location", loc, "jobID", job.ID))
}

func (s *Server) listJobsHandler(w http.ResponseWriter, r *http.Request, e environment) {
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

func (s *Server) inspectJobHandler(w http.ResponseWriter, r *http.Request, e environment) {
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

func (s *Server) removeJobHandler(w http.ResponseWriter, r *http.Request, e environment) {
	vars := mux.Vars(r)
	jobID := db.JobID(vars["jobID"])
	allow, _ := Authorization{s, w, r}.allowRemoveJob(jobID)
	if !allow {
		return
	}

	job, err := s.db.GetJob(jobID)
	if err != nil {
		Response{w}.ClientError(err.Error(), err)
		return
	}

	err = s.db.RemoveJob(jobID)
	if err != nil {
		Response{w}.ClientError(err.Error(), err)
		return
	}
	Response{w}.Ok(jmap("Job", job))
}

func (s *Server) volumeContentHandler(w http.ResponseWriter, r *http.Request, e environment) {
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
		err := s.pier.DownStreamContainerFile(vars["volumeID"], filepath.Join("/root/volume/", fileLocation), s.limits, s.timeouts, w)
		if err != nil {
			Response{w}.ServerError("downloading volume files failed", err)
			return
		}

		Response{w}.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		Response{w}.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	} else { // Return of list of files in a specific location in a volume
		volumeFiles, err := s.pier.ListFiles(db.VolumeID(vars["volumeID"]), fileLocation, s.limits, s.timeouts)
		if err != nil {
			Response{w}.ServerError("streaming container files failed", err)
			return
		}
		Response{w}.Ok(jmap("volumeID", vars["volumeID"], "volumeContent", volumeFiles))
	}
}
