package interactive

import (
	"drvcheck/app"
	"testing"
	"time"
)


func TestDriveElms(t *testing.T) {

	drvcheck.Conf.IsLoaded = true
	drvcheck.Conf.ConfigYaml.Drivers = []string{
		"/dev1",
		"/dev2",
		"/dev3",
	}

	if Delms.initiated {
		t.Errorf("driveElements unknown initialisation")
	}

	tmpCsvModel := drvcheck.GetCsvModelInstance()
	erow1 := drvcheck.ErrRow{}
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
	erow2 := drvcheck.ErrRow{}
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

	erow3 := drvcheck.ErrRow{}
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

	tmpCsvModel.AddRow(erow1, erow2, erow3)
	Delms.initDriveElms(tmpCsvModel)


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