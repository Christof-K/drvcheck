package helper

import (
	"reflect"
	"strconv"
	"time"
)

type row struct {
	Filesystem    string
	Size          int64
	Used          int64
	Avail         int64
	Capacity      string
	IsUsed        int64
	IsFree        int64
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

func (erow *ErrRow) parseMemInt(value string, result *int64) {
	res, err := strconv.ParseInt(value, 0, 64)
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

		if refRowFieldType.Type.String() == "int64" {
			tmp = append(tmp, strconv.FormatInt(refRowFieldValue.Int(), 10))
		} else {
			tmp = append(tmp, refRowFieldValue.String())
		}
	}

	return tmp
}
