package conf

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

const _5MB = 5 * 1024 * 1024

type Configuration struct {
	StaicFileDir          string
	ServerBaseUrl         string
	ServerPort            int
	ResponseTimeoutSecond time.Duration
	FilterMode            string
	FilterFile            string
	ContentLengthLimit    int64
}

func validateConfig() error {
	// TODO
	return nil
}

var CONFIG Configuration = Configuration{
	ContentLengthLimit: _5MB,
}

var inited bool

var (
	ERROR_LOAD_CONF = errors.New("can't load configuration file")
)

func LoadConfig(configFilePath string) error {
	if inited {
		return nil
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
	if validateErr := validateConfig(); validateConfig() != nil {
		return validateErr
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
				log.Println("load configuration file from " + path)
				return path, nil
			}
		}
		return "", ERROR_LOAD_CONF
	}
}
