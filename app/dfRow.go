package drvcheck

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

type Row struct {
	Filesystem    string
	Size          uint64
	Used          uint64
	Avail         uint64
	Capacity      string
	MountedOn     string
	Time          string // time of command execution, not from "df" output
	MemUnit       string
}


type ErrRow struct {
	row Row
	errs []error
}


func (erow *ErrRow) dfFill(args []string) {

	if len(args) < 6 {
		erow.errs = append(erow.errs, errors.New("fill: invalid args"))
		return
	}

	conf, _ := GetConfig()
	erow.StrFill(
		args[0],
		args[1],
		args[2],
		args[3],
		args[4],
		args[5],
		time.Now().Local().Format("2006-01-02 15:04:05"),
		conf.ConfigYaml.Unit,
	)

}


func (erow *ErrRow) StrFill(filesystem, size, used, avail, capacity, mountedOn, time, memUnit string) {
	erow.row.Filesystem = filesystem
	erow.errParseInt(size, &erow.row.Size)
	erow.errParseInt(used, &erow.row.Used)
	erow.errParseInt(avail, &erow.row.Avail)
	erow.row.Capacity = capacity
	erow.row.MountedOn = mountedOn
	erow.row.Time = time
	erow.row.MemUnit = memUnit
}


func (erow *ErrRow) errParseInt(value string, result *uint64) {
	res, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		erow.errs = append(erow.errs, err)
	} 
	*result = res
}

func (erow *ErrRow) _stringify() []string {

	var tmp []string
	conf, _ := GetConfig()
	helms := conf.ConfigYaml.Csv.Header

	for _, helm := range helms {

		refRowFieldType, found := reflect.TypeOf(erow.row).FieldByName(helm)
		if !found {
			panic("CSV header item not found in erow")
		}
		refRowFieldValue := reflect.ValueOf(erow.row).FieldByName(helm)

		switch refRowFieldType.Type {
			case reflect.TypeOf((uint64)(0)):
				tmp = append(tmp, strconv.FormatUint(refRowFieldValue.Uint(), 10))
			default:
				tmp = append(tmp, refRowFieldValue.String())
		}
	}
	
	return tmp
}
