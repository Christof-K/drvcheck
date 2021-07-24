package rowable

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
	Row Row
	Errs []error
}


func (erow *ErrRow) DfFill(args []string, memUnit string) {

	if len(args) < 6 {
		erow.Errs = append(erow.Errs, errors.New("fill: invalid args"))
		return
	}

	erow.StrFill(
		args[0],
		args[1],
		args[2],
		args[3],
		args[4],
		args[5],
		time.Now().Local().Format("2006-01-02 15:04:05"),
		memUnit,
	)

}


func (erow *ErrRow) StrFill(filesystem, size, used, avail, capacity, mountedOn, time, memUnit string) {
	erow.Row.Filesystem = filesystem
	erow.errParseInt(size, &erow.Row.Size)
	erow.errParseInt(used, &erow.Row.Used)
	erow.errParseInt(avail, &erow.Row.Avail)
	erow.Row.Capacity = capacity
	erow.Row.MountedOn = mountedOn
	erow.Row.Time = time
	erow.Row.MemUnit = memUnit
}


func (erow *ErrRow) errParseInt(value string, result *uint64) {
	res, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		erow.Errs = append(erow.Errs, err)
	} 
	*result = res
}

func (erow *ErrRow) Stringify(headerElements []string) []string {

	var tmp []string

	for _, helm := range headerElements {

		refRowFieldType, found := reflect.TypeOf(erow.Row).FieldByName(helm)
		if !found {
			panic("CSV header item not found in erow")
		}
		refRowFieldValue := reflect.ValueOf(erow.Row).FieldByName(helm)

		switch refRowFieldType.Type {
			case reflect.TypeOf((uint64)(0)):
				tmp = append(tmp, strconv.FormatUint(refRowFieldValue.Uint(), 10))
			default:
				tmp = append(tmp, refRowFieldValue.String())
		}
	}
	
	return tmp
}
