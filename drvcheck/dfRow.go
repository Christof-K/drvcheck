package helper

import (
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
	MemUnit		  string
}

const _Filesystem = "Filesystem"
const _Size = "Size"
const _Used = "Used"
const _Avail = "Avail"
const _Capacity = "Capacity"
const _IsUsed = "IsUsed"
const _IsFree = "IsFree"
const _IsUsedPercent = "IsUsedPercent"
const _MountedOn = "MountedOn"
const _Time = "Time"
const _MemUnit = "MemUnit"

type ErrRow struct {
	row row
	errs []error
}

func (erow *ErrRow) fill(args []string) {

	conf, _ := GetConfig()

	erow.row.Filesystem = args[0]
	erow.parseIntWrapper(args[1], &erow.row.Size)
	erow.parseIntWrapper(args[2], &erow.row.Used)
	erow.parseIntWrapper(args[3], &erow.row.Avail)
	erow.row.Capacity = args[4]
	// erow.parseIntWrapper(args[5], &erow.row.IsUsed)
	// erow.parseIntWrapper(args[6], &erow.row.IsFree)
	// erow.row.IsUsedPercent = args[7]
	erow.row.MountedOn = args[5]
	erow.row.Time = time.Now().Local().Format("2006-01-02 15:04:05")
	erow.row.MemUnit = conf.configYaml.Unit
}


func (erow *ErrRow) parseIntWrapper(value string, result *int64) {
	res, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		erow.errs = append(erow.errs, err)
	} else {
		*result = res 
	}
}

func (erow *ErrRow) _strigify() []string {

	var tmp []string
	helms := conf.configYaml.Csv.Header

	for _, elm := range helms {
		// todo: reflect on row field to get its name?
		switch(elm) {
			case _Filesystem:
				tmp = append(tmp, erow.row.Filesystem)
			case _Size:
				tmp = append(tmp, strconv.FormatInt(parseMemInt(erow.row.Size), 10))
			case _Used:
				tmp = append(tmp, strconv.FormatInt(parseMemInt(erow.row.Used), 10))
			case _Avail:
				tmp = append(tmp, strconv.FormatInt(parseMemInt(erow.row.Avail), 10))
			case _Capacity:
				tmp = append(tmp, erow.row.Capacity)
			// --------
			case _IsUsed:
				tmp = append(tmp, strconv.FormatInt(parseMemInt(erow.row.IsUsed), 10))
			case _IsFree:
				tmp = append(tmp, strconv.FormatInt(parseMemInt(erow.row.IsFree), 10))
			case _IsUsedPercent:
			// --------
				tmp = append(tmp, erow.row.IsUsedPercent)
			case _MountedOn:
				tmp = append(tmp, erow.row.MountedOn)
			case _Time:
				tmp = append(tmp, erow.row.Time)
			case _MemUnit:
				tmp = append(tmp, conf.configYaml.Unit)
		}
	}

	return tmp
}

// todo: zapis calosci a nie kazdego rowa osobono - bez sensu
// func (erow *ErrRow) store() []error {

// 	model := CsvModel{}
// 	model.store(*erow)
	
// 	return model.errs
// }
