package helper

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func RunInteractive() {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.SetManagerFunc(mainLayout)

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
	if err != nil { errs = append(errs, err) }

	// navigate
	err2 := gui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		delms.selectNext()
		return nil
	})
	if err2 != nil { errs = append(errs, err2) }

	err3 := gui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		delms.selectPrev()
		return nil
	})
	if err3 != nil { errs = append(errs, err3) }


	return errs
}

const fontRed = "\x1b[0;31m"
var initiate = true

func mainLayout(gui *gocui.Gui) error {
	if view, err := gui.SetView("main", 2, 2, 60, 20); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		fmt.Fprintln(view, fontRed + "Drvcheck")
		fmt.Fprintln(view, "\nChoose drive:")

		if initiate {
			delms.initDriveElms()
			initiate = false
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

	}
	return nil
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
			break;
		}
	}

	return selectedElm
}

func (de *driveElms) selectNext() {
	
	hitNext := false
	for _, elm := range de.elms {
		if hitNext {
			elm.selected = true
			hitNext = false
			break	
		}

		if elm.selected {
			elm.selected = false
			hitNext = true
			continue
		} 
	}

	if hitNext {
		de.elms[0].selected = true
	}

}

func (de *driveElms) selectPrev() {
	hitKey := 0
	for k, elm := range de.elms {

		if elm.selected {
			elm.selected = false
			hitKey = k - 1
			if hitKey < 0 {
				hitKey = len(de.elms) - 1
			}
			break
		} 
	}

	de.elms[hitKey].selected = true
}


type driveElm struct {
	selected bool
	name string
	used []string
}

