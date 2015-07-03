SOURCES := $(shell find . -iname '*.go')

build: gef-service-example

gef-service-example: $(SOURCES)
	golint ./...
	go vet ./...
	go build

run:
	# ---- run devel
	./gef-service-example

clean:
	go clean
