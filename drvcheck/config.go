package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/kardianos/osext"
	"github.com/goccy/go-yaml"
)

type configYaml struct {
	Drivers []string `yaml:",flow"`
	Unit string
	Csv struct {
		Mode string
		Dir string
		Header []string `yaml:",flow"`
	}
}

type PreConfig struct {
	overrideYamlConfigPath string
}

func (pc *PreConfig) SetYamlConfigPath(path string) {
	pc.overrideYamlConfigPath = path
}


type config struct {
	isLoaded bool
	configYaml *configYaml
	PreConfig *PreConfig
}

var Conf = &config{
	isLoaded: false,
	configYaml: &configYaml{},
	PreConfig: &PreConfig{},
}

func GetConfig() (config, error) {
	
	if Conf.isLoaded != true {
		fmt.Println("Reading configuration....")

		var path string
		if Conf.PreConfig.overrideYamlConfigPath == "" {
			tmp_path, _ := osext.ExecutableFolder()
			path = tmp_path
		} else {
			path = Conf.PreConfig.overrideYamlConfigPath
			fmt.Println("overriding config path with", path)
		}


		cyaml, err := os.ReadFile(filepath.Join(path, "config.yaml"))
		if err != nil { 
			fmt.Println(err.Error())
			return *Conf, err
		}
		uerr := yaml.Unmarshal(cyaml, Conf.configYaml)
		if uerr != nil {
			fmt.Println(uerr.Error())
			return *Conf, uerr
		}

		Conf.configYaml.Csv.Dir = path
		Conf.isLoaded = true
	}

	return *Conf, nil
}
