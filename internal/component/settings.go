package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func Settings[S hypp.State](settings fairytale.AdminSettings) *hypp.VNode {
	return html.Div(
		hypp.HProps{"class": "settings"},
		IFrameSizeSelect[S](settings.IFrameSize),
		RotationSelect[S](settings.Rotation),
	)
}
