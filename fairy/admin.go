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

func historyPushState(state *State) {
	href := getHref(js.Global())
	stateQuery := state.toQuery()
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

/*
TODO debug
panic: hypp: dispatchable has unexpected type '<nil>'. Expected type 'StateAndEffects[*fairy.State]', 'Action[*fairy.State]', 'ActionAndPayload[*fairy.State]' or '*fairy.State' wasm_exec.js:22:14
<empty string> wasm_exec.js:22:14
goroutine 6 [running]: wasm_exec.js:22:14
github.com/macabot/hypp.app[...].func3({0x94d80, 0xb504c0}) wasm_exec.js:22:14
	/home/michael/repos/github.com/macabot/hypp/engine.go:614 +0x27 wasm_exec.js:22:14
github.com/macabot/hypp.app[...].func1.1() wasm_exec.js:22:14
	/home/michael/repos/github.com/macabot/hypp/engine.go:581 +0x10 wasm_exec.js:22:14
github.com/macabot/hypp/driver/js.EventTarget.AddEventListener.func1({{}, 0x7ff800010000006d, 0xb4e790}, {0xb504b0, 0x1, 0x1}) wasm_exec.js:22:14
	/home/michael/repos/github.com/macabot/hypp/driver/js/js.go:98 +0x6 wasm_exec.js:22:14
syscall/js.handleEvent() wasm_exec.js:22:14
	/usr/local/go/src/syscall/js/func.go:94 +0x26
*/

func selectTaleByPath(path []int) hypp.Action[*State] {
	return func(state *State, _ hypp.Payload) hypp.Dispatchable {
		if equalPaths(state.Current, path) {
			return state
		}
		newState := state.Clone()
		newState.Current = path
		newState.TaleEvents = nil
		postMessageToIFrame(Message[[]int]{
			Type: MessageSelectTale,
			Data: path,
		})
		return newState
	}
}

func toggleNode(path []int) hypp.Action[*State] {
	return func(state *State, _ hypp.Payload) hypp.Dispatchable {
		newState := state.Clone()
		node := state.Tree
		for _, i := range path {
			node = node.Children()[i]
		}
		node.SetIsOpen(!node.IsOpen())
		return newState
	}
}

func appendTaleEvent(state *State, payload hypp.Payload) hypp.Dispatchable {
	raw := payload.(json.RawMessage)
	fmt.Println("appendTaleEvent", string(raw))
	var taleEvent TaleEvent
	if err := json.Unmarshal(raw, &taleEvent); err != nil {
		panic(fmt.Errorf("fairy: cannot unmarshal appendTaleEvent data '%s': %w", string(raw), err))
	}
	newState := state.Clone()
	newState.TaleEvents = append(newState.TaleEvents, taleEvent)
	return newState
}

func onTaleEvent(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: MessageProps{
			Type:         MessageTaleEvent,
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
	children := make([]*hypp.VNode, len(n.Children()))
	for i, child := range n.Children() {
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
				"active":    n.IsOpen(),
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
			hypp.Text(n.Name()),
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
					"caret-down": n.IsOpen(),
				},
				"onclick": toggleNode(path),
			},
			hypp.Text(n.Name()),
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

func renderRightSide(state *State) *hypp.VNode {
	return html.Div(
		hypp.HProps{"class": "right-side"},
		renderSettings(state.Settings),
		renderIFrame(state.Settings),
		renderPanel(state),
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

func selectIFrameSize(state *State, payload hypp.Payload) hypp.Dispatchable {
	event := payload.(hypp.Event)
	value := event.Target().Value()
	size := mustIFrameSizeFromString(value)

	newState := state.Clone()
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
			"onchange": hypp.Action[*State](selectIFrameSize),
		},
		options...,
	)
}

func toggleLandscape(state *State, _ hypp.Payload) hypp.Dispatchable {
	newState := state.Clone()
	newState.Settings.landscape = !newState.Settings.landscape
	return newState
}

func renderLandscapeToggle(landscape bool) *hypp.VNode {
	return html.Select(
		hypp.HProps{
			"onchange": hypp.Action[*State](toggleLandscape),
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

func renderPanel(state *State) *hypp.VNode {
	panels := []func() *hypp.VNode{
		func() *hypp.VNode { return renderControls(state) },
		func() *hypp.VNode { return renderTaleEvents(state.TaleEvents) },
	}
	controls := 0
	if tale := state.currentTale(); tale != nil {
		controls = len(tale.controls)
	}
	return html.Div(
		hypp.HProps{"class": "panel"},
		renderPanelTabs(
			state.SelectedPanelTab,
			fmt.Sprintf("Controls (%d)", controls),
			fmt.Sprintf("Events (%d)", len(state.TaleEvents)),
		),
		panels[state.SelectedPanelTab](),
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

func selectPanelTab(i int) hypp.Action[*State] {
	return func(state *State, _ hypp.Payload) hypp.Dispatchable {
		newState := state.Clone()
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

func renderControls(state *State) *hypp.VNode {
	tale := state.currentTale()
	var controls []*hypp.VNode
	if tale == nil {
		controls = []*hypp.VNode{hypp.Text("No controls: no tale has been selected")}
	} else {
		controls = make([]*hypp.VNode, len(tale.controls))
		for i, control := range tale.controls {
			controls[i] = control.Render(tale.state, state.Current, i)
		}
	}
	return html.Div(
		hypp.HProps{"class": "controls"},
		controls...,
	)
}

func renderTaleEvents(taleEvents []TaleEvent) *hypp.VNode {
	children := make([]*hypp.VNode, len(taleEvents))
	for i, taleEvent := range taleEvents {
		b, _ := json.Marshal(taleEvent.Event)
		children[i] = html.Li(
			hypp.HProps{"class": "tale-event"},
			html.Span(hypp.HProps{"class": "key"}, hypp.Text(taleEvent.Key)),
			html.Pre(nil, hypp.Text(string(b))),
		)
	}
	var child *hypp.VNode
	if len(taleEvents) == 0 {
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

func runAdmin(state *State) {
	el := js.Global().Get("document").Call("getElementById", "app")
	if el.IsNull() {
		panic("Could not find element with id 'app'.")
	}
	hypp.App(hypp.AppProps[*State]{
		Driver: jsd.Driver{},
		Init:   state,
		View: func(state *State) *hypp.VNode {
			return html.Main(
				nil,
				renderTreeView(state.Tree, state.Current),
				renderRightSide(state),
			)
		},
		DispatchWrapper: func(dispatch hypp.Dispatch) hypp.Dispatch {
			return func(dispatchable hypp.Dispatchable, payload hypp.Payload) {
				switch v := dispatchable.(type) {
				case hypp.StateAndEffects[*State]:
					historyPushState(v.State)
				case *State:
					historyPushState(v)
				}
				dispatch(dispatchable, payload)
			}
		},
		Node: jsd.Node(el),
		Subscriptions: func(state *State) []hypp.Subscription {
			return []hypp.Subscription{
				onTaleEvent(hypp.Action[*State](appendTaleEvent)),
			}
		},
	})

	select {} // run Go forever
}
