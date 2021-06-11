package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	// "reflect"
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

	fileContentByte, _ := os.ReadFile(filename)
	fileContentStr := string(fileContentByte)

	if fileContentStr == "" {
		fileContentStr, _ = BuildHeader()
	}

	fileContentStr = fileContentStr + model._strigify()

	write_error := os.WriteFile(filename, []byte(fileContentStr), 0777)
	if write_error != nil {
		model.errs = append(model.errs, write_error)
		return
	}
}


var delimiter = ";"

// todo: wykrywanie czy config sie zmienil na przestrzeni tego samego pliku
func BuildHeader() (string, error) {
	var strheader string
	conf, err := GetConfig()
	strheader = strheader + strings.Join(conf.configYaml.Csv.Header, delimiter)
	return strheader, err
}


func (model *CsvModel) _strigify() string {
	var tmp []string

	conf, _ := GetConfig()
	helms := conf.configYaml.Csv.Header

	for _, elm := range helms {
		// todo: reflect on row field to get its name?
		switch(elm) {
			case _Filesystem:
				tmp = append(tmp, model.erow.row.Filesystem)
			case _Size:
				tmp = append(tmp, strconv.FormatInt(parseMemInt(model.erow.row.Size), 10))
			case _Used:
				tmp = append(tmp, strconv.FormatInt(parseMemInt(model.erow.row.Used), 10))
			case _Avail:
				tmp = append(tmp, strconv.FormatInt(parseMemInt(model.erow.row.Avail), 10))
			case _Capacity:
				tmp = append(tmp, model.erow.row.Capacity)
			// --------
			case _IsUsed:
				tmp = append(tmp, strconv.FormatInt(parseMemInt(model.erow.row.IsUsed), 10))
			case _IsFree:
				tmp = append(tmp, strconv.FormatInt(parseMemInt(model.erow.row.IsFree), 10))
			case _IsUsedPercent:
			// --------
				tmp = append(tmp, model.erow.row.IsUsedPercent)
			case _MountedOn:
				tmp = append(tmp, model.erow.row.MountedOn)
			case _Time:
				tmp = append(tmp, model.erow.row.Time)
			case _MemUnit:
				tmp = append(tmp, conf.configYaml.Unit)
		}
	}

	return "\n" + strings.Join(tmp, delimiter)
}

func parseMemInt(value int64) int64 {
	var result int64
	conf, _ := GetConfig()
	switch (conf.configYaml.Unit) {
		case "KB":
			result = value
		case "MB":
			result = value / 1024
		case "GB":
			result = value / 1024 / 1024
		case "TB":
			result = value / 1024 / 1024 / 1024
		case "PB":
			result = value / 1024 / 1024 / 1024 / 1024
		case "EB":
			result = value / 1024 / 1024 / 1024 / 1024 / 1024
		case "ZB":
			result = value / 1024 / 1024 / 1024 / 1024 / 1024 / 1024
		case "JB":
			result = value / 1024 / 1024 / 1024 / 1024 / 1024 / 1024 / 1024
	}
	return result
}