package interactive

import (
	rowable "drvcheck/rowable"
	config "drvcheck/config"
	"testing"
	"time"
)


func TestDriveElms(t *testing.T) {

	
	config.Conf.IsLoaded = true
	config.Conf.ConfigYaml.Drivers = []string{
		"/dev1",
		"/dev2",
		"/dev3",
	}
	conf = *config.Conf

	if Delms.initiated {
		t.Errorf("driveElements unknown initialisation")
	}


	erow1 := rowable.ErrRow{}
	erow1.StrFill(
		"/dev1",
		"1048576",
		"10000",
		"1038576",
		"1048576",
		"/Volumes/dev1",
		time.Now().Local().Format("2006-01-02 15:04:05"),
		"KB",
	)
	erow2 := rowable.ErrRow{}
	erow2.StrFill(
		"/dev2",
		"2048576",
		"10000",
		"2038576",
		"2048576",
		"/Volumes/dev2",
		time.Now().Local().Format("2006-01-02 15:04:05"),
		"KB",
	)

	erow3 := rowable.ErrRow{}
	erow3.StrFill(
		"/dev3",
		"2048576",
		"10000",
		"2038576",
		"2048576",
		"/Volumes/dev3",
		time.Now().Local().Format("2006-01-02 15:04:05"),
		"KB",
	)

	var rows []rowable.Row
	rows = append(rows, erow1.Row, erow2.Row, erow3.Row)
	Delms.initDriveElms(rows)


	if Delms.getSelected().name != "/dev1" {
		t.Errorf("Unexpected output")
	}

	Delms.selectNext()

	if Delms.getSelected().name != "/dev2" {
		t.Errorf("Unexpected output")
	}

	Delms.selectNext()

	if Delms.getSelected().name != "/dev3" {
		t.Errorf("Unexpected output")
	}

	Delms.selectNext()

	if Delms.getSelected().name != "/dev1" {
		t.Errorf("Unexpected output")
	}

	Delms.selectPrev()

	if Delms.getSelected().name != "/dev3" {
		t.Errorf("Unexpected output")
	}


	if !Delms.initiated {
		t.Errorf("driveElements should be initiated")	
	}
}