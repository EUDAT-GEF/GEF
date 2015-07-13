FROM ubuntu:14.04
MAINTAINER Emanuel Dima <emanueldima@gmail.com>

# update system and install required packages
RUN apt-get update
RUN apt-get install -q -y software-properties-common vim wget libssl0.9.8

# install and configure irods
RUN wget -nv ftp://ftp.renci.org/pub/irods/releases/4.0.0/irods-icommands-4.0.0-64bit.deb && \
	dpkg -i *.deb && \
	apt-get -f install -y && \
	rm irods-icommands-4.0.0-64bit.deb
COPY .irods/  /.irods
