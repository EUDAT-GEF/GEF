SOURCES := $(shell find . -iname '*.go')

build: gef-docker

gef-docker: $(SOURCES)
	golint ./...
	go vet ./...
	go build

run:
	# ---- run devel
	# (cd vagrant && vagrant up)
	./gef-docker

clean:
	go clean
