GOSRC = ./../../..
GITHUBSRC = ./../..
EUDATSRC = ./..
WEBUI = webui
JSBUNDLE = webui/app/gef-bundle.js

build: dependencies backend webui

backend:
	$(GOPATH)/bin/golint ./gefserver
	go vet ./gefserver
	go build -i -o build/gefserver ./gefserver
	(cd ./gefserver/tests &&  GEF_SECRET_KEY="test" go test -coverpkg "../db","../def/","../pier/...","../server" -timeout 4m)

webui: $(JSBUNDLE)

run_gef:
	(cd gefserver && go run main.go)

run_webui_dev_server:
	(cd $(WEBUI) && node_modules/webpack-dev-server/bin/webpack-dev-server.js -d --hot --https --config webpack.config.devel.js)

certificate:
	@echo "-- Creating self-signed GEF web server certificate in ./ssl/"
	@mkdir -p ssl
	@openssl req -x509 -nodes -newkey rsa:2048 -keyout ssl/server.key -out ssl/server.crt -days 365 -subj "/C=EU/ST=Helsinki/L=Helsinki/O=EUDAT/OU=GEF/CN=gef" 2>&1

clean:
	go clean ./...
	rm -rf gefserver/vendor
	rm -rf $(WEBUI)/node_modules
	rm -f $(JSBUNDLE) $(JSBUNDLE).map

pack: dependencies webui certificate
	mkdir -p build
	mkdir -p build/bin
	docker build -t gefcompile:linux .
	docker run --rm -v $(PWD)/build:/go/src/github.com/EUDAT-GEF/GEF/build gefcompile:linux
	mv build/gef_linux ./build/bin
	cp gefserver/config.json ./build/bin
	cp -r services ./build/
	cp -r ssl ./build/
	tar -cvzf gef-0.2.0.tar.gz build/*
	rm -rf build

dependencies: $(GITHUBSRC)/golang/lint/golint \
			  gefserver/vendor \
			  $(WEBUI)/node_modules

$(WEBUI)/node_modules:
	(cd $(WEBUI) && npm install)

$(GITHUBSRC)/golang/lint/golint:
	go get -u github.com/golang/lint/golint

gefserver/vendor:
	@echo "-- Installing go dependencies"
	(cd gefserver && glide install)

$(JSBUNDLE):
	(cd $(WEBUI) && node_modules/webpack/bin/webpack.js -p)

.PHONY: build backend webui run_gef run_webui_dev_server certificate dependencies clean pack
