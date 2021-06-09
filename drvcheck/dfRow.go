package helper

import (
	"fmt"
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
}

type ErrRow struct {
	row row
	errs []error
}

func (erow *ErrRow) fill(args []string) {
	erow.row.Filesystem = args[0]
	erow.parseIntWrapper(args[1], &erow.row.Size)
	erow.parseIntWrapper(args[2], &erow.row.Used)
	erow.parseIntWrapper(args[3], &erow.row.Avail)
	erow.row.Capacity = args[4]
	erow.parseIntWrapper(args[5], &erow.row.IsUsed)
	erow.parseIntWrapper(args[6], &erow.row.IsFree)
	erow.row.IsUsedPercent = args[7]
	erow.row.MountedOn = args[8]
	erow.row.Time = time.Now().Local().Format("2006-01-02 15:04:05")
}


// todo - nie dziala referencja
func (erow *ErrRow) parseIntWrapper(value string, result *int64) {
	res, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		erow.errs = append(erow.errs, err)
	} else {
		*result = res 
	}
}

func (erow *ErrRow) store() error {
	// todo next
	fmt.Println(erow)
	return nil
}
