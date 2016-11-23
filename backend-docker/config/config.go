package config

import (
	"os"
	"encoding/json"
	"github.com/EUDAT-GEF/GEF/backend-docker/dckr"
	"github.com/EUDAT-GEF/GEF/backend-docker/server"
)

var config configuration

type configuration struct {
	Docker []dckr.Config
	Server server.Config
}



func ReadConfigFile(configfilepath string) (configuration, error) {
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
