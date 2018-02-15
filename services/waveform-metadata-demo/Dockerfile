FROM ubuntu:16.04

MAINTAINER Mathijs Koymans

# Metadata labels
LABEL "eudat.gef.service.name"="ObsPy QC Demo"
LABEL "eudat.gef.service.description"="GEF demo for seismic metadata extraction."
LABEL "eudat.gef.service.version"="1.0.0"

# Add input/output labels
LABEL "eudat.gef.service.input.1.name"="Input Directory"
LABEL "eudat.gef.service.input.1.path"="/var/input"
LABEL "eudat.gef.service.input.1.type"="url"
LABEL "eudat.gef.service.output.1.name"="Output Directory"
LABEL "eudat.gef.service.output.1.path"="/var/output"

RUN apt-get update

# Install some dependencies
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y git python-dev python-pycurl python-simplejson \
    libcurl4-gnutls-dev libssl-dev pkg-config gcc-4.9 \
    python python-psutil python-requests python-jsonschema \
    python-setuptools python-dev build-essential python-pip \
    libxml2-dev libxslt1-dev libz-dev libfreetype6-dev \
    python-numpy python-scipy python-matplotlib

# Create IO directories
RUN mkdir -p /var/input
RUN mkdir -p /var/output

# Make application directory
RUN mkdir -p /opt/collector

COPY ./collector.py /opt/collector

WORKDIR /opt/collector

# Install ObsPy for the QC module
RUN git clone https://github.com/obspy/obspy.git
WORKDIR /opt/collector/obspy
RUN pip install -e .

WORKDIR /opt/collector

ENTRYPOINT ["python", "collector.py"]
