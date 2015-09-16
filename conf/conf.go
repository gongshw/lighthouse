package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
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
	SSLCertificationFile  string
	SSLKeyFile            string
	DisableSSL            bool
}

type ConfigValidateError struct {
	errorMsgs []string
}

func (e ConfigValidateError) Error() string {
	return "fail to load configuration: " + strings.Join(e.errorMsgs, ";")
}

func (e *ConfigValidateError) Append(msg string) {
	e.errorMsgs = append(e.errorMsgs, msg)
}

func (e *ConfigValidateError) HasError() bool {
	return len(e.errorMsgs) > 0
}

func validateConfig() error {
	e := &ConfigValidateError{}

	if CONFIG.StaicFileDir == "" {
		e.Append("StaicFileDir is empty")
	} else if !pathValid(CONFIG.StaicFileDir) {
		e.Append(fmt.Sprintf("StaicFileDir(\"%s\") is unaccessable", CONFIG.StaicFileDir))
	}

	if CONFIG.ServerPort <= 0 || CONFIG.ServerPort > 65535 {
		e.Append(fmt.Sprintf("ServerPort(%d) is illegal", CONFIG.ServerPort))
	}

	if CONFIG.FilterMode != "" {
		if strings.EqualFold(CONFIG.FilterMode, "white") &&
			strings.EqualFold(CONFIG.FilterMode, "black") {
			e.Append(fmt.Sprintf("FilterMode(\"%s\") is illegal", CONFIG.StaicFileDir))
		}
		if !pathValid(CONFIG.FilterFile) {
			e.Append(fmt.Sprintf("FilterFile(\"%s\") is unaccessable", CONFIG.StaicFileDir))
		}
	}

	if e.HasError() {
		return e
	} else {
		return nil
	}
}

func pathValid(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

var CONFIG Configuration = Configuration{
	ContentLengthLimit:    _5MB,
	ResponseTimeoutSecond: 5,
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
	if validateErr := validateConfig(); validateErr != nil {
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
				log.Println("found configuration file at " + path)
				return path, nil
			}
		}
		return "", ERROR_LOAD_CONF
	}
}
