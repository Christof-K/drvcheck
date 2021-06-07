package helper

import (
	"os/exec"
	"strings"
	"Log"
)

func Run() {
	err := driveCheck()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func driveCheck() error {

	conf, cerr := getConfig()
	if cerr != nil { return cerr } 

	cmd := exec.Command("df")
	output, err := cmd.CombinedOutput()
	if err != nil { return err } 

	stringOutput := string(output[:])
	lines := strings.Split(stringOutput, "\n")

	for _, line := range lines {
		for _, vol := range conf.Drivers {
			if strings.Contains(line, vol) {
				row := Row{}
				args := strings.Fields(line)
				fill_err := row.fill(args)
				if fill_err != nil { return fill_err }
				store_err := row.store()
				if store_err != nil { return store_err }
			}
		}
	}

	return nil
}


