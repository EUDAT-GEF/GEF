#!/bin/bash

BUILD=1
RUN=

while [[ $# > 0 ]]
do
key="$1"
# echo $# " :" $key
case $key in
    --build)
    BUILD=1
    ;;
    --run)
    RUN=1
    ;;
    *)
    echo "Unknown option:" $1
    exit 1
    ;;
esac
shift
done

if [ $BUILD ]
then
	echo; echo "---- go lint"
	golint ./...
    echo; echo "---- go vet"
	go vet ./...
    echo; echo "---- go build"
	go build
fi

if [ $RUN ]
then
	echo; echo "---- run devel"
	# echo "-- vagrant up"
	# (cd vagrant && vagrant up)
	echo "-- ./gef-docker"
	"./gef-docker"
fi
