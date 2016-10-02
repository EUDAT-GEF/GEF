SOURCES := $(shell find . -iname '*.go')

build: gef-docker

gef-docker: $(SOURCES)
	golint ./...
	go vet ./...
	# go test ./...
	cd ./src
# update the packages
	GOOS=linux GOARCH=amd64 go install github.com/eudat-gef/gef-docker/dckr
	GOOS=linux GOARCH=amd64 go install github.com/eudat-gef/gef-docker/server
# install the bin
	GOOS=linux GOARCH=amd64 go install github.com/eudat-gef/gef-docker/gef-docker


clean:
	go clean
