package component

import (
	"fmt"

	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func AppPage(s *fairytale.State) *hypp.VNode {
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
		html.Title(nil, hypp.Text(taleToTitle(currentTale))),
	)
	currentTaleNode := CurrentTale(currentTale)
	target := fairytale.TaleInsideBody
	if currentTale != nil {
		target = currentTale.Settings().Target
	}
	var body *hypp.VNode
	switch target {
	case fairytale.TaleInsideBody:
		body = html.Body(nil, currentTaleNode)
	case fairytale.TaleAsBody:
		body = currentTaleNode
	default:
		panic(fmt.Errorf("invalid target %v", target))
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
