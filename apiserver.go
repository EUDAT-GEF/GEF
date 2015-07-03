package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"../eudock"

	"github.com/gorilla/mux"
)

// ErrPreMarshal reports a coding error before marshalling json
var ErrPreMarshal = errors.New("pre-marshalling json error")

// ServerConfig keeps the configuration options needed to make a Server
type ServerConfig struct {
	Address          string
	ReadTimeoutSecs  int
	WriteTimeoutSecs int
}

// Server is a master struct for serving HTTP API requests
type Server struct {
	server http.Server
	docker *eudock.Client
}

// NewServer creates a new Server
func NewServer(cfg ServerConfig, docker eudock.Client) *Server {
	server := &Server{
		http.Server{
			Addr: cfg.Address,
			// timeouts seem to trigger even after a correct read
			// ReadTimeout: 	cfg.ReadTimeoutSecs * time.Second,
			// WriteTimeout: 	cfg.WriteTimeoutSecs * time.Second,

		}, &docker,
	}
	router := mux.NewRouter()
	apirouter := router.PathPrefix("/api").Subrouter()
	apirouter.HandleFunc("/", server.listServicesHandler).Methods("GET")
	apirouter.HandleFunc("/", server.buildServiceImageHandler).Methods("POST")
	apirouter.HandleFunc("/{imageID}", server.inspectServiceHandler).Methods("GET")

	server.server.Handler = router
	return server
}

// Start starts a new http listener
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s Server) listServicesHandler(w http.ResponseWriter, r *http.Request) {
	imgIDs, err := s.docker.ListImages()
	if err != nil {
		log.Println("ERROR: list of docker images: ", err)
		s.serverError(w, err)
		return
	}
	s.okJMap(w, "ImageIDs", imgIDs)
}

func (s Server) inspectServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := eudock.ImageID(vars["imageID"])
	image, err := s.docker.InspectImage(imageID)
	if err != nil {
		log.Println("ERROR: inspect docker image: ", err)
		s.serverError(w, err)
		return
	}
	srv := extractServiceInfo(image.Labels)
	s.okJMap(w, "Image", image, "Service", srv)
}

func (s Server) buildServiceImageHandler(w http.ResponseWriter, r *http.Request) {
	image, err := s.docker.BuildImage("../gef-service-example")
	if err != nil {
		log.Println("ERROR: build docker image: ", err)
		s.serverError(w, err)
		return
	}
	s.okJMap(w, "Image", image)
}

func (s Server) executeServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := eudock.ImageID(vars["imageID"])
	container, err := s.docker.ExecuteImage(imageID)
	if err != nil {
		log.Println("ERROR: execute docker image: ", err)
		s.serverError(w, err)
		return
	}
	s.okJMap(w, "ContainerID", container.ID)
}

func (s Server) okJMap(w http.ResponseWriter, kv ...interface{}) {
	if len(kv) == 0 {
		log.Println("ERROR: jsonmap: empty call")
		s.serverError(w, ErrPreMarshal)
		return
	} else if len(kv)%2 == 1 {
		log.Println("ERROR: jsonmap: unbalanced call")
		s.serverError(w, ErrPreMarshal)
		return
	}
	m := make(map[string]interface{})
	k := ""
	for _, kv := range kv {
		if k == "" {
			if skv, ok := kv.(string); !ok {
				log.Println("ERROR: jsonmap: expected string key")
				s.serverError(w, ErrPreMarshal)
				return
			} else if skv == "" {
				log.Println("ERROR: jsonmap: string key is empty")
				s.serverError(w, ErrPreMarshal)
				return
			} else {
				k = skv
			}
		} else {
			m[k] = kv
			k = ""
		}
	}
	s.okJSON(w, m)
}

func (s Server) okJSON(w http.ResponseWriter, data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		log.Println("ERROR: json marshal: ", err)
		s.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(json)
}

func (s Server) serverError(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintln("Server Error: ", err), 500)
}

// func notFoundHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	w.Write([]byte("Resource not found"))
// }
