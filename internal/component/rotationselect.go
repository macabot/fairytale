package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func OrientationSelect[S hypp.State](selectedOrientation fairytale.Orientation) *hypp.VNode {
	options := make([]*hypp.VNode, len(fairytale.Orientations))
	for i, orientation := range fairytale.Orientations {
		options[i] = html.Option(
			hypp.HProps{
				"value":    orientation.String(),
				"selected": orientation == selectedOrientation,
			},
			hypp.Text(orientation.String()),
		)
	}
	return html.Label(
		nil,
		hypp.Text("Orientation"),
		html.Select(
			hypp.HProps{
				"onchange": dispatch.SelectOrientation[S],
			},
			options...,
		),
	)
}
