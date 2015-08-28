package conf

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Configuration struct {
	StaicFileDir  string
	ServerBaseUrl string
	ServerPort    int
}

var CONFIG Configuration

func InitConfig(configFilePath string) error {
	var err error
	var configFile io.Reader
	configFilePath, err = tryConfigFilePath(configFilePath)
	if err != nil {
		return err
	}
	configFile, err = os.Open(configFilePath)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&CONFIG)
	if err != nil {
		return err
	}
	return nil
}

func tryConfigFilePath(configFilePath string) (string, error) {
	if _, err := os.Stat(configFilePath); err == nil {
		log.Println("Found log file at " + configFilePath)
		return configFilePath, nil
	} else {
		return "", err
	}
}
