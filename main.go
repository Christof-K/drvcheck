package main

import (
	config "drvcheck/config"
	drvcheck "drvcheck/app"
	interactive "drvcheck/interactive"
	"flag"
	"os"
	"strings"
)

const MODE_EXEC = "exec"
const MODE_DEV = "dev"

func main() {

	flag_mode := flag.String("mode", "exec", "Mode exec or dev (override config path)")
	flag.Parse()


	switch *flag_mode {
		case MODE_DEV:
			config.Conf.PreConfig.SetYamlConfigPath(".")
		// case MODE_EXEC: nothing
	}


	// get args params skipping flags
	var args []string
	for _, arg := range os.Args[1:] {
		if strings.Split(arg, "")[0] != "-" {
			args = append(args, arg)
		}
	}

	if len(args) == 0 {
		drvcheck.Run()
	} else {
		for _, arg := range args {
			if strings.Split(arg, "")[0] == "-" {
				continue // skip flags
			}
			switch(arg) {
				case "interactive":
					interactive.Run()
				case "interactive-test":
					interactive.RunTestMode = true
					interactive.Run()
				default:
					drvcheck.Run()
			}
		}
	}

}

func GetConfig() {
	
}
