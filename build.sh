#!/bin/bash

ASSETDIR=src/main/resources/assets
LIBDIR=$ASSETDIR/lib
FONTDIR=$ASSETDIR/fonts
JSDIR=$ASSETDIR/js

if [ ! -e bower_components ]
then
	echo; echo "---- npm & bower"

	mkdir -p $LIBDIR
	mkdir -p $FONTDIR
	mkdir -p $JSDIR

	npm install bower react-tools
	node_modules/bower/bin/bower install jquery bootstrap react react-addons react-router font-awesome

	cp bower_components/bootstrap/dist/css/bootstrap.min.css $LIBDIR/
	cp bower_components/bootstrap/dist/js/bootstrap.min.js $LIBDIR/
	cp bower_components/jquery/dist/jquery.min.js $LIBDIR/
	cp bower_components/jquery/dist/jquery.min.map $LIBDIR/
	cp bower_components/react/react-with-addons.js $LIBDIR/
	cp bower_components/react/react-with-addons.min.js $LIBDIR/
	cp bower_components/react-router/dist/react-router.min.js $LIBDIR/
	cp bower_components/font-awesome/css/font-awesome.min.css $LIBDIR/

	cp bower_components/bootstrap/fonts/*  $FONTDIR/
	cp bower_components/font-awesome/fonts/* $FONTDIR/
fi

echo; echo "---- jsx"
for f in $JSDIR/*.jsx; do 
	cp -v $f $JSDIR/`basename $f .jsx`.js; 
done
node_modules/react-tools/bin/jsx --no-cache-dir $JSDIR $JSDIR

echo; echo "---- go build"
go build src/executor/gefcommand.go

echo; echo "---- mvn clean package"
mvn -q clean package

# Run in production:
# java -jar target/Aggregator2-2.0.0-beta-x.jar server aggregator.yml

# Run for development:
# java -cp src/main/resources:target/Aggregator2-2.0.0-alpha-10.jar eu.clarin.sru.fcs.aggregator.app.Aggregator server aggregator_development.yml

