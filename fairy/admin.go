package fairy

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"syscall/js"

	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
	"github.com/macabot/hypp/tag/html"
)

func consoleDebug(args ...any) {
	js.Global().Get("console").Call("debug", args...)
}

func consoleWarn(args ...any) {
	js.Global().Get("console").Call("warn", args...)
}

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

func historyPushState(s *state) {
	href := getHref(js.Global())
	stateQuery := s.toQuery()
	if equalQuery(href.Query(), stateQuery) {
		return
	}
	href.RawQuery = stateQuery.Encode()
	js.Global().Get("history").Call("pushState", map[string]any{}, "", href.String())
}

func equalPaths(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}

func selectTaleByPath(path []int) hypp.Action[*state] {
	return func(s *state, _ hypp.Payload) hypp.Dispatchable {
		if equalPaths(s.Current, path) {
			return s
		}
		newState := s.clone()
		newState.Current = path
		newState.TaleEvents = nil
		postMessageToIFrame(message[[]int]{
			Type: messageSelectTale,
			Data: path,
		})
		return newState
	}
}

func toggleNode(path []int) hypp.Action[*state] {
	return func(s *state, _ hypp.Payload) hypp.Dispatchable {
		newState := s.clone()
		node := s.Tree
		for _, i := range path {
			node = node.children()[i]
		}
		node.setIsOpen(!node.isOpen())
		return newState
	}
}

func appendTaleEvent(s *state, payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	var event taleEvent
	if err := json.Unmarshal(raw, &event); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal appendTaleEvent data '%s': %w", string(raw), err))
	}
	newState := s.clone()
	newState.TaleEvents = append(newState.TaleEvents, event)
	return newState
}

func onTaleEvent(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: messageProps{
			Type:         messageTaleEvent,
			Dispatchable: dispatchable,
		},
	}
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

func renderNode(n Node, isRoot bool, path []int, current []int) *hypp.VNode {
	children := make([]*hypp.VNode, len(n.children()))
	for i, child := range n.children() {
		childPath := make([]int, len(path)+1)
		copy(childPath, path)
		childPath[len(childPath)-1] = i
		children[i] = renderNode(child, false, childPath, current)
	}
	ul := html.Ul(
		hypp.HProps{
			"class": map[string]bool{
				"nested":    !isRoot,
				"tree-root": isRoot,
				"active":    n.isOpen(),
			},
			"key": pathToKey(path),
		},
		children...,
	)
	selected := equalPaths(path, current)
	if isRoot {
		return ul
	} else if len(children) == 0 {
		return html.Li(
			hypp.HProps{
				"class": map[string]bool{
					"tree-tale": true,
					"selected":  selected,
				},
				"onclick": selectTaleByPath(path),
			},
			hypp.Text(n.name()),
		)
	}
	return html.Li(
		hypp.HProps{
			"class": map[string]bool{
				"selected": selected,
			},
		},
		html.Span(
			hypp.HProps{
				"class": map[string]bool{
					"caret":      true,
					"caret-down": n.isOpen(),
				},
				"onclick": toggleNode(path),
			},
			hypp.Text(n.name()),
		),
		ul,
	)
}

func renderTreeView(tree Node, current []int) *hypp.VNode {
	return html.Div(
		hypp.HProps{"class": "tree-view"},
		renderNode(tree, true, nil, current),
	)
}

func renderRightSide(s *state) *hypp.VNode {
	return html.Div(
		hypp.HProps{"class": "right-side"},
		renderSettings(s.Settings),
		renderIFrame(s.Settings),
		renderPanel(s),
	)
}

func renderSettings(settings AdminSettings) *hypp.VNode {
	return html.Div(
		hypp.HProps{"class": "settings"},
		renderIFrameSizeSelect(settings.iFrameSize),
		renderLandscapeToggle(settings.landscape),
	)
}

func renderIFrameSize(size IFrameSize, selected bool) *hypp.VNode {
	return html.Option(
		hypp.HProps{
			"value":    size.String(),
			"selected": selected,
		},
		hypp.Text(size.String()),
	)
}

func selectIFrameSize(s *state, payload hypp.Payload) hypp.Dispatchable {
	event := payload.(hypp.Event)
	value := event.Target().Value()
	size := mustIFrameSizeFromString(value)

	newState := s.clone()
	newState.Settings.iFrameSize = size
	return newState
}

func renderIFrameSizeSelect(size IFrameSize) *hypp.VNode {
	options := make([]*hypp.VNode, len(IFrameSizes))
	for i, s := range IFrameSizes {
		options[i] = renderIFrameSize(s, s.Equal(size))
	}
	return html.Select(
		hypp.HProps{
			"onchange": hypp.Action[*state](selectIFrameSize),
		},
		options...,
	)
}

func toggleLandscape(s *state, _ hypp.Payload) hypp.Dispatchable {
	newState := s.clone()
	newState.Settings.landscape = !newState.Settings.landscape
	return newState
}

func renderLandscapeToggle(landscape bool) *hypp.VNode {
	return html.Select(
		hypp.HProps{
			"onchange": hypp.Action[*state](toggleLandscape),
		},
		html.Option(
			hypp.HProps{
				"value":    "0",
				"selected": !landscape,
			},
			hypp.Text("Portrait"),
		),
		html.Option(
			hypp.HProps{
				"value":    "1",
				"selected": landscape,
			},
			hypp.Text("Landscape"),
		),
	)
}

func renderIFrame(settings AdminSettings) *hypp.VNode {
	size := settings.iFrameSize
	if settings.landscape {
		size.Swap()
	}
	divProps := hypp.HProps{"class": "current-tale"}
	iFrameProps := hypp.HProps{
		"src": "/",
	}
	if size[0] != 0 && size[1] != 0 {
		divProps["style"] = map[string]string{
			"min-height": fmt.Sprintf("%dpx", size[1]),
		}
		iFrameProps["style"] = map[string]string{
			"width":  fmt.Sprintf("%dpx", size[0]),
			"height": fmt.Sprintf("%dpx", size[1]),
		}
	}
	return html.Div(
		divProps,
		html.Iframe(iFrameProps),
	)
}

func renderPanel(s *state) *hypp.VNode {
	panels := []func() *hypp.VNode{
		func() *hypp.VNode { return renderControls(s) },
		func() *hypp.VNode { return renderTaleEvents(s.TaleEvents) },
	}
	controls := 0
	if tale := s.currentTale(); tale != nil {
		controls = len(tale.myControls)
	}
	return html.Div(
		hypp.HProps{"class": "panel"},
		renderPanelTabs(
			s.SelectedPanelTab,
			fmt.Sprintf("Controls (%d)", controls),
			fmt.Sprintf("Events (%d)", len(s.TaleEvents)),
		),
		panels[s.SelectedPanelTab](),
	)
}

func renderPanelTabs(selectedTab int, names ...string) *hypp.VNode {
	children := make([]*hypp.VNode, len(names))
	for i, name := range names {
		children[i] = renderPanelTab(i, name, i == selectedTab)
	}
	return html.Div(
		hypp.HProps{"class": "panel-tabs"},
		children...,
	)
}

func selectPanelTab(i int) hypp.Action[*state] {
	return func(s *state, _ hypp.Payload) hypp.Dispatchable {
		newState := s.clone()
		newState.SelectedPanelTab = i
		return newState
	}
}

func renderPanelTab(i int, name string, selected bool) *hypp.VNode {
	return html.Span(
		hypp.HProps{
			"class": map[string]bool{
				"panel-tab": true,
				"selected":  selected,
			},
			"onclick": selectPanelTab(i),
		},
		hypp.Text(name),
	)
}

func renderControls(s *state) *hypp.VNode {
	tale := s.currentTale()
	var controls []*hypp.VNode
	if tale == nil {
		controls = []*hypp.VNode{hypp.Text("No controls: no tale has been selected")}
	} else {
		controls = make([]*hypp.VNode, len(tale.myControls))
		for i, control := range tale.myControls {
			controls[i] = control.Render(tale.myState, s.Current, i)
		}
	}
	return html.Div(
		hypp.HProps{"class": "controls"},
		controls...,
	)
}

func renderTaleEvents(events []taleEvent) *hypp.VNode {
	children := make([]*hypp.VNode, len(events))
	for i, taleEvent := range events {
		b, _ := json.Marshal(taleEvent.Event)
		children[i] = html.Li(
			hypp.HProps{"class": "tale-event"},
			html.Span(hypp.HProps{"class": "key"}, hypp.Text(taleEvent.Key)),
			html.Pre(nil, hypp.Text(string(b))),
		)
	}
	var child *hypp.VNode
	if len(events) == 0 {
		child = hypp.Text("[No events have been triggered]")
	} else {
		child = html.Ul(nil, children...)
	}
	return html.Div(
		hypp.HProps{"class": "tale-events"},
		child,
	)
}

func getHref(window js.Value) *url.URL {
	href, err := url.Parse(window.Get("location").Get("href").String())
	if err != nil {
		panic("Could not parse window.location.href as URL.")
	}
	return href
}

func runAdmin(s *state) {
	el := js.Global().Get("document").Call("getElementById", "app")
	if el.IsNull() {
		panic("Could not find element with id 'app'.")
	}
	hypp.App(hypp.AppProps[*state]{
		Driver: jsd.Driver{},
		Init:   s,
		View: func(s *state) *hypp.VNode {
			return html.Main(
				nil,
				renderTreeView(s.Tree, s.Current),
				renderRightSide(s),
			)
		},
		DispatchWrapper: func(dispatch hypp.Dispatch) hypp.Dispatch {
			return func(dispatchable hypp.Dispatchable, payload hypp.Payload) {
				switch v := dispatchable.(type) {
				case hypp.StateAndEffects[*state]:
					historyPushState(v.State)
				case *state:
					historyPushState(v)
				}
				dispatch(dispatchable, payload)
			}
		},
		Node: jsd.Node(el),
		Subscriptions: func(_ *state) []hypp.Subscription {
			return []hypp.Subscription{
				onTaleEvent(hypp.Action[*state](appendTaleEvent)),
			}
		},
	})

	select {} // run Go forever
}
