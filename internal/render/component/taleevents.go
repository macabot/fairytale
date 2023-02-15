package component

import (
	"encoding/json"

	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func TaleEvents(events []state.TaleEvent) *hypp.VNode {
	children := make([]*hypp.VNode, len(events))
	for i, taleEvent := range events {
		b, _ := json.Marshal(taleEvent.Event)
		children[i] = html.Li(
			hypp.HProps{"class": "tale-event"},
			html.Span(hypp.HProps{"class": "key"}, hypp.Text(taleEvent.Key)),
			html.Pre(nil, hypp.Text(string(b))),
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
