package component

import (
	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func PanelTabs(selectedTab int, names ...string) *hypp.VNode {
	children := make([]*hypp.VNode, len(names))
	for i, name := range names {
		children[i] = PanelTab(i, name, i == selectedTab)
	}
	return html.Div(
		hypp.HProps{"class": "panel-tabs"},
		children...,
	)
}

func PanelTab(i int, name string, selected bool) *hypp.VNode {
	return html.Span(
		hypp.HProps{
			"class": map[string]bool{
				"panel-tab": true,
				"selected":  selected,
			},
			"onclick": dispatch.SelectPanelTab(i),
		},
		hypp.Text(name),
	)
}
