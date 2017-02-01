FROM ubuntu:14.04
MAINTAINER Emanuel Dima <emanueldima@gmail.com>

ADD test.txt /test.txt

LABEL "eudat.gef.service.name"="Test Clone"
LABEL "eudat.gef.service.description"="Test Copy input to output"
LABEL "eudat.gef.service.version"="0.1"
LABEL "eudat.gef.service.input.1.name"="Input Directory"
LABEL "eudat.gef.service.input.1.path"="/mydata/input"
LABEL "eudat.gef.service.output.1.name"="Output Directory"
LABEL "eudat.gef.service.output.1.path"="/mydata/output"

CMD ["cp", "/mydata/input/*", "/mydata/output/"]