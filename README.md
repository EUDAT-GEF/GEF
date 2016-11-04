GEF
===

The EUDAT Generic Execution Framework.

[![Build Status](https://travis-ci.org/EUDAT-GEF/GEF.svg?branch=master)](https://travis-ci.org/EUDAT-GEF/GEF)

Installation
------------

1. Make sure you have go (the language tools) installed on your machine.
2. Set a GOPATH, e.g.: `export GOPATH=/Users/myself/Projects/Go`.
3. Use `go get` to clone the GEF repository: `go get github.com/EUDAT-GEF/GEF`.
4. Go to the downloaded repository location: `cd $GOPATH/src/github.com/EUDAT-GEF/GEF`.
5. Build the project: `make build`.
6. Start the frontend and the backend in two separate terminal sessions: `make run_frontend` and `make run_backend`.
7. Go to `http://localhost:4042/gef`. The GEF UI should be online.
