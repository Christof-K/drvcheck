package interactive

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)


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


// func (_stw *StatTypeWidget) getActiveType() StatTypeWidgetType {
// 	activeType := _stw.types[0]
// 	for _, t := range _stw.types {
// 		if t.active {
// 			activeType = t
// 			break
// 		}
// 	}
// 	return activeType
// }

