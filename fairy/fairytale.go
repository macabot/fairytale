package fairy

import (
	"syscall/js"

	"github.com/macabot/fairytale/fairy/internal/state"
	"github.com/macabot/hypp"
)

// Run the Fairy Tale.
func Run(tree state.Node, assets []*hypp.VNode) {
	top := js.Global().Get("top")
	inTopFrame := js.Global().Get("self").Equal(top)
	s := &state.State{Tree: tree, Assets: assets}
	href := getHref(top)
	s.UpdateCurrentFromURL(href)
	if inTopFrame {
		runAdmin(s)
	} else {
		runApp(s)
	}
}
