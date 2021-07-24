package drvcheck

import (
	csv "drvcheck/csv"
	config "drvcheck/config"
	rowable "drvcheck/rowable"
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

	conf, cerr := config.GetConfig()
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

	model := csv.GetCsvModelInstance(
		conf.ConfigYaml.Unit,
		conf.ConfigYaml.Csv.Header,
		conf.ConfigYaml.Csv.Dir,
		conf.ConfigYaml.Csv.Mode,
	)

	for _, line := range lines[1:] {

		valid := false
		
		for _, vol := range conf.ConfigYaml.Drivers {
			if strings.Contains(line, vol) {
				valid = true
				break
			}
		}

		if valid {
			row := rowable.ErrRow{}
			args := strings.Fields(line)
			row.DfFill(args, config.Conf.ConfigYaml.Unit)

			if row.Errs != nil {
				return row.Errs
			}

			model.AddRow(row)
		}

	}

	model.Store()

	return model.Errs
}

