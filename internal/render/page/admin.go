package page

import (
	"github.com/macabot/fairytale/internal/render/component"
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func AdminPage(s *state.State) *hypp.VNode {
	return html.Main(
		nil,
		component.TreeView(s),
		component.RightSide(s),
	)
}
