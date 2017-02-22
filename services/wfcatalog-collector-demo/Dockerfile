FROM ubuntu:14.04

MAINTAINER Mathijs Koymans

LABEL "eudat.gef.service.name"="WFCatalog Collector Demo"
LABEL "eudat.gef.service.description"="GEF demo for seismic metadata extraction"
LABEL "eudat.gef.service.version"="1.0.0"
LABEL "eudat.gef.service.input.1.name"="Input Directory"
LABEL "eudat.gef.service.input.1.path"="/usr/src/collector/wfcatalog/collector/archive"
LABEL "eudat.gef.service.output.1.name"="Output Directory"
LABEL "eudat.gef.service.output.1.path"="/dump"

# Temporary solution
RUN mkdir /dump

RUN apt-get update

# Install some dependencies
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y wget git vim python-dev python-pycurl python-simplejson \
    libcurl4-gnutls-dev libssl-dev pkg-config gcc-4.9 \
    python python-psutil python-requests python-jsonschema \
    python-setuptools python-dev build-essential python-pip \
    libxml2-dev libxslt1-dev libz-dev libfreetype6-dev \
    python-numpy python-scipy python-matplotlib

# Python MongoDB driver
RUN pip install pymongo

# Upgrade decorator or ObsPy cannot be imported
RUN pip install decorator --upgrade

# Make application directory
RUN mkdir -p /usr/src/collector
WORKDIR /usr/src/collector

# Install ObsPy
RUN git clone https://github.com/obspy/obspy.git
WORKDIR /usr/src/collector/obspy
RUN pip install -e .

# Get the WFCatalog Collector source
WORKDIR /usr/src/collector
RUN git clone https://github.com/EIDA/wfcatalog.git

# Copy WFCatalogCollector.py script & configation to app directory
RUN cp /usr/src/collector/wfcatalog/collector/WFCatalogCollector.py /usr/src/collector/.
RUN cp /usr/src/collector/wfcatalog/collector/config.json /usr/src/collector/.

RUN rm -rf /usr/src/collector/wfcatalog/collector/archive/*

CMD ["python", "/usr/src/collector/WFCatalogCollector.py", "--dir", "/usr/src/collector/wfcatalog/collector/archive"]