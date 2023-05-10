package component

import (
	"fmt"

	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func Panel[S hypp.State](s *fairytale.State[S]) *hypp.VNode {
	tale := s.CurrentTale()
	controls := 0
	if tale != nil {
		controls = len(tale.Controls())
	}

	panels := []func() *hypp.VNode{
		func() *hypp.VNode { return Controls(s) },
	}

	return html.Div(
		hypp.HProps{"class": "panel"},
		PanelTabs[S](
			s.SelectedPanelTab(),
			fmt.Sprintf("Controls (%d)", controls),
		),
		panels[s.SelectedPanelTab()](),
	)
}
