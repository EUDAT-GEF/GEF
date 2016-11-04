package api_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/eudat-gef/gef/services/volume-inspector/api"
)

var (
	server   *httptest.Server
	reader   io.Reader
	lsUrl string
)

func init() {
	server = httptest.NewServer(api.Handlers())
	lsUrl = fmt.Sprintf("%s/ls", server.URL)
}

func TestCreateUser(t *testing.T) {
	userJson := `{"folderPath": "./"}`
	reader = strings.NewReader(userJson)
	request, err := http.NewRequest("POST", lsUrl, reader)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 201 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
}