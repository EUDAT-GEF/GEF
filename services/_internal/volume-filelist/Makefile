SOURCES := $(shell find . -iname '*.go')

build: volume-filelist

volume-filelist: $(SOURCES)
	$(GOPATH)/bin/golint ./...
	go vet ./...
	GOOS=linux GOARCH=amd64 go build
	docker build -t volume-filelist .

clean:
	go clean
