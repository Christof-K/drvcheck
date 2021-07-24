package interactive

import (
	"fmt"
	"time"
	csv "drvcheck/csv"

	"github.com/jroimartin/gocui"
)


type DriveSelectorWidget struct {
	name string
	x, y, w, h int
}


func (ds *DriveSelectorWidget) Layout(g *gocui.Gui) error {
	view, err := g.SetView(ds.name, ds.x, ds.y, ds.x+ds.w, ds.y+ds.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	view.Clear()
	view.FgColor = gocui.ColorWhite

	view.Title = "Choose drive"
	fmt.Fprint(view, "\n")

	if !Delms.initiated {
		csvModel := csv.GetCsvModelInstance(
			conf.ConfigYaml.Unit,
			conf.ConfigYaml.Csv.Header,
			conf.ConfigYaml.Csv.Dir,
			conf.ConfigYaml.Csv.Mode,
		)
		rows := csvModel.Read(time.Now().Local().AddDate(0, 0, GraphDaysRangeActive * -1))
		Delms.initDriveElms(rows)
	}

	for _, delm := range Delms.elms {
		var displayName string
		if delm.selected {
			displayName += "[*]"
		} else {
			displayName += "[ ]"
		}
		displayName += " " + delm.name
		fmt.Fprintln(view, displayName)
	}

	return nil
}
