package helper

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type CsvModel struct {
	errs []error
	erows []ErrRow
}

type modelError struct {
	err string
}

func (me *modelError) Error() string {
	return me.err
}


func (model *CsvModel) store() {

	fmt.Println("csvStore - store!")

	if len(model.erows) == 0 {
		fmt.Println("csvStore - nothing to store")
		return
	}

	config, err := GetConfig()

	if err != nil {
		model.errs = append(model.errs, err)
		return
	}

	directory := config.configYaml.Csv.Dir
	filename := directory + "/drvcheck"

	// todo: valid in config.go
	if config.configYaml.Csv.Mode == "daily" {
		filename = filename + "_" + time.Now().Local().Format("2006-01-02") + ".csv"
	} else if config.configYaml.Csv.Mode == "solid" {
		filename = filename + ".csv"
	} else {
		var err error = &modelError{
			err: "Invalid CSV mode! (valid modes: daily, solid) | given: " + config.configYaml.Csv.Mode,
		}
		model.errs = append(model.errs, err)
		return
	}

	fileContentByte, _ := os.ReadFile(filename)
	fileContentStr := string(fileContentByte)

	if fileContentStr == "" {
		fileContentStr, _ = BuildHeader()
	}

	fileContentStr = fileContentStr + model.strigify()

	write_error := os.WriteFile(filename, []byte(fileContentStr), 0777)
	if write_error != nil {
		model.errs = append(model.errs, write_error)
		return
	}
}


var delimiter = ";"

func BuildHeader() (string, error) {
	var strheader string
	conf, err := GetConfig()
	strheader = strheader + strings.Join(conf.configYaml.Csv.Header, delimiter)
	return strheader, err
}


func (model *CsvModel) strigify() string {
	var tmp string

	for _, row := range model.erows {
		tmp += "\n" + strings.Join(row._strigify(), delimiter)
	}

	return tmp
}

func _parseMemInt(value uint64) uint64 {
	conf, _ := GetConfig()
	switch (conf.configYaml.Unit) {
		case "KB":
			return value
		case "MB":
			return value / 1024
		case "GB":
			return value / 1024 / 1024
		case "TB":
			return value / 1024 / 1024 / 1024
		case "PB":
			return value / 1024 / 1024 / 1024 / 1024
		case "EB":
			return value / 1024 / 1024 / 1024 / 1024 / 1024
		case "ZB":
			return value / 1024 / 1024 / 1024 / 1024 / 1024 / 1024
		case "JB":
			return value / 1024 / 1024 / 1024 / 1024 / 1024 / 1024 / 1024
	}
	return value
}