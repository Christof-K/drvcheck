package drvcheck

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/kardianos/osext"
	"gopkg.in/yaml.v2"
)

type configYaml struct {
	Drivers []string `yaml:"drivers,flow"`
	Unit string `yaml:"unit"`
	Csv CsvStruct `yaml:"csv"`
}

type CsvStruct struct {
	Mode string `yaml:"mode"`
	Dir string `yaml:"dir"`
	Header []string `yaml:"header,flow"`
}

type PreConfig struct {
	overrideYamlConfigPath string
}

func (pc *PreConfig) SetYamlConfigPath(path string) {
	pc.overrideYamlConfigPath = path
}


type Config struct {
	isLoaded bool
	ConfigYaml *configYaml
	PreConfig *PreConfig
}

var Conf = &Config{
	isLoaded: false,
	ConfigYaml: &configYaml{},
	PreConfig: &PreConfig{},
}

func GetConfig() (Config, error) {
	
	if !Conf.isLoaded {
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
		
		uerr := yaml.Unmarshal(cyaml, Conf.ConfigYaml)

		if uerr != nil {
			fmt.Println(uerr.Error())
			return *Conf, uerr
		}

		// Conf.configYaml.Csv.Dir = path
		Conf.isLoaded = true
	}

	return *Conf, nil
}
