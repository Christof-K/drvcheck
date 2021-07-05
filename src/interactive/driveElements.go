package interactive

import (
	"drvcheck/src/drvcheck"
	"time"
)

type driveElm struct {
	selected bool
	name     string
	data     []drvcheck.Row
}

type driveElms struct {
	elms []driveElm
}


func (des *driveElms) initDriveElms() {
	
	conf, _ := drvcheck.GetConfig()
	csvModel := drvcheck.GetCsvModelInstance()
	rows := csvModel.Read(time.Now().Local().AddDate(0, 0, defaultDataPeriodDays * -1))

	selb := true
	des.elms = nil

	for _, elm := range conf.ConfigYaml.Drivers {
		
		var elmRows []drvcheck.Row
		for _, r := range rows {
			if r.Filesystem == elm || r.MountedOn == elm {
				elmRows = append(elmRows, r)
			}
		}

		tmp := driveElm{selb, elm, elmRows}
		des.elms = append(des.elms, tmp)
		selb = false
	}
}

func (des *driveElms) getSelected() driveElm {

	var selectedElm driveElm
	for _, elm := range des.elms {
		if elm.selected {
			selectedElm = elm
			break
		}
	}

	return selectedElm
}

func (de *driveElms) selectNext() {

	key := 0
	for k, elm := range de.elms {
	
		if elm.selected {
			key = k + 1
			if key > len(de.elms) - 1 {
				key = 0
			}
			de.elms[k].selected = false
			break
		}
	}
	
	de.elms[key].selected = true
}

func (de *driveElms) selectPrev() {
	key := 0
	for k, elm := range de.elms {

		if elm.selected {
			de.elms[k].selected = false
			key = k - 1
			if key < 0 {
				key = len(de.elms) - 1
			}
			break
		}
	}

	de.elms[key].selected = true
}