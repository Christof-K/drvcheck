package interactive

import (
	"fmt"
	"github.com/jroimartin/gocui"
)


var Delms driveElms
var RunTestMode = false

func Run() {
	
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
		Delms.selectNext()
		return nil
	})
	if err2 != nil {
		errs = append(errs, err2)
	}

	err3 := gui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		Delms.selectPrev()
		return nil
	})
	if err3 != nil {
		errs = append(errs, err3)
	}

	
	graphDaysRange := []int{1, 3, 7, 14, 30, 90}

	gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {


		key := 0
		for k, dr := range graphDaysRange {
			if dr == GraphDaysRangeActive {
				if k + 1 <= len(graphDaysRange) - 1 {
					key = k + 1
				}
			}
		}
		
		GraphDaysRangeActive = graphDaysRange[key]
		return nil
	})

	return errs
}

