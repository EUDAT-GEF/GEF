FROM ubuntu:16.04

MAINTAINER Alexandr Chernov <kstchernov@gmail.com>

LABEL "eudat.gef.service.name"="NLTK POS-tagging"
LABEL "eudat.gef.service.description"="Performs text segmentation (splits into sentences) and POS-tagging"
LABEL "eudat.gef.service.version"="1.0"
LABEL "eudat.gef.service.input.1.name"="First Input Directory"
LABEL "eudat.gef.service.input.1.path"="/root/input1"
LABEL "eudat.gef.service.input.2.name"="Second Input Directory"
LABEL "eudat.gef.service.input.2.path"="/root/input2"
LABEL "eudat.gef.service.output.1.name"="Output Directory"
LABEL "eudat.gef.service.output.1.path"="/root/output"

RUN apt-get update
RUN apt-get install -y python-pip
RUN pip install -U nltk
RUN pip install -U numpy

RUN mkdir /root/input1
RUN mkdir /root/input2
RUN mkdir /root/output

ADD posTagger.py /root

RUN python -m nltk.downloader punkt
RUN python -m nltk.downloader averaged_perceptron_tagger

CMD ["python", "/root/posTagger.py"]