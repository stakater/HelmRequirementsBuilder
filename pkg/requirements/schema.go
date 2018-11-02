package requirements

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/stakater/RequirementsUpdater/pkg/log"
)

var (
	logger = log.New()
)

type Requirements struct {
	Dependencies []Dependency `yaml:"dependencies"`
}

type Dependency struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	Repository string `yaml:"repository"`
	Alias      string `yaml:"alias"`
}

func Read(filePath string) (*Requirements, error) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var r Requirements

	err = yaml.Unmarshal(yamlFile, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func Write(filePath string, reqs *Requirements) error {
	yamlFile, err := yaml.Marshal(reqs)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, yamlFile, 0644)
}
