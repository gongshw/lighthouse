package conf

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Configuration struct {
	StaicFileDir  string
	ServerBaseUrl string
	ServerPort    int
}

var CONFIG Configuration

var inited bool

var (
	ERROR_CONF_INITED = errors.New("configuration inited already")
	ERROR_LOAD_CONF   = errors.New("can't load configuration file")
)

func InitConfig(configFilePath string) error {
	if inited {
		return ERROR_CONF_INITED
	}
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
	inited = true
	return nil
}

func tryConfigFilePath(configFilePath string) (string, error) {
	if configFilePath != "" {
		configFilePath, _ = filepath.Abs(configFilePath)
		if _, err := os.Stat(configFilePath); err == nil {
			log.Println("found log file at " + configFilePath)
			return configFilePath, nil
		} else {
			return "", err
		}
	} else {
		posibleConfigFiles := []string{}
		if p, e := os.Getwd(); e == nil {
			// loop for lighthouse.json from current working directory
			posibleConfigFiles = append(posibleConfigFiles, p+"/lighthouse.json")
		}
		if u, e := user.Current(); e == nil {
			// loop for lighthouse.json from current user's home directory
			posibleConfigFiles = append(posibleConfigFiles, u.HomeDir+"/lighthouse.json")
		}
		for _, path := range posibleConfigFiles {
			if _, err := os.Stat(path); err == nil {
				log.Println("found log file at " + path)
				return path, nil
			}
		}
		return "", ERROR_LOAD_CONF
	}
}
