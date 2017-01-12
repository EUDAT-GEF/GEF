SOURCES := $(shell find . -iname '*.go')

build: volume-filelist

volume-filelist: $(SOURCES)
	$(GOPATH)/bin/golint ./...
	go vet ./...
	#go build ./...
	#go install github.com/EUDAT-GEF/GEF/services/volume-filelist
	# install the bin
	GOOS=linux GOARCH=amd64 go build
	docker build -t eudatgef/volume-filelist .

clean:
	go clean
