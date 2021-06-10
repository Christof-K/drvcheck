package helper

import (
	"log"
	"os/exec"
	"strings"
)

func Run() {
	errors := check()
	if errors != nil {
		for _, err := range errors {
			log.Fatalf(err.Error())
		}
	}
}

type rowHandler func()

func check() []error {

	conf, cerr := GetConfig()
	if cerr != nil {
		return append(make([]error, 0), cerr)
	}
	cmd := exec.Command("df", "-bk")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return append(make([]error, 0), err)
	}

	stringOutput := string(output[:])
	lines := strings.Split(stringOutput, "\n")

	for _, line := range lines {
		for _, vol := range conf.configYaml.Drivers {
			if strings.Contains(line, vol) {

				row := ErrRow{}
				args := strings.Fields(line)
				row.fill(args)

				// todo: error handling refactor

				if row.errs != nil {
					return row.errs
				}
				store_errs := row.store()
				if store_errs != nil {
					return store_errs
				}
				if row.errs != nil {
					return row.errs
				}
			}
		}
	}

	return nil
}
