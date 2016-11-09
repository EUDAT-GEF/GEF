SOURCES := $(shell find . -iname '*.go')

build: volume-inspector

build-linux: volume-inspector-linux

volume-inspector-linux: $(SOURCES)
	$(GOPATH)/bin/golint ./...
	go vet ./...
	go test ./api/api_test.go
# install the bin
	GOOS=linux GOARCH=amd64 go install github.com/EUDAT-GEF/GEF/services/volume-inspector

volume-inspector: $(SOURCES)
	$(GOPATH)/bin/golint ./...
	go vet ./...
	go test ./api/api_test.go

# install the bin
	go install github.com/EUDAT-GEF/GEF/services/volume-inspector

clean:
	go clean
