GOSRC = ./../../..
GITHUBSRC = ./../..
EUDATSRC = ./..
WEBUI = webui
JSBUNDLE = webui/app/gef-bundle.js

build: dependencies webui backend

webui: $(JSBUNDLE)

$(JSBUNDLE):
	(cd $(WEBUI) && node_modules/webpack/bin/webpack.js -p)

backend:
	$(GOPATH)/bin/golint ./...
	go vet ./...
	go build ./...
	GEF_SECRET_KEY="test" go test -timeout 4m ./...

clean:
	go clean ./...
	rm $(JSBUNDLE) $(JSBUNDLE).map
	rm -r $(WEBUI)/node_modules

run_webui_dev_server:
	(cd $(WEBUI) && node_modules/webpack-dev-server/bin/webpack-dev-server.js -d --hot --https --config webpack.config.devel.js)

run_gef:
	(cd gefserver && go run main.go)


certificate:
	@echo "Creating self-signed GEF web server certificate in ./ssl/"
	@mkdir -p ssl
	@openssl req -x509 -nodes -newkey rsa:2048 -keyout ssl/server.key -out ssl/server.crt -days 365 -subj "/C=EU/ST=Helsinki/L=Helsinki/O=EUDAT/OU=GEF/CN=gef" 2>&1


dependencies: $(WEBUI)/node_modules \
	          $(GITHUBSRC)/golang/lint/golint \
	          $(GITHUBSRC)/fsouza/go-dockerclient \
	          $(GITHUBSRC)/gorilla/mux \
	          $(GITHUBSRC)/gorilla/sessions \
	          $(GITHUBSRC)/pborman/uuid \
	          $(GITHUBSRC)/mattn/go-sqlite3 \
	          $(GOSRC)/golang.org/x/oauth2 \
	          $(GOSRC)/gopkg.in/gorp.v1

$(WEBUI)/node_modules:
	(cd $(WEBUI) && npm install)

$(GITHUBSRC)/golang/lint/golint:
	go get -u github.com/golang/lint/golint

$(GITHUBSRC)/fsouza/go-dockerclient:
	go get -u github.com/fsouza/go-dockerclient

$(GITHUBSRC)/gorilla/mux:
	go get -u github.com/gorilla/mux

$(GITHUBSRC)/gorilla/sessions:
	go get -u github.com/gorilla/sessions

$(GITHUBSRC)/pborman/uuid:
	go get -u github.com/pborman/uuid

$(GITHUBSRC)/mattn/go-sqlite3:
	go get -u github.com/mattn/go-sqlite3

$(GOSRC)/golang.org/x/oauth2:
	go get -u golang.org/x/oauth2

$(GOSRC)/gopkg.in/gorp.v1:
	go get -u gopkg.in/gorp.v1

.PHONY: build webui backend run_gef run_webui_dev_server certificates dependencies
