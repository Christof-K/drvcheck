package helper

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const fontRed = "\x1b[0;31m"


//-------- DRIVE STAT WIDGET --------//
type DriveStatWidget struct {
	name string
	x, y int
	w    int
	h    int
}

func (dstat *DriveStatWidget) Layout(g *gocui.Gui) error {
	view, err := g.SetView(dstat.name, dstat.x, dstat.y, dstat.x + dstat.w, dstat.y + dstat.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	view.Clear()

	fmt.Fprintln(view, "Content")
	fmt.Fprintln(view, "todo: stat data chart")

	// todo -------
	// elm := delms.getSelected()
	// elm.used ....

	// todo draw chart

	return nil
}



//-------- DRIVE SELECTOR WIDGET --------//
type driveElm struct {
	selected bool
	name     string
	used     []string
}

type DriveSelectorWidget struct {
	name string
	x, y int
	w    int
	h    int
}

var delms driveElms

type driveElms struct {
	elms []driveElm
}

func (des *driveElms) initDriveElms() {
	conf, _ := GetConfig()

	selb := true
	des.elms = nil

	for _, elm := range conf.configYaml.Drivers {
		// var used []string

		// todo: get data

		tmp := driveElm{selb, elm, []string{""}} // todo used data
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


var ds_initiate = true

func (ds *DriveSelectorWidget) Layout(g *gocui.Gui) error {
	view, err := g.SetView(ds.name, ds.x, ds.y, ds.x+ds.w, ds.y+ds.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	view.Clear()

	fmt.Fprintln(view, fontRed+"Drvcheck")
	fmt.Fprintln(view, "\nChoose drive:")

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

//-------- /DRIVE SELECTOR WIDGET --------//

func RunInteractive() {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	dsw := &DriveSelectorWidget{"dsw", 5, 5, 35, 15}
	dstatw := &DriveStatWidget{"dstatw", 40, 5, 50, 15}
	g.SetManager(dsw, dstatw)

	errs := keyBindingSetup(g)
	if len(errs) > 0 {
		fmt.Println(errs)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}

}

func keyBindingSetup(gui *gocui.Gui) []error {
	var errs []error

	// quit
	err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		return gocui.ErrQuit
	})
	if err != nil {
		errs = append(errs, err)
	}

	// navigate
	err2 := gui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		delms.selectNext()
		return nil
	})
	if err2 != nil {
		errs = append(errs, err2)
	}

	err3 := gui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		delms.selectPrev()
		return nil
	})
	if err3 != nil {
		errs = append(errs, err3)
	}

	return errs
}
