#!/bin/sh
# This file must be deployed in the iRODS server/bin/cmd directory
# and will get executed via the WSO 

echo "geffilter inparams: " $*
echo "geffilter pwd: " `pwd`

STAGE_DIR="$1"
FILE="$2"
FILTER="$3"

cd "$STAGE_DIR"
echo "geffilter pwd2: " `pwd`
ls -al

mkdir out
grep $FILTER $FILE > out/filter.out.txt 
echo
