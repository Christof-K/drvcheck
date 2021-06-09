package helper

import (
	"os"
	"github.com/goccy/go-yaml"
)

type configYaml struct {
	Drivers []string `yaml:",flow"`
	Csv struct {
		Mode string
		Dir string
	}
}


type config struct {
	isLoaded bool
	configYaml *configYaml
}

var conf = &config{
	isLoaded: false,
	configYaml: &configYaml{},
}

func GetConfig() (*config, error) {
	
	if conf.isLoaded != true {
		cyaml, err := os.ReadFile("./config.yaml")
		if err != nil { return nil, err }
		uerr := yaml.Unmarshal(cyaml, conf.configYaml)
		if uerr != nil { return nil, uerr }
		conf.isLoaded = true
	}

	return conf, nil
}
