package csv

import (
	rowable "drvcheck/rowable"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"crypto/md5"
)

// var lock = &sync.Mutex{}
var csvModelInstance *CsvModel

func GetCsvModelInstance(unit string, header []string, dir string, mode string) *CsvModel {

	if csvModelInstance != nil {

		hasher := md5.New()
		hasher.Write([]byte(unit+dir+mode))
		for _, v := range header {
			hasher.Write([]byte(v))
		}
		sum1 := hasher.Sum(nil)
		
		hasher.Reset()
		hasher.Write([]byte(csvModelInstance.config.unit+csvModelInstance.config.dir+csvModelInstance.config.mode))
		for _, v := range csvModelInstance.config.header {
			hasher.Write([]byte(v))
		}
		sum2 := hasher.Sum(nil)

		if string(sum1) != string(sum2) {
			csvModelInstance = nil
		}
		
	}

	if csvModelInstance == nil {
		csvModelInstance = &CsvModel{}
		csvModelInstance.config = configCsvModel{
			unit, header, dir, mode,
		}
	}
	return csvModelInstance
}

type configCsvModel struct {
	unit string
	header []string
	dir string
	mode string
}

type CsvModel struct {
	Errs []error
	erows []rowable.ErrRow
	config configCsvModel
}

type modelError struct {
	err string
}

func (me *modelError) Error() string {
	return me.err
}

func (model *CsvModel) AddRow(row ...rowable.ErrRow) {
	model.erows = append(model.erows, row...)
}

func (model *CsvModel) Read(dateFrom time.Time) []rowable.Row {

	// if interactive.RunTestMode {
	// 	// todo 
	// }

	var rows []rowable.Row

	dateUnitNow := time.Now().Local().Unix()

	if dateFrom.Local().Unix() > dateUnitNow {
		panic("dateFrom cannot be after current time")
	}

	if model.config.mode == "daily" {
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


func (model *CsvModel) rowsFromCsvFile(content string) []rowable.Row {
	
	var rows []rowable.Row

	if content == "" {
		return rows
	}

	lines := strings.Split(content, "\n")
	header := strings.Split(lines[0], Delimiter)
	fileData := lines[1:]


	for _, line := range fileData {
		row := rowable.ErrRow{}
		args := strings.Split(line, Delimiter)
		if len(args) == 0 {
			continue
		}

		
		for hk, hitem := range header {
			switch hitem {
				case "MountedOn":
					row.Row.MountedOn = args[hk]
				case "Filesystem":
					row.Row.Filesystem = args[hk]
				case "Capacity":
					row.Row.Capacity = args[hk]
				case "Size":
					v, _ := strconv.ParseUint(args[hk], 0, 64)
					row.Row.Size = model.parseMemInt(v)
				case "Used":
					v, _ := strconv.ParseUint(args[hk], 0, 64)
					row.Row.Used = model.parseMemInt(v)
				case "Avail":
					v, _ := strconv.ParseUint(args[hk], 0, 64)
					row.Row.Avail = model.parseMemInt(v)
				case "Time":
					row.Row.Time = args[hk]
				case "MemUnit":
					row.Row.MemUnit = args[hk]
			}
		}

		

		rows = append(rows, row.Row)
	}

	

	return rows
}


func (model *CsvModel) Store() {

	fmt.Println("csvStore - store!")

	if len(model.erows) == 0 {
		fmt.Println("csvStore - nothing to store")
		return
	}

	fileContentStr, filename := model.getFile(time.Now().Local())

	if fileContentStr == "" {
		fileContentStr, _ = model.BuildHeader()
	}

	fileContentStr = fileContentStr + model.strigify()

	write_error := os.WriteFile(filename, []byte(fileContentStr), 0777)
	if write_error != nil {
		model.Errs = append(model.Errs, write_error)
		return
	}
}


func (model *CsvModel) getFile(dailyFrom time.Time) (string, string) {


	filename := model.config.dir + "/drvcheck"

	// todo: valid in config.go
	if model.config.mode == "daily" {
		filename = filename + "_" + dailyFrom.Format("2006-01-02") + ".csv"
		// fmt.Println("daily filename read | " + filename)
	} else if model.config.mode == "solid" {
		filename = filename + ".csv"
	} else {
		var err error = &modelError{
			err: "Invalid CSV mode! (valid modes: daily, solid) | given: " + model.config.mode,
		}
		model.Errs = append(model.Errs, err)
		return "", filename
	}

	fileContentByte, err := os.ReadFile(filename)

	if err != nil {
		model.Errs = append(model.Errs, err)
	}

	fileContentStr := string(fileContentByte)
	return fileContentStr, filename
}


var Delimiter = ";"
func (model *CsvModel) BuildHeader() (string, error) {
	var strheader string
	strheader = strheader + strings.Join(model.config.header, Delimiter)
	return strheader, nil
}


func (model *CsvModel) strigify() string {
	var tmp string

	for _, row := range model.erows {
		tmp += "\n" + strings.Join(row.Stringify(model.config.header), Delimiter)
	}

	return tmp
}

func (model *CsvModel) parseMemInt(value uint64) uint64 {
	switch (model.config.unit) {
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