package fairy

import (
	"syscall/js"

	"github.com/macabot/hypp"
)

// Run the Fairy Tale.
func Run(tree Node, assets []*hypp.VNode) {
	top := js.Global().Get("top")
	inTopFrame := js.Global().Get("self").Equal(top)
	s := &state{Tree: tree, Assets: assets}
	href := getHref(top)
	s.updateCurrentFromURL(href)
	if inTopFrame {
		runAdmin(s)
	} else {
		runApp(s)
	}
}
