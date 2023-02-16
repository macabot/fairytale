package fairy

import (
	"net/url"
	"syscall/js"

	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/fairytale/internal/render/page"
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
)

// Run the Fairy Tale.
func Run(tree Node, assets []*hypp.VNode) {
	top := js.Global().Get("top")
	inTopFrame := js.Global().Get("self").Equal(top)
	s := &state.State{Tree: tree.node(), Assets: assets}
	href := getHref(top)
	s.UpdateCurrentFromURL(href)
	if inTopFrame {
		runAdmin(s)
	} else {
		runApp(s)
	}
}

func getHref(window js.Value) *url.URL {
	href, err := url.Parse(window.Get("location").Get("href").String())
	if err != nil {
		panic("Could not parse window.location.href as URL.")
	}
	return href
}

func runAdmin(s *state.State) {
	el := js.Global().Get("document").Call("getElementById", "app")
	if el.IsNull() {
		panic("Could not find element with id 'app'.")
	}
	hypp.App(hypp.AppProps[*state.State]{
		Driver: jsd.Driver{},
		Init:   s,
		View:   page.AdminPage,
		Node:   jsd.Node(el),
		Subscriptions: func(_ *state.State) []hypp.Subscription {
			return []hypp.Subscription{
				dispatch.OnTaleEvent(hypp.Action[*state.State](dispatch.AppendTaleEvent)),
				dispatch.OnHashChange(),
			}
		},
	})

	select {} // run Go forever
}

func runApp(s *state.State) {
	el := js.Global().Get("document").Call("querySelector", "html")
	if el.IsNull() {
		panic("Could not find <html> element.")
	}
	hypp.App(hypp.AppProps[*state.State]{
		Driver: jsd.Driver{},
		Init:   s,
		View:   page.AppPage,
		Node:   jsd.Node(el),
		Subscriptions: func(_ *state.State) []hypp.Subscription {
			return []hypp.Subscription{
				dispatch.OnSelectTale(hypp.Action[*state.State](dispatch.SelectTale)),
				dispatch.OnOperateControl(hypp.Action[*state.State](dispatch.OperateControl)),
				dispatch.OnRefreshApp(hypp.Action[*state.State](dispatch.RefreshApp)),
			}
		},
	})

	select {} // run Go forever
}
