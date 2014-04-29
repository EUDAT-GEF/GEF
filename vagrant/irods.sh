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

    sudo adduser irods sudo
    cp /vagrant/irods.config /etc/irods/

    # Use system-wide postgresql for iRODS
    echo "CREATE ROLE irods PASSWORD 'irodsgef' superuser createdb createrole inherit login;" | sudo su - postgres -c psql -
    cp /vagrant/pg_hba.conf /etc/postgresql/9.3/main/pg_hba.conf
    ln -sf /usr/lib/x86_64-linux-gnu/odbc/psqlodbca.so /usr/lib/postgresql/9.3/lib/libodbcpsql.so
    service postgresql restart

    # java 
    add-apt-repository ppa:webupd8team/java -y
    apt-get update
    # !!! 
    # this also needs input for license, might not work automatically
    apt-get install oracle-java8-installer -y 

    # tomcat
    export JAVA_HOME=/usr/lib/jvm/java-8-oracle
    ln -s /usr/lib/jvm/java-8-oracle /usr/lib/jvm/default-java
    export JAVA_OPTS="-Djava.awt.headless=true -Xmx1g"
    apt-get install tomcat7 tomcat7-admin -y

    # done
    touch /home/vagrant/.irodsprovisioned
fi

# !!!
# RUN THE NEXT INTERACTIVE COMMAND TO FINISH THE PROVISIONING
# -- When prompted for hostname or IP, use 127.0.0.1 NOT localhost
# sudo su - irods -c /var/lib/irods/packaging/setup_database.sh 
