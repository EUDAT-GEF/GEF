GEF
===

The EUDAT Generic Execution Framework.

[![Build Status](https://travis-ci.org/EUDAT-GEF/GEF.svg?branch=master)](https://travis-ci.org/EUDAT-GEF/GEF)

Installation
------------

- Make sure you have `go 1.8` (the language tools), `docker`, and `npm` installed on your machine.
- Set a GOPATH, e.g.: `export GOPATH=/Users/myself/Projects/Go`.
- Use `go get` to clone the GEF repository: `go get -u github.com/EUDAT-GEF/GEF`.
- Go to the downloaded repository location: `cd $GOPATH/src/github.com/EUDAT-GEF/GEF`.
- Build the project: `make build`.
- Create a new self-signed certificate for the GEF server (with `make certificate`) or edit config.json to use your own
- Define the GEF_SECRET_KEY environment variable (a random string is preferred, remember it and keep it in a safe location): `export GEF_SECRET_KEY="E60su8IL6VY6Ca2"`
- Start the system: `make run_gef`.
- Go to `https://localhost:8443`. The GEF UI should be online.

Docker Images
-------------
When GEF connects to a Docker server it builds several custom images. These are necessary for the system to function properly, please do not remove them.

Database
-------------
The GEF is using an SQLite database to store the data. Using a different SQL database is possible.
