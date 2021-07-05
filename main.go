package main

import (
	"drvcheck/src/drvcheck"
	"drvcheck/src/interactive"
	"flag"
	"os"
	"strings"
)

const MODE_EXEC = "exec"
const MODE_DEV = "dev"

func main() {

	// defer func() {
	// 	panic_err := recover()
	// 	if panic_err != nil {
	// 		fmt.Println("Panic", panic_err)
	// 	}
	// }()

	flag_mode := flag.String("mode", "exec", "Mode exec or dev (override config path)")
	flag.Parse()


	switch *flag_mode {
		case MODE_DEV:
			drvcheck.Conf.PreConfig.SetYamlConfigPath(".")
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
				default:
					drvcheck.Run()
			}
		}
	}

}
