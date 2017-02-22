FROM ubuntu:16.04

MAINTAINER Alexandr Chernov <kstchernov@gmail.com>

LABEL "eudat.gef.service.name"="Stanford Parser for English"
LABEL "eudat.gef.service.description"="Parses a given text and produces constituency and dependency trees for each sentence"
LABEL "eudat.gef.service.version"="1.0"
LABEL "eudat.gef.service.input.1.name"="Input Directory"
LABEL "eudat.gef.service.input.1.path"="/root/input"
LABEL "eudat.gef.service.output.1.name"="Output Directory"
LABEL "eudat.gef.service.output.1.path"="/root/output"

RUN apt-get update
RUN apt-get install -y default-jdk curl unzip

RUN curl -Ls http://nlp.stanford.edu/software/stanford-parser-full-2016-10-31.zip > /root/stanford-parser-full-2016-10-31.zip
RUN unzip /root/stanford-parser-full-2016-10-31.zip -d /root
RUN rm /root/stanford-parser-full-2016-10-31.zip
RUN mkdir /root/input
RUN mkdir /root/output

CMD ["/root/stanford-parser-full-2016-10-31/lexparser.sh", "/root/input/*.txt"]