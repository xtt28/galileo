package config

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

const appDataDirName = "io.github.xtt28.galileo"
const configFileName = "config.json"

// GalileoConfig is the JSON schema for a Galileo configuration file.
type GalileoConfig struct {
	// OpenAIKey is the user's OpenAI API key.
	OpenAIKey string `json:"openAIKey"`
}

// GetConfigPath returns the path to the configuration file, which is contained
// within their home or application data directory.
func GetConfigPath() string {
	confDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("could not open user config directory", err)
	}

	confPath := path.Join(confDir, appDataDirName, configFileName)
	log.Printf("user config directory: %s\n", confPath)
	return confPath
}

// ReadConfig reads the configuration file from the preset configuration
// directory and returns a GalileoConfig object with the configuration info.
func ReadConfig() (config GalileoConfig) {
	data, err := os.ReadFile(GetConfigPath())
	if err != nil {
		log.Fatal("could not read configuration file", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("could not parse configuration file", err)
	}
	log.Println("parsed configuration file")

	return
}
