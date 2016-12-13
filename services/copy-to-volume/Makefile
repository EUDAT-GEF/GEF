SOURCES := $(shell find . -iname '*.go')

build: copy-to-volume

copy-to-volume: $(SOURCES)
	$(GOPATH)/bin/golint ./...
	go vet ./...
	GOOS=linux GOARCH=amd64 go build
	docker build -t eudatgef/copy-to-volume .

clean:
	go clean
