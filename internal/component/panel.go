package component

import (
	"fmt"

	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func Panel[S hypp.State](s *fairytale.State[S]) *hypp.VNode {
	tale := s.CurrentTale()
	var events []fairytale.TaleEvent[S]
	controls := 0
	if tale != nil {
		events = tale.Events()
		controls = len(tale.Controls())
	}

	panels := []func() *hypp.VNode{
		func() *hypp.VNode { return Controls(s) },
		func() *hypp.VNode { return TaleEvents(events) },
	}

	return html.Div(
		hypp.HProps{"class": "panel"},
		PanelTabs[S](
			s.SelectedPanelTab(),
			fmt.Sprintf("Controls (%d)", controls),
			fmt.Sprintf("Events (%d)", len(events)),
		),
		panels[s.SelectedPanelTab()](),
	)
}
