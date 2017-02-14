FROM alpine:latest

MAINTAINER Alexandr Chernov <kstchernov@gmail.com>

LABEL "eudat.gef.service.name"="Volume file list"
LABEL "eudat.gef.service.description"="Provides a list of files and folders in a volume"
LABEL "eudat.gef.service.version"="0.1"

ADD volume-filelist /root
RUN mkdir /root/volume

CMD ["/root/volume-filelist"]