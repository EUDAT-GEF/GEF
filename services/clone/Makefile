SOURCES := $(shell find . -iname '*.go')

build: clone

clone: $(SOURCES)
	#golint ./...
	go vet ./...
	GOOS=linux GOARCH=amd64 go build

run:
	# ---- run devel
	./clone

clean:
	go clean