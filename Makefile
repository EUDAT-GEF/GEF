ASSETDIR = src/main/resources/assets
LIBDIR = $(ASSETDIR)/lib
FONTDIR = $(ASSETDIR)/fonts
JSDIR = $(ASSETDIR)/js
WEBUI = frontend/src/main/webui
GOSRC = ./../..

build: $(WEBUI)/node_modules $(GOSRC)/fsouza/go-dockerclient $(GOSRC)/gorilla/mux $(GOSRC)/pborman/uuid
	(cd $(WEBUI) && node_modules/webpack/bin/webpack.js -p)
	(cd frontend && mvn -q package)
	(cd backend-docker && golint ./...)
	(cd backend-docker && go vet ./...)
	(cd backend-docker && go test ./...)
	(cd backend-docker && go build)

$(WEBUI)/node_modules:
	(cd $(WEBUI) && npm install)

$(GOSRC)/fsouza/go-dockerclient:
	go get github.com/fsouza/go-dockerclient

$(GOSRC)/gorilla/mux:
	go get github.com/gorilla/mux

$(GOSRC)/pborman/uuid:
	go get github.com/pborman/uuid

webui_dev_server:
	(cd $(WEBUI) && node_modules/webpack-dev-server/bin/webpack-dev-server.js --config webpack.config.devel.js)

run_frontend:
	@$(eval JAR = $(shell find frontend/target -iname 'GEF-*.jar'))
	java -cp frontend/src/main/resources:$(JAR) eu.eudat.gef.app.GEF server frontend/gefconfig.yml
	# @java -jar $(JAR) server gefconfig.yml

run_backend:
	(cd backend-docker && go run main.go)

clean:
	(cd backend-docker && go clean)
	(cd frontend && mvn -q clean)

.PHONY: build webui_dev_server run_frontend run_backend clean
