package server

import (
	"testing"
	"github.com/EUDAT-GEF/GEF/backend-docker/dckr"
	"net/http/httptest"
	"fmt"
	"net/http"
)

func TestServer(t *testing.T) {
	c := newClient(t)
	s := createServer(c, t)
	srv := httptest.NewServer(s.server.Handler)
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

func createServer(client dckr.Client, t *testing.T) *Server {
	config := Config{
		Address: ":4142",
		ReadTimeoutSecs: 10,
		WriteTimeoutSecs: 10,
	}

	server := NewServer(config, client)

	return server
}

func newClient(t *testing.T) dckr.Client {
	config2 := []dckr.Config{
		dckr.Config{Endpoint: "unix:///var/run/docker.sock"},
		dckr.Config{UseBoot2Docker: true},
	}
	c, err := dckr.NewClientFirstOf(config2)
	if err != nil {
		t.Error(err)
		t.Error("--- client is not valid (this test requires a docker server)")
		t.Fail()
	}

	if !c.IsValid() {
		t.Error("client not valid (unable to ping)")
	}
	return c
}