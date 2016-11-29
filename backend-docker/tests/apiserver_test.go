package tests

import (
	"testing"
	"github.com/EUDAT-GEF/GEF/backend-docker/config"
	"github.com/EUDAT-GEF/GEF/backend-docker/server"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func isInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func checkJSONReply(callURL string, keyList []string, t *testing.T) {
	request, err := http.NewRequest("GET", callURL, nil)
	if err != nil {
		t.Error(err)
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Error("Error code: ", res.StatusCode)
		t.Fail()
	} else {
		htmlData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error("Cannot read response body: ", err)
			t.Fail()
		}
		correctReply := isJSON(string(htmlData))

		if correctReply != true {
			t.Error("Reply is not JSON")
			t.Fail()
		}

		c := make(map[string]interface{})
		err = json.Unmarshal(htmlData, &c)
		if err != nil {
			t.Error("Error while reading JSON: ", err)
			t.Fail()
		}

		for s, _ := range c {
			if isInSlice(s, keyList) !=true {
				t.Error("The following key was not found in JSON: ", s)
				t.Error("Reply is incorrect")
				t.Fail()
			}
		}
	}

}
func TestServer(t *testing.T) {
	settings, err := config.ReadConfigFile(configFilePath)
	if err != nil {
		t.Error("FATAL while reading config files: ", err)
	}
	clientConf = settings.Docker

	c := newClient(t)
	s := server.NewServer(settings.Server, c)

	srv := httptest.NewServer(s.Server.Handler)
	baseURL := srv.URL + "/api/"

	checkIfAPIExist(baseURL, t)
	callListVolumesHandler(baseURL + "volumes", t)
	callListJobsHandler(baseURL + "jobs", t)
	callListServicesHandler(baseURL + "images", t)
}

func checkIfAPIExist(callURL string, t *testing.T) bool {
	request, err := http.NewRequest("GET", callURL, nil)
	if err != nil {
		t.Error(err)
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Error("Error code: ", res.StatusCode)
		t.Fail()
		return false
	} else {
		return true
	}
}

func callListVolumesHandler(callURL string, t *testing.T) {
	volumeKeys := []string{"Volumes"}
	checkJSONReply(callURL, volumeKeys, t)
}

func callListJobsHandler(callURL string, t *testing.T) {
	jobKeys := []string{"Jobs"}
	checkJSONReply(callURL, jobKeys, t)
}

func callListServicesHandler(callURL string, t *testing.T) {
	jobKeys := []string{"Images", "Services"}
	checkJSONReply(callURL, jobKeys, t)
}
