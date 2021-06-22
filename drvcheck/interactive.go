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
	// switching tabs

	return errs
}

const fontRed = "\x1b[0;31m"

func mainLayout(gui *gocui.Gui) error {
	if view, err := gui.SetView("main", 2, 2, 60, 20); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		fmt.Fprintln(view, fontRed + "Drvcheck")
		fmt.Fprintln(view, "\nChoose drive:")

		// todo: drivers listed from config
		delms := driveElms{}
		delms.initDriveElms()

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
	// todo
}

func (de *driveElms) selectPrev() {
	// todo
}


type driveElm struct {
	selected bool
	name string
	used []string
}

