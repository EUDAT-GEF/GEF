SOURCES := $(shell find . -iname '*.go')

build: volume-inspector

build-linux: volume-inspector-linux

gef-container-server-linux: $(SOURCES)
	#golint ./...
	#go vet ./...
	# go test ./...
	cd ./src
# install the bin
	GOOS=linux GOARCH=amd64 go install github.com/EUDAT-GEF/GEF/services/volume-inspector

gef-container-server: $(SOURCES)
	#golint ./...
	#go vet ./...
	# go test ./...
	cd ./src

# install the bin
	go install github.com/EUDAT-GEF/GEF/services/volume-inspector


clean:
	go clean
