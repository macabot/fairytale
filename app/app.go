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
func Run[S hypp.State](options *Options, nodes ...fairytale.Node[S]) {
	wrap := false
	for _, node := range nodes {
		if node.Tale() != nil {
			wrap = true
			break
		}
	}
	if wrap {
		wrapper := fairytale.NewBundle("Fairy Tales", nodes...)
		nodes = []fairytale.Node[S]{wrapper}
	}

	tree := fairytale.NewBundle("", nodes...)
	tree.SetIsOpen(true)
	s := fairytale.NewState[S](tree)
	if options != nil {
		s.SetAssets(options.Assets)
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

func runAdmin[S hypp.State](s *fairytale.State[S]) {
	el := js.Global().Get("document").Call("getElementById", "app")
	if el.IsNull() {
		panic("Could not find element with id 'app'.")
	}
	hypp.App(hypp.AppProps[*fairytale.State[S]]{
		Driver: jsd.Driver{},
		Init:   s,
		View:   component.AdminPage[S],
		Node:   jsd.Node(el),
		Subscriptions: func(_ *fairytale.State[S]) []hypp.Subscription {
			return []hypp.Subscription{
				dispatch.OnTaleEvent(hypp.Action[*fairytale.State[S]](dispatch.AppendTaleEvent[S])),
				dispatch.OnHashChange[S](),
			}
		},
	})

	select {} // run Go forever
}

func runApp[S hypp.State](s *fairytale.State[S]) {
	el := js.Global().Get("document").Call("querySelector", "html")
	if el.IsNull() {
		panic("Could not find <html> element.")
	}
	hypp.App(hypp.AppProps[*fairytale.State[S]]{
		Driver: jsd.Driver{},
		Init:   s,
		View:   component.AppPage[S],
		Node:   jsd.Node(el),
		Subscriptions: func(_ *fairytale.State[S]) []hypp.Subscription {
			return []hypp.Subscription{
				dispatch.OnSelectTale(hypp.Action[*fairytale.State[S]](dispatch.SelectTale[S])),
				dispatch.OnOperateControl(hypp.Action[*fairytale.State[S]](dispatch.OperateControl[S])),
				dispatch.OnRefreshApp(hypp.Action[*fairytale.State[S]](dispatch.RefreshApp[S])),
			}
		},
	})

	select {} // run Go forever
}
