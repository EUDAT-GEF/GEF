FROM ubuntu:14.04
MAINTAINER Emanuel Dima <emanueldima@gmail.com>

# update system and install required packages
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -q -y software-properties-common vim wget libssl0.9.8

# install and configure irods
RUN wget -nv ftp://ftp.renci.org/pub/irods/releases/4.0.0/irods-icommands-4.0.0-64bit.deb
RUN dpkg -i *.deb
RUN apt-get -f install -y
RUN rm irods-icommands-4.0.0-64bit.deb
ADD .irods/.irodsEnv  /.irods/.irodsEnv
ADD .irods/.irodsA    /.irods/.irodsA
