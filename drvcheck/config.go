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


type config struct {
	isLoaded bool
	configYaml *configYaml
}

var conf = &config{
	isLoaded: false,
	configYaml: &configYaml{},
}

func GetConfig() (config, error) {
	
	if conf.isLoaded != true {
		fmt.Println("Reading configuration....")
		path, _ := osext.ExecutableFolder()
		cyaml, err := os.ReadFile(filepath.Join(path, "config.yaml"))
		if err != nil { 
			fmt.Println(err.Error())
			return *conf, err
		}
		uerr := yaml.Unmarshal(cyaml, conf.configYaml)
		if uerr != nil {
			fmt.Println(uerr.Error())
			return *conf, uerr
		}

		// switch to current executable directory
		if conf.configYaml.Csv.Dir == "." {
			path, _ := osext.ExecutableFolder()
			conf.configYaml.Csv.Dir = path
		}

		conf.isLoaded = true
	}

	return *conf, nil
}
