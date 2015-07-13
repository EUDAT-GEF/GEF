package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/gefx/gef-docker/dckr"
	"github.com/gefx/gef-docker/server"
)

var configFilePath = "config.json"

type configuration struct {
	Docker []dckr.Config
	Server server.Config
}

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

	client, err := dckr.NewClientFirstOf(config.Docker)
	if err != nil {
		log.Print(err)
		log.Fatal("Failed to make any docker client, exiting")
	}

	server := server.NewServer(config.Server, client)
	log.Println("Starting GEF server at: ", config.Server.Address)
	server.Start()
}

func readConfigFile(configfilepath string) (configuration, error) {
	var config configuration
	file, err := os.Open(configfilepath)
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}
