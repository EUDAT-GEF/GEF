GEF
===

The EUDAT Generic Execution Framework.

[![Build Status](https://travis-ci.org/EUDAT-GEF/GEF.svg?branch=master)](https://travis-ci.org/EUDAT-GEF/GEF)

Installation
------------

1. Make sure you have `go 1.8` (the language tools), `docker`, and `npm` installed on your machine. For MacOS it is
recommended to have `go 1.8.1` (or newer)
2. Set a GOPATH, e.g.: `export GOPATH=/Users/myself/Projects/Go`.
3. Use `go get` to clone the GEF repository: `go get -u github.com/EUDAT-GEF/GEF`.
4. Go to the downloaded repository location: `cd $GOPATH/src/github.com/EUDAT-GEF/GEF`.
5. Build the project: `make build`.
6. Create a new self-signed certificate for the GEF server (with `make certificate`) or edit config.json to use your own
7. Define the GEF_SECRET_KEY environment variable (a random string is preferred, remember it and keep it in a safe location): `export GEF_SECRET_KEY="E60su8IL6VY6Ca2"`
8. Start the system: `make run_gef`.
9. Go to `https://localhost:8443`. The GEF UI should be online.

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
To active the Swarm mode on a local machine run `docker swarm init --advertise-addr 127.0.0.1`. Executing `docker swarm leave -f`
will turn it off. If you want to run GEF services on a Docker Swarm, you will need to create and configure it first. There is no need to install anything
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


Packaging For Deployment
-------------
For the packaging procedure you will only need `go` and `docker` to be installed on your machine. Having executed the first
4 steps from the `Installation` section, you can build a linux binary by running `make pack` command. The command will
create an archive which you can unpack on a server. After that do `cd build/bin` and run `./gefserver`. The GEF server
should be running now.

Apache Configuration
-------------
If you are using Apache, you may want to use a configuration similar to the one below:
~~~~
Listen 443
NameVirtualHost *:443
<VirtualHost *:443>
    SSLEngine on
    SSLProxyEngine on
    SSLProxyVerify none
    SSLProxyCheckPeerCN off
    SSLProxyCheckPeerName off
    SSLProxyCheckPeerExpire off
    ProxyRequests off
    SSLCertificateFile /etc/apache2/ssl/ssl.crt
    SSLCertificateKeyFile /etc/apache2/ssl/ssl.key
    ProxyPreserveHost On
    ProxyPass / https://127.0.0.1:8443/
    ProxyPassReverse / https://127.0.0.1:8443/
</VirtualHost>
~~~~