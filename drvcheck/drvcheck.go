package helper

import (
	"Log"
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

	conf, cerr := getConfig()
	if cerr != nil { return append(make([]error, 0), cerr) } 

	cmd := exec.Command("df")
	output, err := cmd.CombinedOutput()
	if err != nil { return append(make([]error, 0), err) } 

	stringOutput := string(output[:])
	lines := strings.Split(stringOutput, "\n")

	for _, line := range lines {
		for _, vol := range conf.Drivers {
			if strings.Contains(line, vol) {
				row := ErrRow{}
				args := strings.Fields(line)
				row.fill(args)
				if row.errs != nil { return row.errs }
				row.store()
				if row.errs != nil { return row.errs }
			}
		}
	}

	return nil
}

