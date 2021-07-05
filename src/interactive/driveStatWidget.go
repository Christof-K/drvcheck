package interactive

import (
	"fmt"
	"strconv"

	"github.com/guptarohit/asciigraph"
	"github.com/jroimartin/gocui"
)

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
