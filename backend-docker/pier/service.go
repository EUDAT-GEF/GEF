package pier

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pborman/uuid"

	"github.com/EUDAT-GEF/GEF/backend-docker/pier/internal/dckr"
)

// GefSrvLabelPrefix is the prefix identifying GEF related labels
const GefSrvLabelPrefix = "eudat.gef.service."

// Service describes metadata for a GEF service
type Service struct {
	ID          ServiceID
	imageID     dckr.ImageID
	Name        string
	RepoTag     string
	Description string
	Version     string
	Created     time.Time
	Size        int64
	Input       []IOPort
	Output      []IOPort
}

// ServiceID exported
type ServiceID string

// IOPort is an i/o specification for a service
// The service can only read data from volumes and write to a single volume
// Path specifies where the volumes are mounted
type IOPort struct {
	ID   string
	Name string
	Path string
}

type srvArray []Service

func (sl srvArray) Len() int {
	return len(sl)
}
func (sl srvArray) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}
func (sl srvArray) Less(i, j int) bool {
	return sl[i].Created.After(sl[j].Created)
}

// ServiceList is a shared structure that stores info about all services
type ServiceList struct {
	sync.Mutex
	cache map[ServiceID]Service
}

// NewServiceList exportedk
func NewServiceList() *ServiceList {
	return &ServiceList{
		cache: make(map[ServiceID]Service),
	}
}

func (serviceList *ServiceList) add(service Service) {
	serviceList.Lock()
	defer serviceList.Unlock()
	serviceList.cache[service.ID] = service
}

func (serviceList *ServiceList) list() []Service {
	serviceList.Lock()
	defer serviceList.Unlock()
	all := make([]Service, len(serviceList.cache), len(serviceList.cache))
	i := 0
	for _, service := range serviceList.cache {
		all[i] = service
		i++
	}
	sort.Sort(srvArray(all))

	return all
}

func (serviceList *ServiceList) get(key ServiceID) (Service, bool) {
	serviceList.Lock()
	defer serviceList.Unlock()
	service, ok := serviceList.cache[key]
	return service, ok

}

///////////////////////////////////////////////////////////////////////////////

func newServiceFromImage(image dckr.Image) Service {
	srv := Service{
		ID:      ServiceID(uuid.New()),
		imageID: image.ID,
		RepoTag: image.RepoTag,
		Created: image.Created,
		Size:    image.Size,
	}

	for k, v := range image.Labels {
		if !strings.HasPrefix(k, GefSrvLabelPrefix) {
			continue
		}
		k = k[len(GefSrvLabelPrefix):]
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
				p.ID = fmt.Sprintf("input%d", len(in))
				in = append(in, p)
			}
		}
		srv.Input = in
	}
	{
		out := make([]IOPort, 0, len(srv.Output))
		for _, p := range srv.Output {
			if p.Path != "" {
				p.ID = fmt.Sprintf("output%d", len(out))
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