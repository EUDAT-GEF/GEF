FROM ubuntu:rolling

MAINTAINER Asela Rajapakse <asela.rajapakse@mpimet.mpg.de>

# Metadata labels

LABEL "eudat.gef.service.name"="CDO Demo"
LABEL "eudat.gef.service.description"="Demo service for calling Climate Data Operators on NetCDF data."
LABEL "eudat.gef.service.version"="1.0.0"

LABEL "eudat.gef.service.input.1.name"="Input Directory"
LABEL "eudat.gef.service.input.1.path"="/input"
LABEL "eudat.gef.service.input.1.type"="url"
LABEL "eudat.gef.service.output.1.name"="Output Directory"
LABEL "eudat.gef.service.output.1.path"="/output"

# Installing Spack package manager

RUN apt-get update && apt-get install -y \
    apt-utils \
    python \
    git \
    gfortran \
    curl \
    build-essential \
 && rm -rf /var/lib/apt/lists/* \
 && git clone https://github.com/spack/spack \
 && cd spack \
 && git checkout releases/v0.11.0 

ENV PATH=$PATH:/spack/bin
ENV PATH=$PATH:/spack/opt/spack/linux-ubuntu17.10-x86_64/gcc-7.2.0/cdo-1.9.1-4gqys4cpbwwkgvwyf5xp7c7a45ehe56l/bin
RUN ["/bin/bash", "-c", ". /spack/share/spack/setup-env.sh"]

# Installing Climate Data Operators

RUN ["/bin/bash", "-c", "spack install cdo~szip~hdf5~udunits2~libxml2~proj~fftw3+netcdf grib2=none ^netcdf~mpi ^hdf5~mpi"]

# Creating input and output directories

RUN mkdir /input
RUN mkdir /output

# Copy CDO invocation script into the image

COPY ./cdo_gather_lwp.sh /

ENTRYPOINT ["/bin/bash", "cdo_gather_lwp.sh"]
