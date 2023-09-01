package templates

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/rodolfoksveiga/moco-menu/pkg/menu"
)

type Init struct {
	TemplatePath string
}

func (init Init) Load() *[]Template {
	templatesFile, err := os.Open(init.TemplatePath)
	menu.ExitOnError(err, "Could not find template file.")
	defer templatesFile.Close()

	configBytes, err := ioutil.ReadAll(templatesFile)
	menu.ExitOnError(err, "Could not read template file.")

	var templates []Template
	err = json.Unmarshal(configBytes, &templates)
	menu.ExitOnError(err, "Could not map template file.")

	return &templates
}
