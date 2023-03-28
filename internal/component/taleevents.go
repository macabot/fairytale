package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func TaleEvents[S hypp.State](events []fairytale.TaleEvent[S]) *hypp.VNode {
	children := make([]*hypp.VNode, len(events))
	for i, taleEvent := range events {
		children[i] = html.Li(
			hypp.HProps{"class": "tale-event"},
			html.Span(hypp.HProps{"class": "key"}, hypp.Text(taleEvent.Label)),
		)
	}
	var child *hypp.VNode
	if len(events) == 0 {
		child = hypp.Text("[No events have been triggered]")
	} else {
		child = html.Ul(nil, children...)
	}
	return html.Div(
		hypp.HProps{"class": "tale-events"},
		child,
	)
}
