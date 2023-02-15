package page

import (
	"fmt"

	"github.com/macabot/fairytale/internal/render/component"
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func AppPage(s *state.State) *hypp.VNode {
	var assets []*hypp.VNode
	currentTale := s.GetTale(s.Current)
	if currentTale != nil {
		assets = s.Assets
	}
	headChildren := append(
		assets,
		html.Meta(hypp.HProps{"charset": "utf-8"}),
		html.Meta(hypp.HProps{
			"name":    "viewport",
			"content": "width=device-width, initial-scale=1.0",
		}),
		html.Title(nil, hypp.Text(state.TaleToTitle(currentTale))),
	)
	currentTaleNode := component.CurrentTale(currentTale)
	target := state.TaleInsideBody
	if currentTale != nil {
		target = currentTale.Settings().Target
	}
	var body *hypp.VNode
	switch target {
	case state.TaleInsideBody:
		body = html.Body(nil, currentTaleNode)
	case state.TaleAsBody:
		body = currentTaleNode
	default:
		panic(fmt.Errorf("invalid target %v", currentTale.Settings().Target))
	}
	return html.Html(
		nil,
		html.Head(
			nil,
			headChildren...,
		),
		body,
	)
}
