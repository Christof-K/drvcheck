package drvcheck

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var csvModelInstance *CsvModel

func GetCsvModelInstance() *CsvModel {
	if csvModelInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if csvModelInstance == nil {
			csvModelInstance = &CsvModel{}
		}
	}
	return csvModelInstance
}

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

func (model *CsvModel) AddRow(row ...ErrRow) {
	model.erows = append(model.erows, row...)
}

func (model *CsvModel) Read(dateFrom time.Time) []Row {

	var rows []Row

	dateUnitNow := time.Now().Local().Unix()

	if dateFrom.Local().Unix() > dateUnitNow {
		panic("dateFrom cannot be after current time")
	}

	config, _ := GetConfig()
	if config.ConfigYaml.Csv.Mode == "daily" {
		for {
			tmp, _ := model.getFile(dateFrom.Local())
			rows = append(rows, model.rowsFromCsvFile(tmp)...)
			
			dateFrom = dateFrom.AddDate(0, 0, 1)
			if dateFrom.Local().Unix() >= dateUnitNow {
				break
			}
		}

	} else {
		tmp, _ := model.getFile(time.Now().Local())
		rows = model.rowsFromCsvFile(tmp)
	}

	return rows
}


func (model *CsvModel) rowsFromCsvFile(content string) []Row {
	
	var rows []Row

	if content == "" {
		return rows
	}

	lines := strings.Split(content, "\n")
	header := strings.Split(lines[0], Delimiter)
	fileData := lines[1:]


	for _, line := range fileData {
		row := ErrRow{}
		args := strings.Split(line, Delimiter)
		if len(args) == 0 {
			continue
		}

		
		for hk, hitem := range header {
			switch hitem {
				case "MountedOn":
					row.row.MountedOn = args[hk]
				case "Filesystem":
					row.row.Filesystem = args[hk]
				case "Capacity":
					row.row.Capacity = args[hk]
				case "Size":
					v, _ := strconv.ParseUint(args[hk], 0, 64)
					row.row.Size = _parseMemInt(v)
				case "Used":
					v, _ := strconv.ParseUint(args[hk], 0, 64)
					row.row.Used = _parseMemInt(v)
				case "Avail":
					v, _ := strconv.ParseUint(args[hk], 0, 64)
					row.row.Avail = _parseMemInt(v)
				case "Time":
					row.row.Time = args[hk]
				case "MemUnit":
					row.row.MemUnit = args[hk]
			}
		}

		

		rows = append(rows, row.row)
	}

	

	return rows
}


func (model *CsvModel) store() {

	fmt.Println("csvStore - store!")

	if len(model.erows) == 0 {
		fmt.Println("csvStore - nothing to store")
		return
	}

	fileContentStr, filename := model.getFile(time.Now().Local())

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


func (model *CsvModel) getFile(dailyFrom time.Time) (string, string) {
	config, err := GetConfig()

	if err != nil {
		model.errs = append(model.errs, err)
		return "", ""
	}

	filename := config.ConfigYaml.Csv.Dir + "/drvcheck"

	// todo: valid in config.go
	if config.ConfigYaml.Csv.Mode == "daily" {
		filename = filename + "_" + dailyFrom.Format("2006-01-02") + ".csv"
		// fmt.Println("daily filename read | " + filename)
	} else if config.ConfigYaml.Csv.Mode == "solid" {
		filename = filename + ".csv"
	} else {
		var err error = &modelError{
			err: "Invalid CSV mode! (valid modes: daily, solid) | given: " + config.ConfigYaml.Csv.Mode,
		}
		model.errs = append(model.errs, err)
		return "", filename
	}

	fileContentByte, err := os.ReadFile(filename)

	if err != nil {
		model.errs = append(model.errs, err)
	}

	fileContentStr := string(fileContentByte)
	return fileContentStr, filename
}


var Delimiter = ";"

func BuildHeader() (string, error) {
	var strheader string
	conf, err := GetConfig()
	strheader = strheader + strings.Join(conf.ConfigYaml.Csv.Header, Delimiter)
	return strheader, err
}


func (model *CsvModel) strigify() string {
	var tmp string

	for _, row := range model.erows {
		tmp += "\n" + strings.Join(row._stringify(), Delimiter)
	}

	return tmp
}

func _parseMemInt(value uint64) uint64 {
	conf, _ := GetConfig()
	switch (conf.ConfigYaml.Unit) {
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