FROM ubuntu:14.04
MAINTAINER Emanuel Dima <emanueldima@gmail.com>
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -q -y wget libssl0.9.8

RUN wget -nv ftp://ftp.renci.org/pub/irods/releases/4.0.0/irods-icommands-4.0.0-64bit.deb
RUN dpkg -i *.deb
RUN apt-get -f install -y
#RUN cp /vagrant/irods.config /etc/irods/

#RUN add-apt-repository ppa:webupd8team/java -y
#RUN apt-get update

# !!!
# this also needs input for license, might not work automatically
#RUN apt-get install oracle-java8-installer -y
