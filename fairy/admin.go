package fairy

import (
	"net/url"
	"strconv"
	"syscall/js"

	"github.com/macabot/fairytale/fairy/internal/dispatch"
	"github.com/macabot/fairytale/fairy/internal/render/page"
	"github.com/macabot/fairytale/fairy/internal/state"
	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
)

func equalQuery(a, b url.Values) bool {
	if len(a) != len(b) {
		return false
	}
	equalValues := func(u, v []string) bool {
		if len(u) != len(v) {
			return false
		}
		for i, x := range u {
			if x != v[i] {
				return false
			}
		}
		return true
	}
	for key, valuesA := range a {
		valuesB := b[key]
		if !equalValues(valuesA, valuesB) {
			return false
		}
	}
	return true
}

func subPath(full, sub []int) bool {
	if len(sub) > len(full) {
		return false
	}
	for i, x := range sub {
		if x != full[i] {
			return false
		}
	}
	return true
}

func pathToKey(p []int) string {
	k := ""
	for _, i := range p {
		if k != "" {
			k += "-"
		}
		k += strconv.Itoa(i)
	}
	if k == "" {
		return "root"
	}
	return k
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
