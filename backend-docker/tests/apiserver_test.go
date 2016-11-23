package tests

import (
	"testing"
	"github.com/EUDAT-GEF/GEF/backend-docker/dckr"
	"github.com/EUDAT-GEF/GEF/backend-docker/config"
	"github.com/EUDAT-GEF/GEF/backend-docker/server"
	"net/http/httptest"
	"fmt"
	"net/http"
)

func TestServer(t *testing.T) {
	settings, err := config.ReadConfigFile(configFilePath)
	if err != nil {
		t.Error("FATAL while reading config files: ", err)
	}
	clientConf = settings.Docker

	c := newClient(t)
	s := createServer(c, t, settings.Server)
	srv := httptest.NewServer(s.Server.Handler)
	lsUrl := fmt.Sprintf("%s/api/", srv.URL)


	//reader := strings.NewReader(userJson)
	request, err := http.NewRequest("GET", lsUrl, nil)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

}

func createServer(client dckr.Client, t *testing.T, servConf server.Config) *server.Server {
	/*config := server.Config{
		Address: ":4142",
		ReadTimeoutSecs: 10,
		WriteTimeoutSecs: 10,
	}*/

	server := server.NewServer(servConf, client)

	return server
}