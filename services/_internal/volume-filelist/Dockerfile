FROM golang:latest

MAINTAINER Alexandr Chernov <kstchernov@gmail.com>

LABEL "eudat.gef.service.name"="Volume file list"
LABEL "eudat.gef.service.description"="Provides a list of files and folders in a volume"
LABEL "eudat.gef.service.version"="0.1"

WORKDIR /root
ADD main.go /root
RUN go build main.go && mv main volume-filelist && rm main.go
RUN mkdir /root/volume

CMD ["/root/volume-filelist"]