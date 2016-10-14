package server

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/eudat-gef/gef-docker/dckr"
)

// GefSrvLabelPrefix is the prefix identifying GEF related labels
const GefSrvLabelPrefix = "eudat.gef.service."

// Service describes metadata for a GEF service
type Service struct {
	ID          dckr.ImageID
	Name        string
	RepoTag     string
	Description string
	Version     string
	Input       []IOPort
	Output      []IOPort
}

// IOPort is an i/o specification for a service
// The service can only read data from volumes and write to a single volume
// Path specifies where the volumes are mounted
type IOPort struct {
	VolumeID string
	Name     string
	Path     string
}

// Job is an instance of a running service
type Job struct {
	dckr.Container
	Service Service
}

func makeJob(container dckr.Container) Job {
	r := Job{Container: container, Service: extractServiceInfo(container.Image)}
	return r
}

///////////////////////////////////////////////////////////////////////////////

func extractServiceInfo(image dckr.Image) Service {
	srv := Service{
		ID:      image.ID,
		RepoTag: image.RepoTag,
	}

	for k, v := range image.Labels {
		if !strings.HasPrefix(k, GefSrvLabelPrefix) {
			continue
		}
		k = k[len(GefSrvLabelPrefix):]
		// fmt.Println(k, " -> ", v)
		ks := strings.Split(k, ".")
		if len(ks) == 0 {
			continue
		}
		switch ks[0] {
		case "name":
			srv.Name = v
		case "description":
			srv.Description = v
		case "version":
			srv.Version = v
		case "input":
			addVecValue(&srv.Input, ks[1:], v)
		case "output":
			addVecValue(&srv.Output, ks[1:], v)
		default:
			log.Println("Unknown GEF service label: ", k, "=", v)
		}
	}

	{
		in := make([]IOPort, 0, len(srv.Input))
		for _, p := range srv.Input {
			if p.Path != "" {
				p.VolumeID = fmt.Sprintf("input%d", len(in))
				in = append(in, p)
			}
		}
		srv.Input = in
	}
	{
		out := make([]IOPort, 0, len(srv.Output))
		for _, p := range srv.Output {
			if p.Path != "" {
				p.VolumeID = fmt.Sprintf("output%d", len(out))
				out = append(out, p)
			}
		}
		srv.Output = out
	}

	return srv
}

func addVecValue(vec *[]IOPort, ks []string, value string) {
	if len(ks) < 2 {
		log.Println("ERROR: GEF service label I/O key error (need 'port number . key name')", ks)
		return
	}
	id, err := strconv.ParseUint(ks[0], 10, 8)
	if err != nil {
		log.Println("ERROR: GEF service label: expecting integer argument for IOPort, instead got: ", ks)
	}
	for len(*vec) < int(id)+1 {
		*vec = append(*vec, IOPort{})
	}
	switch ks[1] {
	case "name":
		(*vec)[id].Name = value
	case "path":
		(*vec)[id].Path = value
	}
}
