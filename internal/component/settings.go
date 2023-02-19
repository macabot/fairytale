package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func Settings(settings fairytale.AdminSettings) *hypp.VNode {
	return html.Div(
		hypp.HProps{"class": "settings"},
		IFrameSizeSelect(settings.IFrameSize),
		RotationSelect(settings.Rotation),
	)
}
