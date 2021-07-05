package interactive

import (
	"fmt"

	"github.com/jroimartin/gocui"
)


var ds_initiate = true

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
	fmt.Fprintln(view, "\n")

	if ds_initiate {
		delms.initDriveElms()
		ds_initiate = false
	}

	for _, delm := range delms.elms {
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