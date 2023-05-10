package component

import (
	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func PanelTabs[S hypp.State](selectedTab int, names ...string) *hypp.VNode {
	children := make([]*hypp.VNode, len(names))
	for i, name := range names {
		children[i] = PanelTab[S](i, name, i == selectedTab)
	}
	return html.Div(
		hypp.HProps{"class": "panel-tabs"},
		children...,
	)
}

func PanelTab[S hypp.State](i int, name string, selected bool) *hypp.VNode {
	return html.Span(
		hypp.HProps{
			"class": map[string]bool{
				"panel-tab": true,
				"selected":  selected,
			},
			"onclick": dispatch.SelectPanelTabAction[S](i),
		},
		hypp.Text(name),
	)
}
