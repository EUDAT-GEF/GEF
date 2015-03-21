#!/bin/bash -x

IRODS_FTP=ftp://ftp.renci.org/pub/irods/releases/4.0.0

if [ ! -e /home/vagrant/.irodsprovisioned ]; then
    apt-get update
    apt-get upgrade -y
    apt-get install -q -y curl build-essential python-pip git python-dev postgresql odbc-postgresql unixodbc-dev libssl0.9.8 super
 
	wget -nv $IRODS_FTP/irods-icat-4.0.0-64bit.deb
	wget -nv $IRODS_FTP/irods-database-plugin-postgres-1.0.deb

	dpkg -i *.deb
	apt-get -f install -y

    adduser irods sudo
    cp /vagrant/irods.config /etc/irods/

    # Use system-wide postgresql for iRODS
    echo "CREATE ROLE irods PASSWORD 'irodsgef' superuser createdb createrole inherit login;" | sudo su - postgres -c psql -
    cp /vagrant/pg_hba.conf /etc/postgresql/9.3/main/pg_hba.conf
    ln -sf /usr/lib/x86_64-linux-gnu/odbc/psqlodbca.so /usr/lib/postgresql/9.3/lib/libodbcpsql.so
    service postgresql restart

    # # install java 
    # add-apt-repository -y ppa:webupd8team/java
    # apt-get update
    # echo oracle-java8-installer shared/accepted-oracle-license-v1-1 select true | /usr/bin/debconf-set-selections
    # apt-get install -y oracle-java8-installer
    # export JAVA_OPTS="-Djava.awt.headless=true -Xmx1g"
    # export JAVA_HOME=/usr/lib/jvm/java-8-oracle
    # ln -s /usr/lib/jvm/java-8-oracle /usr/lib/jvm/default-java

    # tomcat (not used, use local tomcat instead)
    # apt-get install tomcat7 tomcat7-admin -y

    # docker
    apt-get install -q -y docker.io
    usermod -a -G docker vagrant
    ln -sf /usr/bin/docker.io /usr/local/bin/docker

    cd /vagrant
    for dir in docker-*; do
        cp -r $dir /home/vagrant/
        docker build -t $dir $dir
    done

    # 
    cp /vagrant/gefcommand /var/lib/irods/iRODS/server/bin/cmd/

    # done
    touch /home/vagrant/.irodsprovisioned
fi

# !!!
# RUN THE NEXT INTERACTIVE COMMAND TO REALLY FINISH THE PROVISIONING
# -- When prompted for hostname or IP, use 127.0.0.1 NOT localhost
# sudo su - irods -c /var/lib/irods/packaging/setup_database.sh

# For manually restarting the irods server, use:
# sudo su - irods -c /var/lib/irods/iRODS/irodsctl restart
