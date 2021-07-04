package helper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/jroimartin/gocui"
)

const defaultDataPeriodDays = 20


//-------- STAT TYPE WIDGET --------//
var stw StatTypeWidget
type StatTypeWidget struct {
	name string
	x, y, w, h int
	types []StatTypeWidgetType
}

type StatTypeWidgetType struct {
	name string
	active bool
}

func (_stw *StatTypeWidget) Layout(g *gocui.Gui) error {

	view, err := g.SetView(_stw.name, _stw.x, _stw.y, _stw.x + _stw.w, _stw.y + _stw.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	view.Clear()
	view.Title = "Stat type (use TAB to switch)"

	var nameString []string
	for _, t := range _stw.types {
		if t.active {
			nameString = append(nameString, "*" + t.name)
		} else {
			nameString = append(nameString, " " + t.name)
		}
	}
	fmt.Fprint(view, strings.Join(nameString, " ∙ "))
	return nil
}

func (_stw *StatTypeWidget) getActiveType() StatTypeWidgetType {
	activeType := _stw.types[0]
	for _, t := range _stw.types {
		if t.active {
			activeType = t
			break
		}
	}
	return activeType
}

//-------- DRIVE STAT WIDGET --------//
type DriveStatWidget struct {
	name string
	x, y, w, h int
}

func (dstat *DriveStatWidget) Layout(g *gocui.Gui) error {

	view, err := g.SetView(dstat.name, dstat.x, dstat.y, dstat.x + dstat.w, dstat.y + dstat.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	view.Clear()

	// conf, _ := GetConfig()

	elm := delms.getSelected()
	view.Title = "Stats of " + elm.name + " | last " + strconv.FormatUint(defaultDataPeriodDays, 10) + " days"// + conf.configYaml.Unit
	// view.BgColor = gocui.ColorBlue
	
	fmt.Fprintln(view, "\n")


	var graphData []float64
	for _, r := range elm.data {
		graphData = append(graphData, float64(r.Used * 100 / r.Size))
	}
	
	if len(graphData) > 0 {
		graph := asciigraph.Plot(
			graphData,
			asciigraph.Height(12),
			asciigraph.Width(93),
			asciigraph.Caption("Percent of usage"),
			asciigraph.Precision(0),
		)
		fmt.Fprintln(view, graph)
	}

	return nil
}



//-------- DRIVE SELECTOR WIDGET --------//
type driveElm struct {
	selected bool
	name     string
	data     []row
}

type DriveSelectorWidget struct {
	name string
	x, y, w, h int
}

var delms driveElms

type driveElms struct {
	elms []driveElm
}

func (des *driveElms) initDriveElms() {
	
	conf, _ := GetConfig()
	csvModel := GetCsvModelInstance()
	rows := csvModel.read(time.Now().Local().AddDate(0, 0, defaultDataPeriodDays * -1))

	selb := true
	des.elms = nil


	for _, elm := range conf.configYaml.Drivers {
		
		var elmRows []row
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


var ds_initiate = true

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

//-------- /DRIVE SELECTOR WIDGET --------//

func RunInteractive() {
	
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()


	dsw := &DriveSelectorWidget{"dsw", 5, 5, 35, 16}
	dstatw := &DriveStatWidget{"dstatw", 40, 5, 100, 16}

	stw = StatTypeWidget{
		name: "stw",
		x: 40,
		y: 18,
		w: 50,
		h: 2,
	}
	stw_names := []string{"Used", "Avail"}
	for k, name := range stw_names {
		stw.types = append(stw.types, StatTypeWidgetType{name, (k==0)})
	}

	g.SetManager(
		dsw, // todo: vertical scrollable
		dstatw,
		// &stw, // disabled - todo graph, text (scrollable, goto date)
	)

	errs := keyBindingSetup(g)
	if len(errs) > 0 {
		fmt.Println(errs)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Println("gocui: main loop panic err")
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

	gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		key := 0
		for k, t := range stw.types {
			if t.active {
				stw.types[k].active = false
				key = k + 1
				if key > len(stw.types) - 1 {
					key = 0
				}
				break
			}
		}
		stw.types[key].active = true
		return nil
	})

	return errs
}
