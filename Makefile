ASSETDIR = src/main/resources/assets
LIBDIR = $(ASSETDIR)/lib
FONTDIR = $(ASSETDIR)/fonts
JSDIR = $(ASSETDIR)/js

build: jsx target

bower:
	echo "---- npm & bower"
	mkdir -p $LIBDIR
	mkdir -p $FONTDIR
	mkdir -p $JSDIR
	echo
	npm install bower react-tools
	node_modules/bower/bin/bower install jquery bootstrap react react-addons react-router font-awesome
	echo
	cp bower_components/bootstrap/dist/css/bootstrap.min.css $LIBDIR/
	cp bower_components/bootstrap/dist/js/bootstrap.min.js $LIBDIR/
	cp bower_components/jquery/dist/jquery.min.js $LIBDIR/
	cp bower_components/jquery/dist/jquery.min.map $LIBDIR/
	cp bower_components/react/react-with-addons.js $LIBDIR/
	cp bower_components/react/react-with-addons.min.js $LIBDIR/
	cp bower_components/react-router/dist/react-router.min.js $LIBDIR/
	cp bower_components/font-awesome/css/font-awesome.min.css $LIBDIR/
	echo
	cp bower_components/bootstrap/fonts/*  $FONTDIR/
	cp bower_components/font-awesome/fonts/* $FONTDIR/

jsx:
	./buildjsx.sh

target:
	mvn -q package

run:
	@echo "The GEF front end depends on an iRODS server + an gef-docker server!"
	@$(eval JAR = $(shell find target -iname 'GEF-*.jar'))
	@java -cp src/main/resources:$(JAR) eu.eudat.gef.app.GEF server gefconfig.yml

run_production:
	@$(eval JAR = $(shell find target -iname 'GEF-*.jar'))
	@java -jar $(JAR) server gefconfig.yml

clean:
	go clean
	mvn -q clean

.PHONY: build bower jsx run run_production clean
