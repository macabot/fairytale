package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func IFrameSizeSelect[S hypp.State](selectedSize fairytale.IFrameSize) *hypp.VNode {
	options := make([]*hypp.VNode, len(fairytale.IFrameSizes))
	for i, size := range fairytale.IFrameSizes {
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
				"onchange": dispatch.SelectIFrameSize[S],
			},
			options...,
		),
	)
}
