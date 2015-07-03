package main

import (
	"log"
	"strconv"
	"strings"
)

// GefSrvLabelPrefix is the prefix identifying GEF related labels
const GefSrvLabelPrefix = "eudat.gef.service."

// Service describes metadata for a GEF service
type Service struct {
	Name        string
	Description string
	Version     string
	Input       []IOPort
	Output      []IOPort
}

// IOPort is an i/o specification for a service
type IOPort struct {
	Path string
}

///////////////////////////////////////////////////////////////////////////////

func extractServiceInfo(labels map[string]string) Service {
	srv := Service{}

	for k, v := range labels {
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
				in = append(in, p)
			}
		}
		srv.Input = in
	}
	{
		out := make([]IOPort, 0, len(srv.Output))
		for _, p := range srv.Output {
			if p.Path != "" {
				out = append(out, p)
			}
		}
		srv.Output = out
	}
	return srv
}

func addVecValue(vec *[]IOPort, ks []string, value string) {
	if len(ks) == 0 {
		log.Println("ERROR: GEF service label I/O key empty")
		return
	}
	if len(ks) > 1 {
		log.Println("ERROR: GEF service label I/O key has too many parts: ", ks)
		return
	}
	id, err := strconv.ParseUint(ks[0], 10, 8)
	if err != nil {
		log.Println("ERROR: GEF service label: expecting integer argument for IOPort, instead got: ", ks)
	}
	for len(*vec) < int(id)+1 {
		*vec = append(*vec, IOPort{})
	}
	(*vec)[id].Path = value
}
