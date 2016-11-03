#!/bin/bash
cd `dirname $0`

echo " --- Starting docker virtual machine"
(cd vagrant-docker-server && vagrant up)
if [ $? -ne 0 ]; then
	echo "docker vm failed, exiting"
	exit 1
fi

./onchanged-gef-docker.sh

if [ "$1" == "-w" ]; then
	echo " --- Start watching ../gef-docker"
	fswatch -0 -o ../gef-docker/bin/linux_amd64/gef-docker -docker | xargs -0 -n 1 -I {} ./onchanged-gef-docker.sh
fi
