package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/EUDAT-GEF/GEF/gefserver/db"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
)

func (s *Server) decorate(fn func(http.ResponseWriter, *http.Request), actionType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		user, err := s.getCurrentUser(r)
		if err != nil {
			log.Println("User error: ", err)
		}
		var userEnv environment
		userEnv.Timeouts = s.timeouts
		userEnv.Limits = s.limits

		var sysStatistics statistics

		sysStatistics.UserRunningJobs = 0
		if user != nil {
			sysStatistics.UserRunningJobs = s.db.CountUserRunningJobs(user.ID)
		}
		sysStatistics.TotalRunningJobs = s.db.CountRunningJobs()

		allow, closefn := signalEvent(actionType, user, userEnv, sysStatistics, r)
		if !allow {
			Response{w}.DirectiveError()
		} else {
			defer closefn()
			fn(w, r)
		}
	}
}

type eventSystem struct {
	address string
	ID      uint64
}

var eventSys eventSystem

type eventPayload struct {
	Event       event       `json:"event"`
	User        *db.User    `json:"user"`
	Resource    resource    `json:"resource"`
	Environment environment `json:"environment"`
	Statistics  statistics  `json:"statistics"`
}

type event struct {
	ID        uint64 `json:"id"`
	Time      string `json:"time"`   // iso string
	Action    string `json:"action"` // one of: "DataAnalysis"...
	Preceding bool   `json:"preceding"`
}

type resource struct {
	URL    string `json:"uri"`    // "/api/jobs",
	Method string `json:"method"` // "POST",
}

type environment struct {
	Timeouts def.TimeoutConfig `json:"timeouts"`
	Limits   def.LimitConfig   `json:"limits"`
}

type statistics struct {
	UserRunningJobs  int64 `json:"userRunningJobs"`
	TotalRunningJobs int64 `json:"totalRunningJobs"`
}

// InitEventSystem initializes the event system, or disables it if the address
// is empty
func InitEventSystem(address string) {
	if address == "" {
		log.Println("Event system disabled")
		return
	}
	if !strings.HasSuffix(address, "/") {
		address += "/"
	}
	eventSys = eventSystem{
		address: address + "api/events",
	}
}

func (es eventSystem) dispatch(action string, user *db.User, userEnv environment, sysStatistics statistics, r *http.Request, preceding bool) (bool, *http.Response) {
	if es.address == "" {
		return true, nil
	}

	payload := eventPayload{
		Event: event{
			ID:        atomic.AddUint64(&es.ID, 1),
			Time:      time.Now().Format(time.RFC3339),
			Action:    action,
			Preceding: preceding,
		},
		User: user,
		Resource: resource{
			URL:    r.RequestURI,
			Method: r.Method,
		},
		Environment: userEnv,
		Statistics:  sysStatistics,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		log.Printf("event system: json marshal ERROR: %#v\n", err)
	} else {
		resp, err := http.Post(es.address, "application/json", bytes.NewBuffer(json))
		if err != nil {
			if e, ok := err.(*url.Error); ok {
				if e, ok := e.Err.(*net.OpError); ok {
					log.Printf("event system: post url net ERROR: %#v\n", e)
				}
				log.Printf("event system: post url ERROR: %#v\n", e)
			}
			log.Printf("event system: post ERROR: %#v\n", err)
			return true,  resp
		}
		if 400 <= resp.StatusCode && resp.StatusCode < 500 {
			log.Printf("event system: post client ERROR: %#v\n", resp)
			return false, resp
		}
		return true, resp
	}
	return true, nil
}

func signalEvent(action string, user *db.User, userEnv environment, sysStatistics statistics, r *http.Request) (allow bool, closefn func()) {
	fn := func() {}
	ret, resp := eventSys.dispatch(action, user, userEnv, sysStatistics, r, true)
	if resp != nil {
		defer resp.Body.Close()

		newEnv := environment{}
		err := json.NewDecoder(resp.Body).Decode(&newEnv)
		if err != nil {
			log.Printf("event system response parsing error: %#v\n", err)
		}
		if ret {
			fn = func() {
				eventSys.dispatch(action, user, userEnv, sysStatistics, r, false)
			}
		}
	}
	return ret, fn
}
