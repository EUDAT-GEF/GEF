FROM ubuntu:16.04

MAINTAINER Alexandr Chernov <kstchernov@gmail.com>

LABEL "eudat.gef.service.name"="NLTK"
LABEL "eudat.gef.service.description"="NLTK is a leading platform for building Python programs to work with human language data"
LABEL "eudat.gef.service.version"="1.0"
LABEL "eudat.gef.service.input.1.name"="Input Directory"
LABEL "eudat.gef.service.input.1.path"="/root/input"
LABEL "eudat.gef.service.output.1.name"="Output Directory"
LABEL "eudat.gef.service.output.1.path"="/root/output"

RUN apt-get update
RUN apt-get install -y python-pip
RUN pip install -U nltk
RUN pip install -U numpy

RUN mkdir /root/input
RUN mkdir /root/output
#ADD nltk.py /root
ADD test.py /root
ADD test.csv /root
RUN python -m nltk.downloader all

CMD ["python", "/root/test.py"]