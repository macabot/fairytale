package component

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
	"github.com/macabot/hypp/tag/html"
)

func AdminPage(s *fairytale.State) *hypp.VNode {
	return html.Main(
		nil,
		TreeView(s),
		RightSide(s),
	)
}
