package interactive

import (
	"fmt"
	"strconv"

	"github.com/guptarohit/asciigraph"
	"github.com/jroimartin/gocui"
)

var GraphDaysRangeActive int

type DriveStatWidget struct {
	name string
	x, y, w, h int
	// DaysRangeActive int
}

func (dstat *DriveStatWidget) Layout(g *gocui.Gui) error {

	view, err := g.SetView(dstat.name, dstat.x, dstat.y, dstat.x + dstat.w, dstat.y + dstat.h)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	view.Clear()

	// conf, _ := GetConfig()

	elm := delms.getSelected()
	if GraphDaysRangeActive == 0 {
		GraphDaysRangeActive = 7
	}
	view.Title = "Stats of " + elm.name + " | last " + strconv.FormatUint(uint64(GraphDaysRangeActive), 10) + " days"// + conf.configYaml.Unit
	// view.BgColor = gocui.ColorBlue
	
	fmt.Fprint(view, "\n")


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
