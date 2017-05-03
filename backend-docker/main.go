package main

import (
	"flag"
	"log"

	"github.com/EUDAT-GEF/GEF/backend-docker/db"
	"github.com/EUDAT-GEF/GEF/backend-docker/def"
	"github.com/EUDAT-GEF/GEF/backend-docker/pier"
	"github.com/EUDAT-GEF/GEF/backend-docker/server"
)

var configFilePath = "config.json"

func main() {

	flag.StringVar(&configFilePath, "config", configFilePath, "configuration file")
	flag.Parse()

	config, err := def.ReadConfigFile(configFilePath)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}

	d, err := db.InitDb()
	if err != nil {
		log.Fatal("FATAL: ", def.Err(err, "Cannot initialize the database engine"))
	}

	defer d.Db.Close()

	var p *pier.Pier
	p, err = pier.NewPier(config.Docker, config.TmpDir, config.Limits, &d)
	if err != nil {
		log.Fatal("FATAL: ", def.Err(err, "Cannot create Pier"))
	}

	server.InitEventSystem(config.EventSystem.Address)
	server, err := server.NewServer(config.Server, p, config.TmpDir, &d)
	if err != nil {
		log.Fatal("FATAL: ", def.Err(err, "Cannot create API server"))
	}

	log.Println("Starting GEF server at: ", config.Server.Address)
	err = server.Start()
	if err != nil {
		log.Fatal("FATAL: ", def.Err(err, "Cannot start API server"))
	}
}
