package component

import (
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func RightSide(s *state.State) *hypp.VNode {
	return html.Div(
		hypp.HProps{"class": "right-side"},
		Settings(s.Settings),
		IFrame(s.CurrentTale(), s.Settings),
		Panel(s),
	)
}
