package component

import (
	"fmt"

	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func Panel[S hypp.State](s *fairytale.State[S]) *hypp.VNode {
	panels := []func() *hypp.VNode{
		func() *hypp.VNode { return Controls(s) },
		func() *hypp.VNode { return TaleEvents(s.TaleEvents()) },
	}
	controls := 0
	if tale := s.CurrentTale(); tale != nil {
		controls = len(tale.Controls())
	}
	return html.Div(
		hypp.HProps{"class": "panel"},
		PanelTabs[S](
			s.SelectedPanelTab(),
			fmt.Sprintf("Controls (%d)", controls),
			fmt.Sprintf("Events (%d)", len(s.TaleEvents())),
		),
		panels[s.SelectedPanelTab()](),
	)
}
