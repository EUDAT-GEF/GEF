FROM ubuntu:latest
MAINTAINER Xavier Pivan <xavier.pivan.ds@gmail.com>

ENV HOME /root

LABEL "eudat.gef.service.name"="nco-lib"
LABEL "eudat.gef.service.description"="Execute nco command on ncdf" 
LABEL "eudat.gef.service.version"="1.0"
LABEL "eudat.gef.service.input.1.name"="Input Directory" 
LABEL "eudat.gef.service.input.1.path"="/root/input"
LABEL "eudat.gef.service.input.1.type"="url"
LABEL "eudat.gef.service.output.1.name"="Output Directory" 
LABEL "eudat.gef.service.output.1.path"="/root/output"

RUN apt-get update && apt-get -yq install gcc \
								build-essential \
								wget \
								libhdf5-serial-dev \
								tar \
								bzip2 \
								m4 \
								zlib1g-dev \
								libopenmpi-dev



RUN wget ftp://ftp.hdfgroup.org/HDF5/current/src/hdf5-1.10.0-patch1.tar.bz2;
RUN wget ftp://ftp.unidata.ucar.edu/pub/netcdf/netcdf-4.4.1.1.tar.gz;

#Build HDF5

RUN tar xjvf hdf5-1.10.0-patch1.tar.bz2 && \
	cd hdf5-1.10.0-patch1 && \
    CC=mpicc ./configure --enable-parallel --prefix=/usr/local/hdf5 && \
    make -j4 && \
    make install && \
    cd .. && \
    rm -rf /hdf5-1.10.0-patch1 /hdf5-1.10.0-patch1.tar.bz2

RUN tar xzvf netcdf-4.4.1.1.tar.gz && \
	cd netcdf-4.4.1.1 && \
    ./configure --prefix=/usr/local/netcdf \ 
                CC=mpicc \
                LDFLAGS=-L/usr/local/hdf5/lib \
                CFLAGS=-I/usr/local/hdf5/include && \
    make -j4 && \
    make install && \
    cd .. && \
    rm -rf netcdf-4.4.1.1 netcdf-4.4.1.1.tar.gz

RUN apt-get -yq install nco

RUN mkdir /root/input
RUN mkdir /root/output

CMD ncwa -O -a time /root/input/*.bin /root/output/result.bin




