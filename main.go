package main

import (
	"flag"

	"../eudock"
)
import "log"

var configFilePath = "config.json"

var config configuration

func main() {
	flag.StringVar(&configFilePath, "config", configFilePath, "configuration file")
	flag.Parse()

	var err error
	config, err = readConfigFile(configFilePath)
	if err != nil {
		log.Fatal("FATAL while reading config files: ", err)
	}
	if len(config.Docker) == 0 {
		log.Fatal("FATAL: empty docker configuration list:\n", config)
	}

	client, err := eudock.NewClientFirstOf(config.Docker)
	if err != nil {
		log.Print(err)
		log.Fatal("Failed to make any docker client, exiting")
	}

	server := NewServer(config.Server, client)
	log.Println("Starting GEF server at: ", config.Server.Address)
	server.Start()
}
