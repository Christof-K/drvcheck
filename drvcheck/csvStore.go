package helper

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type CsvModel struct {
	errs []error
	erow ErrRow
}

type modelError struct {
	err string
}

func (me *modelError) Error() string {
	return me.err
}


func (model *CsvModel) store(erow ErrRow) {

	fmt.Println("csvStore - store!")

	model.erow = erow
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

	// todo: table header jesli pusty pliczek
	fileContentByte, _ := os.ReadFile(filename)
	fileContentStr := string(fileContentByte)
	fileContentStr = fileContentStr + model._strigify()

	write_error := os.WriteFile(filename, []byte(fileContentStr), 0777)
	if write_error != nil {
		model.errs = append(model.errs, write_error)
		return
	}
}

func (model *CsvModel) _strigify() string {
	var result string

	result = "\n"
	result = result + model.erow.row.Filesystem + ","
	result = result + strconv.FormatInt(model.erow.row.Size, 10) + ","
	result = result + strconv.FormatInt(model.erow.row.Used, 10) + ","
	result = result + strconv.FormatInt(model.erow.row.Avail, 10) + ","
	result = result + model.erow.row.Capacity + ","
	result = result + strconv.FormatInt(model.erow.row.IsUsed, 10) + ","
	result = result + strconv.FormatInt(model.erow.row.IsFree, 10) + ","
	result = result + model.erow.row.IsUsedPercent + ","
	result = result + model.erow.row.MountedOn + ","
	result = result + model.erow.row.Time

	return result
}