#!/bin/bash
cd `dirname $0`

GEF_DOCKER_DIR="/Users/wqiu/Projects/GEF/gef-docker"

echo " --- replacing gef-docker in vm : " `date "+%H:%M:%S"`

cd vagrant-docker-server
vagrant ssh-config > .sshconfig

ssh -F .sshconfig vagrant@default 'killall gef-docker'
echo "killing old version of gef-docker"
scp -F .sshconfig "$GEF_DOCKER_DIR/bin/linux_amd64/gef-docker" vagrant@default:/home/vagrant/
scp -F .sshconfig "$GEF_DOCKER_DIR/config.json" vagrant@default:/home/vagrant/
ssh -F .sshconfig vagrant@default '/home/vagrant/gef-docker' &
echo " --- ... done"
echo
