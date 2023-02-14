package component

import (
	"github.com/macabot/fairytale/fairy/internal/dispatch"
	"github.com/macabot/fairytale/fairy/internal/state"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func IFrameSizeSelect(selectedSize state.IFrameSize) *hypp.VNode {
	options := make([]*hypp.VNode, len(state.IFrameSizes))
	for i, size := range state.IFrameSizes {
		options[i] = html.Option(
			hypp.HProps{
				"value":    size.String(),
				"selected": size.Equal(selectedSize),
			},
			hypp.Text(size.String()),
		)
	}
	return html.Label(
		nil,
		hypp.Text("Size"),
		html.Select(
			hypp.HProps{
				"onchange": hypp.Action[*state.State](dispatch.SelectIFrameSize),
			},
			options...,
		),
	)
}
