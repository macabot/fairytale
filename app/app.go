package app

import (
	"net/url"
	"syscall/js"

	"github.com/macabot/fairytale"
	"github.com/macabot/fairytale/internal/component"
	"github.com/macabot/fairytale/internal/dispatch"
	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
)

type Options struct {
	Assets []*hypp.VNode
}

// Run fairytale.
func Run(options *Options, nodes ...fairytale.Node) {
	wrap := false
	for _, node := range nodes {
		if node.Tale() != nil {
			wrap = true
			break
		}
	}
	if wrap {
		wrapper := fairytale.NewBundle("Fairy Tales", nodes...)
		nodes = []fairytale.Node{wrapper}
	}

	tree := fairytale.NewBundle("", nodes...)
	tree.SetIsOpen(true)
	s := &fairytale.State{
		Tree:   tree,
		Assets: options.Assets,
	}

	top := js.Global().Get("top")
	inTopFrame := js.Global().Get("self").Equal(top)
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

func runAdmin(s *fairytale.State) {
	el := js.Global().Get("document").Call("getElementById", "app")
	if el.IsNull() {
		panic("Could not find element with id 'app'.")
	}
	hypp.App(hypp.AppProps[*fairytale.State]{
		Driver: jsd.Driver{},
		Init:   s,
		View:   component.AdminPage,
		Node:   jsd.Node(el),
		Subscriptions: func(_ *fairytale.State) []hypp.Subscription {
			return []hypp.Subscription{
				dispatch.OnTaleEvent(hypp.Action[*fairytale.State](dispatch.AppendTaleEvent)),
				dispatch.OnHashChange(),
			}
		},
	})

	select {} // run Go forever
}

func runApp(s *fairytale.State) {
	el := js.Global().Get("document").Call("querySelector", "html")
	if el.IsNull() {
		panic("Could not find <html> element.")
	}
	hypp.App(hypp.AppProps[*fairytale.State]{
		Driver: jsd.Driver{},
		Init:   s,
		View:   component.AppPage,
		Node:   jsd.Node(el),
		Subscriptions: func(_ *fairytale.State) []hypp.Subscription {
			return []hypp.Subscription{
				dispatch.OnSelectTale(hypp.Action[*fairytale.State](dispatch.SelectTale)),
				dispatch.OnOperateControl(hypp.Action[*fairytale.State](dispatch.OperateControl)),
				dispatch.OnRefreshApp(hypp.Action[*fairytale.State](dispatch.RefreshApp)),
			}
		},
	})

	select {} // run Go forever
}
