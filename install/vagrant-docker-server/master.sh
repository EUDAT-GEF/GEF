#!/bin/bash

IRODS_FTP=ftp://ftp.renci.org/pub/irods/releases/4.0.0
PROVISIONED=/home/vagrant/.provisioned

if [ ! -e  $PROVISIONED ]; then
    apt-get update
    apt-get install -q -y curl build-essential python-pip git python-dev libssl0.9.8 super

    echo; echo "--- install docker"
    # install newer docker (need labels and volume plugins)
    apt-key adv --keyserver hkp://pgp.mit.edu:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
    echo "deb https://apt.dockerproject.org/repo ubuntu-trusty main" | tee /etc/apt/sources.list.d/docker.list
    apt-get update
    apt-get purge -y lxc-docker* docker*
    apt-get install -q -y docker-engine

    usermod -a -G docker vagrant

    echo; echo "--- testing docker"
    docker run busybox ls -al /

    echo; echo "--- install irodsFs"
    # install irods client
    # this irods package contains the 'irodsFs' utility
	wget -nv $IRODS_FTP/irods-icat-4.0.0-64bit.deb

	dpkg -i *.deb
	apt-get -f install -y

    cp -r /vagrant/.irods  /home/vagrant/.irods

    # mounting irods storage with fuse (for user root)
    # clear text p4ssw0rd here, only for development!
    cp -r /vagrant/.irods  /root/.irods
    sudo iinit rodsgef
    sudo mkdir /data
    sudo irodsFs -o allow_other,ro /data

    echo; echo "--- testing docker with irods volume"
    docker run -v /data:/data_1:ro busybox ls -al /data_1

    echo; echo "--- done"
    chown -R vagrant:vagrant /home/vagrant
    touch $PROVISIONED

    if [ -e /vagrant/gef-docker ]; then
        echo; echo "--- starting gef-docker"
        cp /vagrant/gef-docker /vagrant/config.json /home/vagrant/
        chown -R vagrant:vagrant /home/vagrant
        /home/vagrant/gef-docker
    fi
fi

# # install java
# add-apt-repository -y ppa:webupd8team/java
# apt-get update
# echo oracle-java8-installer shared/accepted-oracle-license-v1-1 select true | /usr/bin/debconf-set-selections
# apt-get install -y oracle-java8-installer
# export JAVA_OPTS="-Djava.awt.headless=true -Xmx1g"
# export JAVA_HOME=/usr/lib/jvm/java-8-oracle
# ln -s /usr/lib/jvm/java-8-oracle /usr/lib/jvm/default-java
