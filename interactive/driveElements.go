package interactive

import (
	"drvcheck/app"
	"time"
)

type driveElm struct {
	selected bool
	name     string
	data     []drvcheck.Row
}

type driveElms struct {
	elms []driveElm
	initiated bool
}


func (drvElms *driveElms) initDriveElms(csvModel *drvcheck.CsvModel) {

	conf,_ := drvcheck.GetConfig()
	
	drvElms.elms = nil
	drvElms.initiated = true
	
	rows := csvModel.Read(time.Now().Local().AddDate(0, 0, GraphDaysRangeActive * -1))

	for k, elm := range conf.ConfigYaml.Drivers {
		
		var elmRows []drvcheck.Row
		for _, r := range rows {
			if r.Filesystem == elm || r.MountedOn == elm {
				elmRows = append(elmRows, r)
			}
		}

		tmp := driveElm{(k==0), elm, elmRows}
		drvElms.elms = append(drvElms.elms, tmp)
	}
}

func (drvElms *driveElms) getSelected() driveElm {

	var selectedElm driveElm
	for _, elm := range drvElms.elms {
		if elm.selected {
			selectedElm = elm
			break
		}
	}

	return selectedElm
}

func (drvElms *driveElms) selectNext() {

	key := 0
	for k, elm := range drvElms.elms {
	
		if elm.selected {
			key = k + 1
			if key > len(drvElms.elms) - 1 {
				key = 0
			}
			drvElms.elms[k].selected = false
			break
		}
	}
	
	drvElms.elms[key].selected = true
}

func (drvElms *driveElms) selectPrev() {

	key := 0
	for k, elm := range drvElms.elms {
		if elm.selected {
			drvElms.elms[k].selected = false
			key = k - 1
			if key < 0 {
				key = len(drvElms.elms) - 1
			}
			break
		}
	}

	drvElms.elms[key].selected = true
}