GEF
===

The EUDAT Generic Execution Framework.

[![Build Status](https://travis-ci.org/EUDAT-GEF/GEF.svg?branch=master)](https://travis-ci.org/EUDAT-GEF/GEF)


Server Deployment
-----------------
- Get the latest GEF binary release from https://github.com/EUDAT-GEF/GEF/releases
- Unpack the tar file in a dedicated folder, and set it as the current directory
- Inspect and change the default server configuration in the  `./config.json` file
- Create a new self-signed certificate for the GEF server (with `make certificate`) or edit `config.json` to use your own
- Define the GEF_SECRET_KEY environment variable (a random string is preferred, remember it and keep it in a safe location): e.g. `export GEF_SECRET_KEY="E60su8IL6VY6Ca2"`
- Define the GEF_B2ACCESS_CONSUMER_KEY and GEF_B2ACCESS_SECRET_KEY environment variables for connection to the B2ACCESS service.
- Run `./gefserver`


Installation for development
----------------------------

- Make sure you have the following tools installed on your machine:
    - the `go 1.8` language tools: https://golang.org. (For MacOS it is recommended to have at least `go 1.8.1`)
    - the `glide` tool for go: https://github.com/Masterminds/glide
    - `docker`: https://www.docker.com
    - `npm`: https://www.npmjs.com

- Set a GOPATH, e.g.: `export GOPATH=/Users/myself/Projects/Go`.
- Use `go get` to clone the GEF repository: `go get -u github.com/EUDAT-GEF/GEF`.
- Go to the downloaded repository location: `cd $GOPATH/src/github.com/EUDAT-GEF/GEF`.
- Build the project: `make build`.
- Create a new self-signed certificate for the GEF server (with `make certificate`) or edit `config.json` to use your own
- Define the GEF_SECRET_KEY environment variable (a random string is preferred, remember it and keep it in a safe location): e.g. `export GEF_SECRET_KEY="E60su8IL6VY6Ca2"`
- Define the GEF_B2ACCESS_CONSUMER_KEY and GEF_B2ACCESS_SECRET_KEY environment variables for connection to the B2ACCESS service.
- Start the system: `make run_gef`.
- Go to `https://localhost:8443`. The GEF UI should be online.

Internal Docker Images
----------------------
When GEF connects to a Docker server it builds several custom images. These are necessary for the system to function
properly, please do not remove them.

GEF Demo Services
-----------------
In the root of the repository there is a folder called `services`: it contains GEF demo services. They are not essential
for the system and therefore are not available by default (after the installation process). However, they can help to better
understand how to build your own GEF services and how they work. If you want to use those services, you can build them
through the web interface (by uploading corresponding Dockerfiles and image-related files).

Database
--------
The GEF is using an SQLite database to store the data. Using a different SQL database is possible.

Docker Swarm Mode
-----------------
If you want to run GEF services on a Docker Swarm, you will need to create and configure it first. There is no need to install anything else,
as long as you have Docker installed. To active the Swarm mode on a local machine run `docker swarm init --advertise-addr 127.0.0.1`.
Executing `docker swarm leave -f` will turn it off. `config.json` file has to be modified accordingly: endpoint should be changed
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
------------------------
For the packaging procedure you will only need to have the GEF source code and `docker` to be installed on your machine. You can build a linux binary
by running `make pack` command. The command will create an archive which you can unpack on a server. After that do `cd build/bin` and run `./gefserver`.
The GEF server should be running now.

Apache Configuration
--------------------
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

API description
--------------------
| URL | Method | Input | Output | Description |
| ---: |:-------- | :------ | :------- | :------ |
| /api/info | GET |  | API version information in JSON | Information about API (welcome page), can be used to check if backend is running |
| /api/user | GET |  | JSON with the information about the current user | Returns information about the current user |
| /api/user/tokens | POST | Form data with the name of a token {tokenName} | JSON with the new token | Adds a new token for the current user |
| /api/user/tokens | GET |  | JSON with the list of all user tokens | List all tokens for the current user |
| /api//user/tokens/{tokenID} | DELETE | {tokenID} an id of a token | Server response code | Removes a specific token from the current user |
| /api/roles | GET |  | JSON with the list of all roles | Lists all available roles |
| /api/roles/{roleID} | GET | {roleID} an id of a role | JSON with the list of users | Returns a list of users to which a certain role was assigned |
| /api/roles/{roleID} | POST | {roleID} an id of a role | Server response code | Assigns a specific role to the current user |
| /api/roles/{roleID}/{userID} | DELETE | {roleID} an id of a role assigned to the user with the user id {userID} | Server response code | Removes a role from a user |
| /api/builds | POST |  | JSON object with information about the location and build ID | Creates a temporary folder when an image has to be created. It returns a buildID identifier and a folder location. This folder is used to store a Dockerfile and files needed for the image. BuildID is a string like a UID in Java (which is generated when required and it is unique) |
| /api/builds/{buildID} | POST | {buildID} build identifier | JSON object with information about the image and the corresponding service | Builds an image provided a buildID (which points to the folder with the Dockerfile), returns JSON with the information about the image and the new service (partly taken from the metadata) |
| /api/services | GET |  | JSON with the list of all services | Lists all available services  |
| /api/services/{serviceID} | GET | {serviceID} an id of a service | JSON with information about a specific service | Returns information about a specific service |
| /api/services/{serviceID} | PUT | {serviceID} an id of a service, form data with new service metadata | JSON with information about a specific service | Modifies metadata of a specific service |
| /api/services/{serviceID} | DELETE | {serviceID} an id of a service | JSON with service information | Deletes a specific job |
| /api/jobs | POST | serviceID and pid | JSON object with information about the location and jobID | Executes a job. We submit a form to this URL when we want to run a job. PID is resolved, files are downloaded and saved to an input volume |
| /api/jobs | GET |  | JSON with the list of jobs | Lists available jobs |
| /api/jobs/{jobID} | GET | {jobID} id of a job | JSON with job information | Information about a specific job |
| /api/jobs/{jobID} | DELETE | {jobID} id of a job | JSON with job information | Deletes a specific job |
| /api/j/volumes/{volumeID}/{path:.*} | GET | {volumeID} is an id of a volume, {path} is a path inside this volume (root folder by default) | JSON object (nested) with the list of the files and folders in a given volume | Lists all files and folders (recursively) in a given volume |
