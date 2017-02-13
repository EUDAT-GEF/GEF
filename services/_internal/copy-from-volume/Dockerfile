# Base image: Minimal Alpine Linux image
FROM alpine:latest

MAINTAINER Alexandr Chernov <kstchernov@gmail.com>

LABEL "eudat.gef.service.name"="Copy from a volume"
LABEL "eudat.gef.service.description"="Copies files from a running container to a volume"
LABEL "eudat.gef.service.version"="0.1"

RUN mkdir /root/volume

CMD ["cp"]
