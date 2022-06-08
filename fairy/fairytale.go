package fairy

import (
	"syscall/js"

	"github.com/macabot/hypp"
)

func Run(tree Node, assets []*hypp.VNode) {
	top := js.Global().Get("top")
	inTopFrame := js.Global().Get("self").Equal(top)
	state := &State{Tree: tree, Assets: assets}
	href := getHref(top)
	state.updateFromQuery(href.Query())
	if inTopFrame {
		runAdmin(state)
	} else {
		runApp(state)
	}
}
