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
)

type myResponseWriter struct {
	w    http.ResponseWriter
	code int
}

func (w myResponseWriter) Header() http.Header {
	return w.w.Header()
}
func (w myResponseWriter) Write(data []byte) (int, error) {
	return w.w.Write(data)
}
func (w myResponseWriter) WriteHeader(code int) {
	w.code = code
	w.w.WriteHeader(code)
}

type handlerfn func(http.ResponseWriter, *http.Request)

func decorate(actionType string, fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		allow, closefn := signalEvent(actionType, r)
		if !allow {
			Response{w}.DirectiveError()
		} else {
			writer := myResponseWriter{w, 0}
			defer closefn(writer)
			fn(writer, r)
		}
	}
}

type eventSystem struct {
	address string
	ID      uint64
}

var eventSys eventSystem

type eventPayload struct {
	Event       event                  `json:"event"`
	User        user                   `json:"user"`
	Resource    resource               `json:"resource"`
	Response    response               `json:"response,omitempty"`
	Environment map[string]interface{} `json:"environment"`
}

type event struct {
	ID        uint64 `json:"id"`
	Time      string `json:"time"`   // iso string
	Action    string `json:"action"` // one of: "DataAnalysis"...
	Preceding bool   `json:"preceding"`
}

type user struct {
	Email string `json:"email"` // "email@example.com",
	Name  string `json:"name"`  // "John Smith",
}

type resource struct {
	URL    string `json:"uri"`    // "/api/jobs",
	Method string `json:"method"` // "POST",
}

type response struct {
	Code int `json:"code"`
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

func (es eventSystem) dispatch(action string, r *http.Request, preceding bool, statusCode int) bool {
	if es.address == "" {
		return true
	}

	payload := eventPayload{
		Event: event{
			ID:        atomic.AddUint64(&es.ID, 1),
			Time:      time.Now().Format(time.RFC3339),
			Action:    action,
			Preceding: preceding,
		},
		User: user{
			Email: "",
			Name:  "",
		},
		Resource: resource{
			URL:    r.RequestURI,
			Method: r.Method,
		},
		Environment: map[string]interface{}{},
	}
	if !preceding {
		payload.Response = response{
			Code: statusCode,
		}
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
			return true
		}
		if 400 <= resp.StatusCode && resp.StatusCode < 500 {
			log.Printf("event system: post client ERROR: %#v\n", resp)
			return false
		}
	}
	return true
}

func signalEvent(action string, r *http.Request) (allow bool, closefn func(myResponseWriter)) {
	fn := func(w myResponseWriter) {}
	ret := eventSys.dispatch(action, r, true, 0)
	if ret {
		fn = func(w myResponseWriter) {
			eventSys.dispatch(action, r, false, w.code)
		}
	}
	return ret, fn
}
