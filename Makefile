SOURCES := $(shell find . -iname '*.go')

build: gef-service-example

gef-service-example: $(SOURCES)
	golint ./...
	go vet ./...
	GOOS=linux GOARCH=amd64 go build

run:
	# ---- run devel
	./gef-service-example

clean:
	go clean
