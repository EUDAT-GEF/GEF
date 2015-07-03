package main

import (
	"encoding/json"
	"os"

	"../eudock"
)

type configuration struct {
	Docker []eudock.Config
	Server ServerConfig
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
