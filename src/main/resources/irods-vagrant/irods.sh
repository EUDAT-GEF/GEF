#!/bin/bash -x

if [ ! -e /home/vagrant/.irodsprovisioned ]; then
    apt-get update
    apt-get upgrade -y
    apt-get install -q -y curl build-essential python-pip git python-dev postgresql odbc-postgresql unixodbc-dev
 
	wget ftp://ftp.renci.org/pub/irods/releases/4.0.0/irods-icat-4.0.0-64bit.deb
	wget ftp://ftp.renci.org/pub/irods/releases/4.0.0/irods-database-plugin-postgres-1.0.deb
#	wget ftp://ftp.renci.org/pub/irods/releases/4.0.0/irods-resource-4.0.0-64bit.deb
#	wget ftp://ftp.renci.org/pub/irods/releases/4.0.0/irods-dev-4.0.0-64bit.deb
#	wget ftp://ftp.renci.org/pub/irods/releases/4.0.0/irods-runtime-4.0.0-64bit.deb
#	wget ftp://ftp.renci.org/pub/irods/releases/4.0.0/irods-icommands-4.0.0-64bit.deb

	dpkg -i *.deb
	apt-get -f install -y

#    cp /vagrant/irods.config $IRODS_DIR/config/

#    chown -R vagrant:vagrant $IRODS_DIR
#    chmod -R a+rx $IRODS_DIR

#     Use system-wide postgresql for iRODS
#    echo "CREATE ROLE vagrant PASSWORD 'md5ce5f2d27bc6276a03b0328878c1dc0e2' SUPERUSER CREATEDB CREATEROLE INHERIT LOGIN;" | su - postgres -c psql -
#    psql createuser does not allow password via cmdline
#    su postgres -c "createuser vagrant -s -w"

#    cp /vagrant/pg_hba.conf /etc/postgresql/9.1/main/pg_hba.conf
#    ln -sf /usr/lib/x86_64-linux-gnu/odbc/psqlodbca.so /usr/lib/postgresql/9.1/lib/libodbcpsql.so
#    service postgresql restart

#    su vagrant -c "cd $IRODS_DIR && export USE_LOCALHOST=1 && ./scripts/configure && make && ./scripts/finishSetup --noask"

#    su vagrant -c "mkdir -p /home/vagrant/.irods"

#    echo "export PATH=\$PATH:$IRODS_DIR:$IRODS_DIR/clients/icommands/bin" >> /home/vagrant/.bashrc

    touch /home/vagrant/.irodsprovisioned
fi
