package main

import (
	"flag"
	"log"
	"github.com/EUDAT-GEF/GEF/backend-docker/dckr"
	"github.com/EUDAT-GEF/GEF/backend-docker/server"
	"github.com/EUDAT-GEF/GEF/backend-docker/config"
)

var configFilePath = "config/config.json"

func main() {
	flag.StringVar(&configFilePath, "config", configFilePath, "configuration file")
	flag.Parse()

	//var err error
	settings, err := config.ReadConfigFile(configFilePath)
	if err != nil {
		log.Fatal("FATAL while reading config files: ", err)
	}
	if len(settings.Docker) == 0 {
		log.Fatal("FATAL: empty docker configuration list:\n", settings)
	}

	client, err := dckr.NewClientFirstOf(settings.Docker)
	if err != nil {
		log.Print(err)
		log.Fatal("Failed to make any docker client, exiting")
	}

	server := server.NewServer(settings.Server, client)
	log.Println("Starting GEF server at: ", settings.Server.Address)
	err = server.Start()
	if err != nil {
		log.Println("GEF server failed: ", err)
	}
}


