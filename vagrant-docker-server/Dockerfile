FROM ubuntu:14.04
MAINTAINER Emanuel Dima <emanueldima@gmail.com>

RUN apt-get update && \
    apt-get install -q -y wget build-essential python-pip git python-dev postgresql odbc-postgresql unixodbc-dev libssl0.9.8 super fuse

RUN wget -nv ftp://ftp.renci.org/pub/irods/releases/4.0.0/irods-icat-4.0.0-64bit.deb && \
    dpkg -i *.deb && \
    apt-get -f install -y && \
    rm irods-icat-4.0.0-64bit.deb

COPY .irods /root/.irods
RUN iinit rodsgef

CMD ["/bin/bash"]
