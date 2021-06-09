package helper

import (
	"io/ioutil"
	"gopkg.in/yaml.v3"
)

type configYaml struct {
	Drivers []string `yaml:"Drivers,flow"`
	Csv struct {
		Mode string `yaml:"mode"`
		Dir string `yaml:"dir"`
	}
}


type config struct {
	isLoaded bool
	configYaml
}

var conf = &config{isLoaded: false}

func getConfig() (*config, error) {
	if !conf.isLoaded {
		cyaml, err := ioutil.ReadFile("./config.yaml")
		if err != nil { return nil, err }
		cy := configYaml{}
		uerr := yaml.Unmarshal(cyaml, &cy)
		if uerr != nil { return nil, uerr }
		conf.configYaml = cy
		conf.isLoaded = true
	}
	return conf, nil
}
