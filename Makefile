GOSRC = ./../../..
GITHUBSRC = ./../..
EUDATSRC = ./..
WEBUI = frontend/webui
INTERNALSERVICES = services/_internal
GOFLAGS=

build: dependencies webui backend

webui: $(WEBUI)/
	(cd $(WEBUI) && node_modules/webpack/bin/webpack.js -p)

backend:
	$(GOPATH)/bin/golint ./...
	go vet ./...
	go test -timeout 30s $(GOFLAGS) ./...
	go build $(GOFLAGS) ./...

run_webui_dev_server:
	(cd $(WEBUI) && node_modules/webpack-dev-server/bin/webpack-dev-server.js --config webpack.config.devel.js)

run_gef:
	(cd backend-docker && go run $(GOFLAGS) main.go)

certificate:
	@echo "Creating self-signed GEF web server certificate in ./ssl/"
	@mkdir -p ssl
	@openssl req -x509 -nodes -newkey rsa:2048 -keyout ssl/server.key -out ssl/server.crt -days 365 -subj "/C=EU/ST=Helsinki/L=Helsinki/O=EUDAT/OU=GEF/CN=gef" 2>&1

dependencies: $(WEBUI)/node_modules \
	          $(GITHUBSRC)/golang/lint/golint \
	          $(GITHUBSRC)/fsouza/go-dockerclient \
	          $(GITHUBSRC)/gorilla/mux \
	          $(GITHUBSRC)/pborman/uuid \
	          $(GITHUBSRC)/mattn/go-sqlite3 \
	          $(GOSRC)/gopkg.in/gorp.v1

$(WEBUI)/node_modules:
	(cd $(WEBUI) && npm install)

$(GITHUBSRC)/golang/lint/golint:
	go get -u github.com/golang/lint/golint

$(GITHUBSRC)/fsouza/go-dockerclient:
	go get -u github.com/fsouza/go-dockerclient

$(GITHUBSRC)/gorilla/mux:
	go get -u github.com/gorilla/mux

$(GITHUBSRC)/pborman/uuid:
	go get -u github.com/pborman/uuid

$(GITHUBSRC)/mattn/go-sqlite3:
	go get -u github.com/mattn/go-sqlite3

$(GOSRC)/gopkg.in/gorp.v1:
	go get -u gopkg.in/gorp.v1

.PHONY: build webui backend run_gef run_webui_dev_server certificates dependencies
