package config

import (
	"flag"
	"github.com/goccy/go-yaml"
	"github.com/qrest/gomisc/serror"
	"log"
	"os"
)

// ReadConfig reads the config file from the given file path
func ReadConfig(configFilePath string, config interface{}) error {
	// Open config file
	file, err := os.Open(configFilePath)
	if err != nil {
		return serror.New(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(serror.New(err))
		}
	}(file)

	// Init new YAML decode
	d := yaml.NewDecoder(file, yaml.Strict())

	// Start YAML decoding from file
	if err := d.Decode(config); err != nil {
		return serror.New(err)
	}

	return nil
}

// WriteConfig writes the config value to the given file path
func WriteConfig(filePath string, config interface{}) error {
	marshalledConfig, err := yaml.Marshal(&config)
	if err != nil {
		return serror.New(err)
	}

	if err := os.WriteFile(filePath, marshalledConfig, 0600); err != nil {
		return serror.New(err)
	}

	return nil
}

// SetConfigFlags sets the CLI flags for accessing and generating the configuration file
func SetConfigFlags(defaultConfigName string, filePath *string, createConfigFile *bool) {
	flag.StringVar(filePath, "config", defaultConfigName, "config file path")
	flag.BoolVar(createConfigFile, "createConfig", false,
		"creates a default config file '"+defaultConfigName+"' (default: false)")
}
