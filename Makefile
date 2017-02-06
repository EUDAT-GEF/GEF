GOSRC = ./../..
EUDATSRC = ./..
WEBUI = frontend/src/main/webui
EPICPID = ../EpicPID

build: dependencies
	(cd $(WEBUI) && node_modules/webpack/bin/webpack.js -p)
	(cd $(EPICPID) && mvn package install)
	(cd frontend && mvn -q package)
	(cd services/_internal/volume-filelist && go build && docker build -t volume-filelist .)
	$(GOPATH)/bin/golint ./...
	go vet ./...
	go test ./...
	go build ./...

dependencies: $(WEBUI)/node_modules $(GOSRC)/golang/lint/golint $(GOSRC)/fsouza/go-dockerclient $(GOSRC)/gorilla/mux $(GOSRC)/pborman/uuid $(EUDATSRC)/EpicPID

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

$(EUDATSRC)/EpicPID:
	(cd $(EUDATSRC) && git clone https://github.com/EUDAT-GEF/EpicPID)
	(cd $(EUDATSRC)/EpicPID && mvn package install)

webui_dev_server:
	(cd $(WEBUI) && node_modules/webpack-dev-server/bin/webpack-dev-server.js --config webpack.config.devel.js)

run_frontend:
	@$(eval JAR = $(shell find frontend/target -iname 'GEF-*.jar'))
	java -cp frontend/src/main/resources:$(JAR) eu.eudat.gef.app.GEF server frontend/gefconfig.yml
	# @java -jar $(JAR) server gefconfig.yml

run_backend:
	(cd backend-docker && go run main.go)

clean:
	go clean ./...
	(cd frontend && mvn -q clean)

.PHONY: build dependencies webui_dev_server run_frontend run_backend clean
