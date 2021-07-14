package drvcheck

import (
	"fmt"
	"testing"
)

func TestCsvModelRowsFromCsvFile(t *testing.T) {

	model := GetCsvModelInstance()

	if rows := model.rowsFromCsvFile(""); len(rows) != 0 {
		t.Errorf("Unexpected output")
	}

	tmpStrRow := "MountedOn\n/dev1\n/dev2\n/dev3\n/dev4\n/dev5"
	if rows := model.rowsFromCsvFile(tmpStrRow); len(rows) != 5 {
		t.Errorf("Unexpected output")
	}

	Conf.IsLoaded = true
	Conf.ConfigYaml.Unit = "GB"
	tmpStrRow2 := "Size\n1048576"
	if rows := model.rowsFromCsvFile(tmpStrRow2); rows[0].Size != 1 {
		fmt.Println(rows)
		t.Errorf("Unexpected output")
	}

}

func TestBuildHeader(t *testing.T) {
	Conf.IsLoaded = true
	Conf.ConfigYaml.Csv.Header = []string{
		"Filesystem", "Size",
	}
	if h, _ := BuildHeader(); h != "Filesystem"+Delimiter+"Size" {
		t.Errorf("Unexpected output")
	}
}
