package drvcheck

import (
	"strconv"
	"testing"
)


func configOverride() {
	Conf.isLoaded = true
	Conf.ConfigYaml.Unit = "GB"
	Conf.ConfigYaml.Csv.Mode = "daily"
	Conf.ConfigYaml.Csv.Dir = "."
	Conf.ConfigYaml.Csv.Header = []string{
		"Filesystem", "Size", "Used", "Avail", "Capacity", "MountedOn", "Time", "MemUnit",
	}
}



func TestDfRow(t *testing.T) {

	configOverride()

	mockupRow := [][]string{
		{"\\test\\test", "3000", "10", "2990", "90%", "."},
		{"/test/test", "0", "0", "0", "0%", "/dev"},
		
	}

	mockupRowErr := [][]string{
		{"", "", "", "", "", ""},
		// {nil, nil, nil, nil, nil},
	}

	// not expeting any errors
	if err1 := _testDfRow(mockupRow, t); err1 != nil {
		for _, e := range err1 {
			t.Error(e.Error())
		}
	}

	// expecting df fill errors
	if err2 := _testDfRow(mockupRowErr, t); err2 == nil {
		t.Error("mockupRowErr - expecting fill errors")
	}


}

func _testDfRow(rows [][]string, t *testing.T) []error {
	var errs []error

	for _, mr := range rows {
		erow := ErrRow{}
		erow.dfFill(mr)
		
		if len(erow.errs) > 0 {
			errs = append(errs, erow.errs...)
		}

		if erow.row.Filesystem != mr[0] {
			t.Error("Invalid file system")
		}

		if tmp, _ := strconv.ParseUint(mr[1], 0, 64); erow.row.Size != tmp {
			t.Error("Invalid size")
		}

		if tmp, _ := strconv.ParseUint(mr[2], 0, 64); erow.row.Used != tmp {
			t.Error("Invalid used")
		}

		if tmp, _ := strconv.ParseUint(mr[3], 0, 64); erow.row.Avail != tmp {
			t.Error("Invalid avail")
		}

		if erow.row.Capacity != mr[4] {
			t.Error("Invalid capacity")
		}

		if erow.row.MountedOn != mr[5] {
			t.Error("Invalid mountedOn")
		}

		// mem unit is took from config
		if Conf.ConfigYaml.Unit != erow.row.MemUnit {
			t.Error("Invalid MemUnit")
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}


func TestDfRowErrParseInt(t *testing.T) {

	configOverride()

	erow := ErrRow{}
	erow.dfFill([]string{"/dev/d0", "3000", "10", "2990", "90%", "."})

	var val1 uint64
	if erow.errParseInt("30", &val1); val1 != 30 {
		t.Error("parseMemInt invalid value")
	}

	var val2 uint64
	if erow.errParseInt("xxx", &val2); len(erow.errs) == 0 {
		t.Error("parseMemInt expected parseUint error")
	}

}


func Test_stringify(t *testing.T) {

	configOverride()

	erow := ErrRow{}
	rowArr := []string{"/dev/d0", "3000", "10", "2990", "90%", "."}
	erow.dfFill(rowArr)

	erow._stringify()

}
