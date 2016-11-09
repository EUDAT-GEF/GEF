# Base image: Minimal Alpine Linux image with Java
FROM alpine:latest

MAINTAINER Alexandr Chernov <kstchernov@gmail.com>

LABEL "eudat.gef.service.name"="Volume inspector"
LABEL "eudat.gef.service.description"="Provides a list of files and folders in a given folder"
LABEL "eudat.gef.service.version"="0.1"

# specify the port number the container should expose
EXPOSE 8282

ADD volume-inspector /root

# run the application
CMD ["/root/volume-inspector"]