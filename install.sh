#!/bin/bash

echo "# 1. Build the source code"
echo "make"
echo
echo "# 2. Provision the irods server with vagrant"
echo "# after install you may need to manually run additional commands"
echo "# connect to the machine and make sure the irods server works"
echo "(cd install/develserver/vagrant-irods-server && vagrant up)"
echo
echo "# Now provision the docker server with vagrant and start the gef server"
echo "(cd install/develserver/vagrant-docker-server && vagrant up)"
