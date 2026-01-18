// Config File parsing
package config

import (
	"balancr/internal/pool"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type configFile struct {
	Backends []*pool.Backend `yaml:"backends"`
}

func GetBackends() []*pool.Backend {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Error loading the yaml file")
	}
	var cfg configFile

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("Error parsing the yaml")
	}
	return cfg.Backends
}
