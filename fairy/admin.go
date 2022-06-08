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

func postMessage[T any](message Message[T]) {
	origin := js.Global().Get("window").Get("location").Get("origin")
	iframeEl := js.Global().Get("document").Call("querySelector", "iframe")
	b, err := json.Marshal(message)
	if err != nil {
		panic(fmt.Errorf("fairy: cannot JSON marshal message with type '%d': %w", message.Type, err))
	}
	iframeEl.Get("contentWindow").Call("postMessage", string(b), origin)
}

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

func selectTaleByPath(path []int) hypp.Action[*State] {
	return func(state *State, _ hypp.Payload) hypp.Dispatchable {
		if equalPaths(state.Current, path) {
			return state
		}
		newState := state.clone()
		newState.Current = path
		postMessage(Message[[]int]{
			Type: MessageSelectTale,
			Data: path,
		})
		return newState
	}
}

func toggleNode(path []int) hypp.Action[*State] {
	return func(state *State, _ hypp.Payload) hypp.Dispatchable {
		newState := state.clone()
		node := state.Tree
		for _, i := range path {
			node = node.Children()[i]
		}
		node.SetIsOpen(!node.IsOpen())
		return newState
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
		renderControls(state),
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

	newState := state.clone()
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
	newState := state.clone()
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
	})

	select {} // run Go forever
}
