package helper

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

type row struct {
	Filesystem    string
	Size          uint64
	Used          uint64
	Avail         uint64
	Capacity      string
	IsUsed        uint64
	IsFree        uint64
	IsUsedPercent string
	MountedOn     string
	Time          string
	MemUnit       string
}


type ErrRow struct {
	row row
	errs []error
}

/** @see fillFromCsv **/
func (erow *ErrRow) dfFill(args []string) {

	if len(args) < 6 {
		erow.errs = append(erow.errs, errors.New("fill: invalid args"))
		return
	}

	conf, _ := GetConfig()


	erow._fill(
		args[0],
		args[1],
		args[2],
		args[3],
		args[4],
		args[5],
		time.Now().Local().Format("2006-01-02 15:04:05"),
		conf.configYaml.Unit,
	)

}


func (erow *ErrRow) _fill(filesystem, size, used, avail, capacity, mountedOn, time, memUnit string) {
	erow.row.Filesystem = filesystem
	erow.parseMemInt(size, &erow.row.Size)
	erow.parseMemInt(used, &erow.row.Used)
	erow.parseMemInt(avail, &erow.row.Avail)
	erow.row.Capacity = capacity
	erow.row.MountedOn = mountedOn
	erow.row.Time = time
	erow.row.MemUnit = memUnit
}


func (erow *ErrRow) parseMemInt(value string, result *uint64) {
	res, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		erow.errs = append(erow.errs, err)
	} else {
		*result = _parseMemInt(res)
	}
}

func (erow *ErrRow) _strigify() []string {

	var tmp []string
	conf, _ := GetConfig()
	helms := conf.configYaml.Csv.Header

	for _, elm := range helms {

		refRowFieldType, found := reflect.TypeOf(erow.row).FieldByName(elm)
		if !found {
			continue // todo: err?
		}
		refRowFieldValue := reflect.ValueOf(erow.row).FieldByName(elm)

		switch refRowFieldType.Type {
			case reflect.TypeOf((uint64)(0)):
				tmp = append(tmp, strconv.FormatUint(refRowFieldValue.Uint(), 10))
			default:
				tmp = append(tmp, refRowFieldValue.String())
		}
	}

	return tmp
}
