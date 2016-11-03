#!/bin/bash

ASSETDIR=src/main/resources/assets
JSDIR=$ASSETDIR/js

echo; echo "---- jsx"
for f in $JSDIR/*.jsx; do
	cp -v $f $JSDIR/`basename $f .jsx`.js;
done
node_modules/react-tools/bin/jsx --no-cache-dir $JSDIR $JSDIR
