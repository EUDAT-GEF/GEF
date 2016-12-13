# Base image: Minimal Alpine Linux image
FROM alpine:latest

MAINTAINER Alexandr Chernov <kstchernov@gmail.com>

LABEL "eudat.gef.service.name"="Copy data to a volume"
LABEL "eudat.gef.service.description"="Copies data from a specific folder in a container to a given volume"
LABEL "eudat.gef.service.version"="0.1"

ADD copy-to-volume /root
RUN mkdir /root/volume
RUN mkdir /root/buffer

# run the application
CMD ["/root/copy-to-volume"]