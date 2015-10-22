#!/bin/bash
cd `dirname $0`

echo " --- replacing gef-docker in vm : " `date "+%H:%M:%S"`

cd vagrant-docker-server
vagrant ssh-config > .sshconfig

ssh -F .sshconfig vagrant@default 'killall gef-docker'
scp -F .sshconfig ../../gef-docker/gef-docker vagrant@default:/home/vagrant/
scp -F .sshconfig ../../gef-docker/config.json vagrant@default:/home/vagrant/
ssh -F .sshconfig vagrant@default '/home/vagrant/gef-docker' &
echo " --- ... done"
echo
