package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier"
	"github.com/EUDAT-GEF/GEF/backend-docker/server"
)

func TestServer(t *testing.T) {
	config, err := def.ReadConfigFile(configFilePath)
	checkMsg(t, err, "reading config files")

	var p *pier.Pier
	p, err = pier.NewPier(config.Docker, config.TmpDir)
	checkMsg(t, err, "creating new pier")

	var srv *httptest.Server
	{
		s, err := server.NewServer(config.Server, p, config.TmpDir)
		checkMsg(t, err, "creating api server")
		srv = httptest.NewServer(s.Server.Handler)
	}
	defer srv.Close()
	baseURL := srv.URL + "/api/"

	checkRunRequest(t, "GET", baseURL, 200)

	json := checkGetJSON(t, baseURL+"services")
	services, ok := json["Services"]
	expect(t, ok, "Services not found in returned json")
	expect(t, services != nil, "nil Services in returned json")

	json = checkGetJSON(t, baseURL+"jobs")
	jobs, ok := json["Jobs"]
	expect(t, ok, "Jobs not found in returned json")
	expect(t, jobs != nil, "nil Jobs in returned json")
}

func checkGetJSON(t *testing.T, url string) map[string]interface{} {
	res := checkRunRequest(t, "GET", url, 200)
	defer res.Body.Close()

	htmlData, err := ioutil.ReadAll(res.Body)
	check(t, err)

	var j map[string]interface{}
	err = json.Unmarshal([]byte(htmlData), &j)
	check(t, err)
	return j
}

func checkRunRequest(t *testing.T, method string, url string, expectedCode int) *http.Response {
	request, err := http.NewRequest(method, url, nil)
	check(t, err)
	res, err := http.DefaultClient.Do(request)
	check(t, err)
	expect(t, res.StatusCode == expectedCode,
		fmt.Sprintf("unexpected http request status code: %d instead of %d",
			res.StatusCode, expectedCode))
	return res
}

func isInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
