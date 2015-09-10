package server

import (
	"encoding/json"
	"fmt"
	"log"

	"net/http"
)

// Response encapsulates a http ResponseWriter
type Response struct {
	http.ResponseWriter
}

// Location sets location header
func (w Response) Location(loc string) Response {
	w.Header().Set("Location", loc)
	return w
}

// ServerError returns a server error
func (w Response) ServerError(message string, err error) {
	str := fmt.Sprintf("API Server ERROR: %s %s", message, err.Error())
	http.Error(w, str, 500)
}

// Ok sets the body for a ok response
func (w Response) Ok(body interface{}) {
	if jsonMap, ok := body.(jsonMap); ok {
		json, err := json.Marshal(jsonMap.m)
		if err != nil {
			w.ServerError("json marshal: ", err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(json)
		// log.Println("ok json:", string(json))
	} else if str, ok := body.(string); ok {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(str))
		// log.Println("ok string:", str)
	} else {
		log.Printf("ERROR: unexpected Ok body type: %T\n", body)
		http.Error(w, fmt.Sprintln("Server Error: unexpected Ok body type"), 500)
	}
}

type jsonMap struct {
	m map[string]interface{}
}

func jmap(kv ...interface{}) jsonMap {
	if len(kv) == 0 {
		log.Println("ERROR: jsonmap: empty call")
		return jsonMap{}
	} else if len(kv)%2 == 1 {
		log.Println("ERROR: jsonmap: unbalanced call")
		return jsonMap{}
	}
	m := make(map[string]interface{})
	k := ""
	for _, kv := range kv {
		if k == "" {
			if skv, ok := kv.(string); !ok {
				log.Println("ERROR: jsonmap: expected string key")
				return jsonMap{m}
			} else if skv == "" {
				log.Println("ERROR: jsonmap: string key is empty")
				return jsonMap{m}
			} else {
				k = skv
			}
		} else {
			m[k] = kv
			k = ""
		}
	}
	return jsonMap{m}
}
