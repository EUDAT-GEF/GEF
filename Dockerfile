FROM ubuntu:14.04
MAINTAINER Emanuel Dima <emanueldima@gmail.com>

LABEL "eudat.gef.service.name"="Clone"
LABEL "eudat.gef.service.description"="Copy input to output"
LABEL "eudat.gef.service.version"="0.1"
LABEL "eudat.gef.service.input.1.name"="Text Input"
LABEL "eudat.gef.service.input.1.path"="/data/input1/"
LABEL "eudat.gef.service.output.1.name"="Text Output"
LABEL "eudat.gef.service.output.1.path"="/data/output1/"

COPY gef-service-example /gef-service-example
CMD ["/gef-service-example"]
