SOURCES := $(shell find . -iname '*.go')
VAGRANT := ./install/develserver/vagrant-docker-server

build: gef-docker

gef-docker: $(SOURCES)
	golint ./...
	go vet ./...
	# go test ./...
	GOOS=linux GOARCH=amd64 go build

install: gef-docker
	cp ./gef-docker ./config.json $(VAGRANT)

clean:
	go clean
	rm -f $(VAGRANT)/gef-docker $(VAGRANT)/config.json
