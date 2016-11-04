SOURCES := $(shell find . -iname '*.go')

build: gef-container-server

build-linux: gef-container-server-linux

gef-container-server-linux: $(SOURCES)
	#golint ./...
	#go vet ./...
	# go test ./...
	cd ./src
# install the bin
	GOOS=linux GOARCH=amd64 go install github.com/gef-container-server

gef-container-server: $(SOURCES)
	#golint ./...
	#go vet ./...
	# go test ./...
	cd ./src

# install the bin
	go install github.com/gef-container-server


clean:
	go clean
