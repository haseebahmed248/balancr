// Config File parsing
package config

import (
	"balancr/internal/logger"
	"balancr/internal/pool"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type configFile struct {
	Backends []*pool.Backend `yaml:"backends"`
}

func GetBackends() []*pool.Backend {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		logger.Log("Error loading the yaml file", "ERROR")
		return nil
	}
	var cfg configFile

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logger.Log("Error parsing the yaml", "ERROR")
		return nil

	}
	return cfg.Backends
}
