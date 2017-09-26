package main

import (
	"flag"
	"log"

	"github.com/EUDAT-GEF/GEF/gefserver/db"
	"github.com/EUDAT-GEF/GEF/gefserver/def"
	"github.com/EUDAT-GEF/GEF/gefserver/pier"
	"github.com/EUDAT-GEF/GEF/gefserver/server"
)

var configFilePath = "config.json"

func main() {
	flag.StringVar(&configFilePath, "config", configFilePath, "configuration file")
	flag.Parse()

	config, err := def.ReadConfigFile(configFilePath)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	if config.Limits.CPUPeriod <= 0 {
		log.Fatal("FATAL: ", def.Err(nil, "CPUPeriod is not set in the config file"))
	}

	d, err := db.InitDb()
	if err != nil {
		log.Fatal("FATAL: ", def.Err(err, "Cannot initialize the database engine"))
	}
	defer d.Close()

	var p *pier.Pier
	p, err = pier.NewPier(&d, config.Pier, config.TmpDir, config.Timeouts)
	if err != nil {
		log.Fatal("FATAL: ", def.Err(err, "Cannot create Pier"))
	}

	// add the main docker connection with no owner (UserID = 0)
	_, err = p.AddDockerConnection(0, config.Docker)
	if err != nil {
		log.Fatal("FATAL: ", def.Err(err, "Cannot set docker connection"))
	}

	server.InitEventSystem(config.EventSystem.Address)
	server, err := server.NewServer(config, p, &d)
	if err != nil {
		log.Fatal("FATAL: ", def.Err(err, "Cannot create API server"))
	}

	log.Println("Starting GEF server at: ", config.Server.Address)
	err = server.Start()
	if err != nil {
		log.Fatal("FATAL: ", def.Err(err, "Cannot start API server"))
	}
}
