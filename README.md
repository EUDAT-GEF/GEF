GEF
===

The EUDAT Generic Execution Framework.

[![Build Status](https://travis-ci.org/EUDAT-GEF/GEF.svg?branch=master)](https://travis-ci.org/EUDAT-GEF/GEF)

Installation
------------

- Make sure you have `go` (the language tools), `java`, `maven`, and `docker` installed on your machine.
- Set a GOPATH, e.g.: `export GOPATH=/Users/myself/Projects/Go`.
- Use `go get` to clone the GEF repository: `go get -u github.com/EUDAT-GEF/GEF`.
- Go to the downloaded repository location: `cd $GOPATH/src/github.com/EUDAT-GEF/GEF`.
- Build the project: `make build`.
- Start the frontend and the backend in two separate terminal sessions: `make run_frontend` and `make run_backend`.
- Go to `http://localhost:4042`. The GEF UI should be online.

Docker Images
-------------
During the installation process several docker images will be created: they are necessary for the system to function
properly, please do not remove them. 