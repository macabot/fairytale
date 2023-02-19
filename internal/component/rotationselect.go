package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func RotationSelect(selectedRotation fairytale.Rotation) *hypp.VNode {
	options := make([]*hypp.VNode, len(fairytale.Rotations))
	for i, rotation := range fairytale.Rotations {
		options[i] = html.Option(
			hypp.HProps{
				"value":    rotation.String(),
				"selected": rotation == selectedRotation,
			},
			hypp.Text(rotation.String()),
		)
	}
	return html.Label(
		nil,
		hypp.Text("Rotation"),
		html.Select(
			hypp.HProps{
				"onchange": hypp.Action[*fairytale.State](dispatch.SelectRotation),
			},
			options...,
		),
	)
}
