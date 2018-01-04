package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

func (s *Server) decorate(apifn func(http.ResponseWriter, *http.Request, environment), actionType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)

		var userEnv environment
		userEnv.Timeouts = s.timeouts
		userEnv.Limits = s.limits

		if eventSys.address == "" {
			apifn(w, r, userEnv)
			return
		}

		dbuser, err := s.getCurrentUser(r)
		if err != nil {
			log.Println("User error: ", err)
		}

		var user *userdat
		if dbuser != nil {
			user = &userdat{
				ID:      dbuser.ID,
				Name:    dbuser.Name,
				Email:   dbuser.Email,
				Created: dbuser.Created,
				Updated: dbuser.Updated,
			}

			var dbroles []db.Role
			dbroles, err = s.db.GetUserRoles(dbuser.ID)
			if err != nil {
				log.Printf("event system: ERROR retrieving user roles: %#v\n", err)
			} else {
				for _, r := range dbroles {
					user.Roles = append(user.Roles, roledat{
						ID:            r.ID,
						Name:          r.Name,
						Description:   r.Description,
						CommunityID:   r.CommunityID,
						CommunityName: r.CommunityName,
					})
				}
			}
		}

		var sysStatistics statistics

		sysStatistics.UserRunningJobs = 0
		if user != nil {
			sysStatistics.UserRunningJobs = s.db.CountUserRunningJobs(user.ID)
		}
		sysStatistics.TotalRunningJobs = s.db.CountRunningJobs()

		allow, closefn, retEnv := sendEvent(actionType, user, userEnv, sysStatistics, r)
		if !allow {
			Response{w}.DirectiveError()
		} else {
			defer closefn()
			update(userEnv, retEnv)
			apifn(w, r, userEnv)
		}
	}
}

func update(userEnv environment, newEnv map[string]interface{}) {
	// TODO: !!!
	// TODO: update userEnv based on the newEnv values
	// TODO: !!!
}

type eventSystem struct {
	address string
	ID      uint64
}

var eventSys eventSystem

type eventPayload struct {
	Event       event       `json:"event"`
	User        *userdat    `json:"user"`
	Resource    resource    `json:"resource"`
	Environment environment `json:"environment"`
	Statistics  statistics  `json:"statistics"`
}

type userdat struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Roles   []roledat `json:"role"`
}

type roledat struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	CommunityID   int64  `json:"communityid"`
	CommunityName string `json:"communityname"`
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
	log.Println("Event System enabled: %v", address)
	eventSys = eventSystem{
		address: address + "api/events",
	}
}

func sendEvent(action string, user *userdat, userEnv environment, sysStatistics statistics, r *http.Request) (allow bool, closefn func(), newEnv map[string]interface{}) {
	closefn = func() {}
	jsonBody := eventSys.callDirectiveEngine(action, user, userEnv, sysStatistics, r, true)
	if jsonBody != nil && len(jsonBody) > 0 {
		err := json.NewDecoder(bytes.NewReader(jsonBody)).Decode(&newEnv)
		if err != nil {
			log.Printf("event system response parsing error: %#v\n", err)
		}
		log.Printf("    event system: directive response: %v\n", newEnv)

		allow = true
		if jsonAllow, exists := newEnv["allow"]; exists {
			if allowPrimitive, ok := jsonAllow.(bool); ok {
				allow = allowPrimitive
			} else if allowArray, ok := jsonAllow.([]interface{}); ok {
				for _, allowValue := range allowArray {
					a, ok := allowValue.(bool)
					if ok {
						allow = allow && a
					} else {
						log.Printf("event system: response parsing error: allow array containing non-boolean: %#v\n", allowArray)
					}
				}
			} else {
				log.Printf("event system: response parsing error: unknown allow type: %#v\n", jsonAllow)
			}
		}

		if allow {
			closefn = func() {
				eventSys.callDirectiveEngine(action, user, userEnv, sysStatistics, r, false)
			}
		}
	}
	return allow, closefn, newEnv
}

func (es eventSystem) callDirectiveEngine(action string, user *userdat, userEnv environment, sysStatistics statistics, r *http.Request, preceding bool) []byte {
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
			return nil
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("event system: body read ERROR: %#v\n", err)
			return nil
		}
		if 400 <= resp.StatusCode && resp.StatusCode < 500 {
			log.Printf("event system: post client ERROR: %#v\n", resp)
			return nil
		}
		return body
	}
	return nil
}
