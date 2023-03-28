package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func RightSide[S hypp.State](s *fairytale.State[S]) *hypp.VNode {
	return html.Div(
		hypp.HProps{"class": "right-side"},
		Settings[S](s.Settings()),
		IFrame(s.CurrentTale(), s.Settings()),
		Panel(s),
	)
}
