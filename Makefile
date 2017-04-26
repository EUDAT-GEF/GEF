GOSRC = ./../..
EUDATSRC = ./..
WEBUI = frontend/webui
INTERNALSERVICES = services/_internal
GOFLAGS=-ldflags -s

build: dependencies webui frontend containers backend

webui:
	(cd $(WEBUI) && node_modules/webpack/bin/webpack.js -p)

backend:
	$(GOPATH)/bin/golint ./...
	go vet ./...
	go test $(GOFLAGS) ./...
	go build $(GOFLAGS) ./...

dependencies: $(WEBUI)/node_modules $(GOSRC)/golang/lint/golint $(GOSRC)/fsouza/go-dockerclient $(GOSRC)/gorilla/mux $(GOSRC)/pborman/uuid $(GOSRC)/gopkg.in/gorp.v1 $(GOSRC)github.com/mattn/go-sqlite3

$(WEBUI)/node_modules:
	(cd $(WEBUI) && npm install)

$(GOSRC)/golang/lint/golint:
	go get -u github.com/golang/lint/golint

$(GOSRC)/fsouza/go-dockerclient:
	go get -u github.com/fsouza/go-dockerclient

$(GOSRC)/gorilla/mux:
	go get -u github.com/gorilla/mux

$(GOSRC)/pborman/uuid:
	go get -u github.com/pborman/uuid

$(GOSRC)/gopkg.in/gorp.v1:
	go get -u gopkg.in/gorp.v1

$(GOSRC)github.com/mattn/go-sqlite3:
	go get -u github.com/mattn/go-sqlite3

webui_dev_server:
	(cd $(WEBUI) && node_modules/webpack-dev-server/bin/webpack-dev-server.js --config webpack.config.devel.js)

run_gef:
	(cd backend-docker && go run $(GOFLAGS) main.go)

.PHONY: build dependencies webui frontend backend webui_dev_server run_frontend run_backend clean
