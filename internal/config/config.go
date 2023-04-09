package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	ConfigFilePath string
}

func (config Config) Load() (*AuthConfig, *string) {
	configFile, err := os.Open(config.ConfigFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		errMsg := "Could not find config file."
		return nil, &errMsg
	}
	defer configFile.Close()

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		fmt.Println("Error:", err)
		errMsg := "Could not read config file."
		return nil, &errMsg
	}

	var authConfig AuthConfig
	err = json.Unmarshal(configBytes, &authConfig)
	if err != nil {
		fmt.Println("Error:", err)
		errMsg := "Could not map config file."
		return nil, &errMsg
	}

	return &authConfig, nil
}
