ASSETDIR = src/main/resources/assets
LIBDIR = $(ASSETDIR)/lib
FONTDIR = $(ASSETDIR)/fonts
JSDIR = $(ASSETDIR)/js

build: install jsx target

install:
	(cd src/main/webui && npm install)

jsx:
	(cd src/main/webui && node_modules/webpack/bin/webpack.js --config webpack.config.devel.js -d)

target:
	mvn -q package

run:
	@echo "The GEF front end depends on an iRODS server + an gef-docker server!"
	@$(eval JAR = $(shell find target -iname 'GEF-*.jar'))
	@java -cp src/main/resources:$(JAR) eu.eudat.gef.app.GEF server gefconfig.yml

run_production:
	(cd src/main/webui && npm install)
	(cd src/main/webui && node_modules/webpack/bin/webpack.js -p)
	@$(eval JAR = $(shell find target -iname 'GEF-*.jar'))
	@java -jar $(JAR) server gefconfig.yml

clean:
	go clean
	mvn -q clean

.PHONY: build install jsx run run_production clean
