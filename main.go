package main

import (
	helper "drvcheck/drvcheck"
	"fmt"
	"os"
)

func main() {

	defer func() {
		panic_err := recover()
		if panic_err != nil {
			fmt.Println("Panic", panic_err)
		}
	}()


	var args []string
	args = os.Args[1:]

	fmt.Println(args)

	if len(args) == 0 {
		helper.Run()
	} else {
		for _, arg := range args {
			switch(arg) {
				case "dev":
					helper.Conf.PreConfig.SetYamlConfigPath(".")
				// default:
			}
			helper.Run()
		}
	}

}
