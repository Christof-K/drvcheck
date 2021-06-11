package helper

import (
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
	row  row
	errs []error
}

func (erow *ErrRow) fill(args []string) {

	conf, _ := GetConfig()

	erow.row.Filesystem = args[0]
	erow.parseMemInt(args[1], &erow.row.Size)
	erow.parseMemInt(args[2], &erow.row.Used)
	erow.parseMemInt(args[3], &erow.row.Avail)
	erow.row.Capacity = args[4]
	erow.row.MountedOn = args[5]
	erow.row.Time = time.Now().Local().Format("2006-01-02 15:04:05")
	erow.row.MemUnit = conf.configYaml.Unit
}

func (erow *ErrRow) parseMemInt(value string, result *uint64) { // todo: reflect na polu zamiast zwracac pointer?
	res, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		erow.errs = append(erow.errs, err)
	} else {
		*result = _parseMemInt(res)
	}
}

func (erow *ErrRow) _strigify() []string {

	var tmp []string
	helms := conf.configYaml.Csv.Header

	for _, elm := range helms {

		refRowFieldType, found := reflect.TypeOf(erow.row).FieldByName(elm)
		if found == false {
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
