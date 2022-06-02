package fairytale

import (
	"syscall/js"

	"github.com/macabot/fairytale/fairy"
	"github.com/macabot/hypp"
)

func Run(tree fairy.Node, assets []*hypp.VNode) {
	inTopFrame := js.Global().Get("self").Equal(js.Global().Get("top"))
	if inTopFrame {
		fairy.RunAdmin(&fairy.AdminState{Tree: tree})
	} else {
		fairy.RunApp(&fairy.AppState{Tree: tree, Assets: assets})
	}
}
