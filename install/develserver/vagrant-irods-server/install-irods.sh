#!/bin/bash

IRODS_FTP=ftp://ftp.renci.org/pub/irods/releases/4.0.0
PROVISIONED=/home/vagrant/.irodsprovisioned

if [ ! -e  $PROVISIONED ]; then
    apt-get update
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

    # done
    chown -R vagrant:vagrant /home/vagrant
    touch $PROVISIONED

    echo "!!!"
    echo "RUN THE NEXT INTERACTIVE COMMAND TO REALLY FINISH THE PROVISIONING"
    echo "-- When prompted for hostname or IP, use 127.0.0.1 NOT localhost"
    echo "$ sudo su - irods -c /var/lib/irods/packaging/setup_database.sh"
    echo
    echo "For manually restarting the irods server, use:"
    echo "$ sudo su - irods -c '/var/lib/irods/iRODS/irodsctl restart'"
fi
