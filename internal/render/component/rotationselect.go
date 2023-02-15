package component

import (
	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func RotationSelect(selectedRotation state.Rotation) *hypp.VNode {
	options := make([]*hypp.VNode, len(state.Rotations))
	for i, rotation := range state.Rotations {
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
				"onchange": hypp.Action[*state.State](dispatch.SelectRotation),
			},
			options...,
		),
	)
}
