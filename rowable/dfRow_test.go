package rowable

import (
	config "drvcheck/config"
	"strconv"
	"testing"
)


func configOverride() {
	config.Conf.IsLoaded = true
	config.Conf.ConfigYaml.Unit = "GB"
	config.Conf.ConfigYaml.Csv.Mode = "daily"
	config.Conf.ConfigYaml.Csv.Dir = "."
	config.Conf.ConfigYaml.Csv.Header = []string{
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
		erow.DfFill(mr, config.Conf.ConfigYaml.Unit)
		
		if len(erow.Errs) > 0 {
			errs = append(errs, erow.Errs...)
		}

		if erow.Row.Filesystem != mr[0] {
			t.Error("Invalid file system")
		}

		if tmp, _ := strconv.ParseUint(mr[1], 0, 64); erow.Row.Size != tmp {
			t.Error("Invalid size")
		}

		if tmp, _ := strconv.ParseUint(mr[2], 0, 64); erow.Row.Used != tmp {
			t.Error("Invalid used")
		}

		if tmp, _ := strconv.ParseUint(mr[3], 0, 64); erow.Row.Avail != tmp {
			t.Error("Invalid avail")
		}

		if erow.Row.Capacity != mr[4] {
			t.Error("Invalid capacity")
		}

		if erow.Row.MountedOn != mr[5] {
			t.Error("Invalid mountedOn")
		}

		// mem unit is took from config
		if config.Conf.ConfigYaml.Unit != erow.Row.MemUnit {
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
	erow.DfFill([]string{"/dev/d0", "3000", "10", "2990", "90%", "."}, config.Conf.ConfigYaml.Unit)

	var val1 uint64
	if erow.errParseInt("30", &val1); val1 != 30 {
		t.Error("parseMemInt invalid value")
	}

	var val2 uint64
	if erow.errParseInt("xxx", &val2); len(erow.Errs) == 0 {
		t.Error("parseMemInt expected parseUint error")
	}

}


func Test_stringify(t *testing.T) {

	configOverride()

	erow := ErrRow{}
	rowArr := []string{"/dev/d0", "3000", "10", "2990", "90%", "."}
	erow.DfFill(rowArr, config.Conf.ConfigYaml.Unit)
	erow.Stringify(config.Conf.ConfigYaml.Csv.Header)

}
