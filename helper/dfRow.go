package helper

import (
	"fmt"
	"strconv"
)

type Row struct {
	Filesystem    string
	Size          int64
	Used          int64
	Avail         int64
	Capacity      string
	IsUsed        int64
	IsFree        int64
	IsUsedPercent string
	MountedOn     string
	Time 		  string
}

// type DfRow interface {
// 	fill([]string) error
// 	store() error
// }

func (row *Row) fill(args []string) error {

	row.Filesystem = args[0]

	tmp, err := strconv.ParseInt(args[1], 0, 64)
	if err != nil {
		return err
	}
	row.Size = tmp

	tmp2, err := strconv.ParseInt(args[2], 0, 64)
	if err != nil {
		return err
	}
	row.Used = tmp2

	tmp3, err := strconv.ParseInt(args[3], 0, 64)
	if err != nil {
		return err
	}
	row.Avail = tmp3

	row.Capacity = args[4]

	tmp4, err := strconv.ParseInt(args[5], 0, 64)
	if err != nil {
		return err
	}
	row.IsUsed = tmp4

	tmp5, err := strconv.ParseInt(args[6], 0, 64)
	if err != nil {
		return err
	}
	row.IsFree = tmp5

	row.IsUsedPercent = args[7]
	row.MountedOn = args[8]
	row.Time = "" // todo: ----

	return nil
}

func (row *Row) store() error {
	// todo next
	fmt.Println(row)
	return nil
}
