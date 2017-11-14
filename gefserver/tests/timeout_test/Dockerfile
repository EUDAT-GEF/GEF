FROM ubuntu:16.04
MAINTAINER Alexandr Chernov <kstchernov@gmail.com>

LABEL "eudat.gef.service.name"="Executing apt-get update command"
LABEL "eudat.gef.service.description"="apt-get update keeps the container running and eventually is supposed to time out"
LABEL "eudat.gef.service.version"="0.1"
LABEL "eudat.gef.service.version"="1.0"
LABEL "eudat.gef.service.input.1.name"="Input Directory"
LABEL "eudat.gef.service.input.1.path"="/root/input"
LABEL "eudat.gef.service.input.1.type"="url"
LABEL "eudat.gef.service.output.1.name"="Output Directory"
LABEL "eudat.gef.service.output.1.path"="/root/output"

RUN mkdir /root/input
RUN mkdir /root/output

CMD ["apt-get", "update"]
