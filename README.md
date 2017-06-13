GEF
===

The EUDAT Generic Execution Framework.

[![Build Status](https://travis-ci.org/EUDAT-GEF/GEF.svg?branch=master)](https://travis-ci.org/EUDAT-GEF/GEF)

Installation
------------

- Make sure you have `go 1.8` (the language tools), `docker`, and `npm` installed on your machine. For MacOS it is
recommended to have `go 1.8.1` (or newer)
- Set a GOPATH, e.g.: `export GOPATH=/Users/myself/Projects/Go`.
- Use `go get` to clone the GEF repository: `go get -u github.com/EUDAT-GEF/GEF`.
- Go to the downloaded repository location: `cd $GOPATH/src/github.com/EUDAT-GEF/GEF`.
- Build the project: `make build`.
- Create a new self-signed certificate for the GEF server (with `make certificate`) or edit config.json to use your own
- Define the GEF_SECRET_KEY environment variable (a random string is preferred, remember it and keep it in a safe location): `export GEF_SECRET_KEY="E60su8IL6VY6Ca2"`
- Start the system: `make run_gef`.
- Go to `https://localhost:8443`. The GEF UI should be online.

Internal Docker Images
-------------
When GEF connects to a Docker server it builds several custom images. These are necessary for the system to function
properly, please do not remove them.

GEF Demo Services
-------------
In the root of the repository there is a folder called `services`: it contains GEF demo services. They are not essential
for the system and therefore are not available by default (after the installation process). However, they can help to better 
understand how to build your own GEF services and how they work. If you want to use those services, you can build them
through the web interface (by uploading corresponding Dockerfiles and image-related files).

Database
-------------
The GEF is using an SQLite database to store the data. Using a different SQL database is possible.

Docker Swarm Mode
-------------
If you want to run GEF services on a Docker Swarm, you will need to create and configure it first. There is no need to install anything
else, as long as you have Docker installed. `config.json` file has to be modified accordingly: endpoint should be changed
to `tcp://[IP_ADDRESS_OF_THE_MANAGER_MACHINE]:[PORT_WHERE_DOCKER_IS_RUNNING]` and `TLSVerify`, `CertPath`,
`KeyPath`, `CAPath` should be set in the `Docker` section of the config file.
- This tutorial explains how to create a swarm: https://rominirani.com/docker-swarm-tutorial-b67470cf8872
- Follow the instructions at: https://docs.docker.com/engine/security/https/ to create certificates or use the existing `tls.sh` script
(which is much more convenient). 
If you create a virtual machine in VirtualBox, you should configure the network properly:
- Open `VirtualBox Manager`
- Right click on a virtual machine and select `Settings`
- Go to `Network` tab
- First network adapter should be attached to `NAT`, the second one - to `Host-only adapter` (with the name `vboxnet1`)
- Configure port forwarding for the first adapter similar to the following table:

| Name | Protocol | Host IP | Host Port | Guest IP | Guest Port |
| :---: |:--------:| :------:| :-------: | :------: | :---------: |
| Daemon | TCP | 127.0.0.1 | 59047 | 10.0.2.15 | 2376 |
| swarm | TCP | 127.0.0.1 | 59046 | 10.0.2.15 | 3024 |
| ssh | TCP | 127.0.0.1 | 59045 |  | 22 |
