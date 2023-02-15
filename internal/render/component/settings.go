package component

import (
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func Settings(settings state.AdminSettings) *hypp.VNode {
	return html.Div(
		hypp.HProps{"class": "settings"},
		IFrameSizeSelect(settings.IFrameSize),
		RotationSelect(settings.Rotation),
	)
}
