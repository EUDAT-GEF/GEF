FROM ubuntu:14.04
MAINTAINER Alexandr Chernov <kstchernov@gmail.com>

LABEL "eudat.gef.service.name"="Test multiple inputs"
LABEL "eudat.gef.service.description"="Test Copy input to output"
LABEL "eudat.gef.service.version"="0.1"
LABEL "eudat.gef.service.input.1.name"="First Input Directory"
LABEL "eudat.gef.service.input.1.path"="/root/input1"
LABEL "eudat.gef.service.input.2.name"="Second Input Directory"
LABEL "eudat.gef.service.input.2.path"="/root/input2"
LABEL "eudat.gef.service.output.1.name"="First Output Directory"
LABEL "eudat.gef.service.output.1.path"="/root/output1"
LABEL "eudat.gef.service.output.2.name"="Second Output Directory"
LABEL "eudat.gef.service.output.2.path"="/root/output2"

ADD copy.sh /root/copy.sh
RUN mkdir /root/input1
RUN mkdir /root/input2
RUN mkdir /root/output1
RUN mkdir /root/output2
CMD ["/bin/bash", "/root/copy.sh"]
