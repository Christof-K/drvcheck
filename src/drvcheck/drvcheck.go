package drvcheck

import (
	"fmt"
	"os/exec"
	"strings"
)

func Run() {
	errors := check()
	for _, err := range errors {
		fmt.Println(err.Error())
	}
}

func check() []error {

	fmt.Println("Run check...")

	conf, cerr := GetConfig()
	if cerr != nil {
		return append(make([]error, 0), cerr)
	}
	if len(conf.ConfigYaml.Drivers) == 0 {
		fmt.Println("Edit config.yaml and specify drivers")
		return nil
	}
	cmd := exec.Command("df", "-bk")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return append(make([]error, 0), err)
	}

	stringOutput := string(output[:])
	lines := strings.Split(stringOutput, "\n")

	model := GetCsvModelInstance()

	for _, line := range lines[1:] {

		valid := false
		
		for _, vol := range conf.ConfigYaml.Drivers {
			if strings.Contains(line, vol) {
				valid = true
				break
			}
		}

		if valid {
			row := ErrRow{}
			// args := strings.Split(line, "\t") // 
			args := strings.Fields(line)
			row.dfFill(args)

			if row.errs != nil {
				return row.errs
			}

			model.erows = append(model.erows, row)
		}

	}

	model.store()

	return model.errs
}

