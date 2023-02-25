package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func Controls[S hypp.State](s *fairytale.State[S]) *hypp.VNode {
	tale := s.CurrentTale()
	var controls []*hypp.VNode
	if tale == nil {
		controls = []*hypp.VNode{hypp.Text("No controls: no tale has been selected")}
	} else {
		controls = make([]*hypp.VNode, len(tale.Controls()))
		for i, control := range tale.Controls() {
			controls[i] = control.Render(tale.State(), s.Current(), i)
		}
	}
	return html.Div(
		hypp.HProps{"class": "controls"},
		controls...,
	)
}
