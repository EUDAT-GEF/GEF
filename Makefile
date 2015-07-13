SOURCES := $(shell find . -iname '*.go')

build: gef-docker

gef-docker: $(SOURCES)
	golint ./...
	go vet ./...
	# go test ./...
	go build
	cp ./gef-docker ./vagrant/

clean:
	go clean
	rm ./vagrant/gef-docker
