package def

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration keeps the configs for the entire application
type Configuration struct {
	Docker DockerConfig
	// Default Docker limits
	Limits LimitConfig

	// Default GEF internal timeouts
	Timeouts TimeoutConfig

	Pier        PierConfig
	Server      ServerConfig
	EventSystem EventSystemConfig

	// TmpDir is the directory to keep session files in
	// If the path is relative, it will be used as a subfolder of the system temporary directory
	TmpDir string
}

// DockerConfig configuration for building docker clients
type DockerConfig struct {
	Description string
	Endpoint    string
	TLSVerify   bool
	CertPath    string
	KeyPath     string
	CAPath      string
}

// PierConfig configuration for pier
type PierConfig struct {
	InternalServicesFolder string
}

// ServerConfig keeps the configuration options needed to make a Server
type ServerConfig struct {
	Address                string
	ReadTimeoutSecs        int
	WriteTimeoutSecs       int
	TLSCertificateFilePath string
	TLSKeyFilePath         string
	B2Access               B2AccessConfig
	B2Drop                 B2DropConfig
}

// B2AccessConfig exported
type B2AccessConfig struct {
	BaseURL     string
	RedirectURL string
}

// B2DropConfig exported
type B2DropConfig struct {
	BaseURL string
}

// EventSystemConfig keeps the configuration options needed to make an EventSystem
type EventSystemConfig struct {
	Address string
}

// LimitConfig keeps the configuration options to limit resources used by a docker container while its execution
type LimitConfig struct {
	CPUShares  int64
	CPUPeriod  int64
	CPUQuota   int64
	Memory     int64
	MemorySwap int64
}

type TimeoutConfig struct {
	DataStaging      int64
	VolumeInspection int64
	FileDownload     int64
	Preparation      int64
	JobExecution     int64
}

func (c DockerConfig) String() string {
	tls := ""
	if c.TLSVerify {
		tls = "with TLS"
	}
	return fmt.Sprintf("%s %s -- %s", c.Endpoint, tls, c.Description)
}

// ReadConfigFile reads a configuration file
func ReadConfigFile(configFilepath string) (Configuration, error) {
	var config Configuration

	file, err := os.Open(configFilepath)
	if err != nil {
		return config, Err(err, "Cannot open config file %s", configFilepath)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, Err(err, "Cannot read config file %s", configFilepath)
	}

	if config.Docker.Endpoint == "" {
		return config, Err(nil, "Incorrect Docker endpoint in file: %s", configFilepath)
	}

	return config, nil
}
