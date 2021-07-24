package csv

import (
	config "drvcheck/config"
	"fmt"
	"testing"
)

func TestCsvModelRowsFromCsvFile(t *testing.T) {

	config.Conf.ConfigYaml.Unit = "GB"

	model := GetCsvModelInstance(
		config.Conf.ConfigYaml.Unit,
		[]string{"MountedOn", "Size"},
		"/",
		"daily",
	)

	if rows := model.rowsFromCsvFile(""); len(rows) != 0 {
		t.Errorf("Unexpected output")
	}

	tmpStrRow := "MountedOn\n/dev1\n/dev2\n/dev3\n/dev4\n/dev5"
	if rows := model.rowsFromCsvFile(tmpStrRow); len(rows) != 5 {
		t.Errorf("Unexpected output")
	}

	config.Conf.IsLoaded = true
	
	tmpStrRow2 := "Size\n1048576"
	if rows := model.rowsFromCsvFile(tmpStrRow2); rows[0].Size != 1 {
		fmt.Println(rows)
		t.Errorf("Unexpected output")
	}

}

func TestBuildHeader(t *testing.T) {

	model := GetCsvModelInstance(
		config.Conf.ConfigYaml.Unit,
		[]string{
			"mountOn", "Size",
		},
		"/",
		"daily",
	)
	

	if h, _ := model.BuildHeader(); h != "mountOn"+Delimiter+"Size" {
		t.Errorf("Unexpected output: " + h)
	}
}
