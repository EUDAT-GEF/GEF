#!/bin/bash
cd `dirname $0`

echo " --- Starting irods virtual machine"
(cd vagrant-irods-server && vagrant up)
if [ $? -ne 0 ]; then
	echo "irods vm failed, exiting"
	exit 1
fi

echo " --- Starting docker virtual machine"
(cd vagrant-docker-server && vagrant up)
if [ $? -ne 0 ]; then
	echo "docker vm failed, exiting"
	exit 1
fi

./onchanged-gef-docker.sh

if [ "$1" == "-w" ]; then
	echo " --- Start watching ../gef-docker"
	fswatch -0 -o ../gef-docker/gef-docker | xargs -0 -n 1 -I {} ./onchanged-gef-docker.sh
fi
