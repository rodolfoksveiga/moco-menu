package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/rodolfoksveiga/moco-menu/pkg/menu"
)

type Init struct {
	ConfigPath string
}

func (init Init) Load() *Config {
	configFile, err := os.Open(init.ConfigPath)
	menu.ExitOnError(err, "Could not find config file.")
	defer configFile.Close()

	configBytes, err := ioutil.ReadAll(configFile)
	menu.ExitOnError(err, "Could not read config file.")

	var config Config
	err = json.Unmarshal(configBytes, &config)
	menu.ExitOnError(err, "Could not map config file.")

	return &config
}
