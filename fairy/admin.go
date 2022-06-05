package fairy

import (
	"encoding/json"
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
	"github.com/macabot/hypp/tag/html"
)

type AdminState struct {
	hypp.EmptyState
	Tree     Node
	Current  []int
	Settings AdminSettings
}

type IFrameSize [2]int

var (
	SizeDesktop        = IFrameSize{0, 0}
	Size_iPhone_11_Pro = IFrameSize{375, 812}
)

var IFrameSizes = [...]IFrameSize{
	SizeDesktop,
	Size_iPhone_11_Pro,
}

func (i *IFrameSize) Swap() {
	i[0], i[1] = i[1], i[0]
}

func (i IFrameSize) Equal(other IFrameSize) bool {
	return i[0] == other[0] && i[1] == other[1]
}

func (i IFrameSize) String() string {
	switch i {
	case SizeDesktop:
		return "Desktop"
	case Size_iPhone_11_Pro:
		return "iPhone 11 Pro"
	default:
		panic(fmt.Errorf("unknown IFrameSize: %d, %d", i[0], i[1]))
	}
}

func IFrameSizeFromString(s string) IFrameSize {
	switch s {
	case "Desktop":
		return SizeDesktop
	case "iPhone 11 Pro":
		return Size_iPhone_11_Pro
	default:
		panic(fmt.Errorf("cannot convert '%s' to IFrameSize", s))
	}
}

type AdminSettings struct {
	iFrameSize IFrameSize
	landscape  bool
}

func (s AdminState) getTale(path []int) *Tale {
	node := s.Tree
	for _, i := range path {
		node = node.Children()[i]
	}
	return node.Tale()
}

func (s AdminState) currentTale() *Tale {
	return s.getTale(s.Current)
}

func (s AdminState) clone() *AdminState {
	return &s
}

func postMessage[T any](message Message[T]) hypp.Effect {
	origin := js.Global().Get("window").Get("location").Get("origin")
	return hypp.Effect{
		Effecter: func(_ hypp.Dispatch, _ hypp.Payload) {
			iframeEl := js.Global().Get("document").Call("querySelector", "iframe")
			b, err := json.Marshal(message)
			if err != nil {
				panic(fmt.Errorf("fairy: cannot JSON marshal message with type '%d': %w", message.Type, err))
			}
			iframeEl.Get("contentWindow").Call("postMessage", string(b), origin)
		},
	}
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

func selectTaleByPath(path []int) hypp.Action[*AdminState] {
	return func(state *AdminState, _ hypp.Payload) hypp.Dispatchable {
		if equalPaths(state.Current, path) {
			return state
		}
		newState := state.clone()
		newState.Current = path
		return hypp.StateAndEffects[*AdminState]{
			State: newState,
			Effects: []hypp.Effect{
				postMessage(Message[[]int]{
					Type: MessageSelectTale,
					Data: path,
				}),
			},
		}
	}
}

func toggleNode(path []int) hypp.Action[*AdminState] {
	return func(state *AdminState, _ hypp.Payload) hypp.Dispatchable {
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

func renderNode(n Node, isRoot bool, path []int) *hypp.VNode {
	children := make([]*hypp.VNode, len(n.Children()))
	for i, child := range n.Children() {
		childPath := make([]int, len(path)+1)
		copy(childPath, path)
		childPath[len(childPath)-1] = i
		children[i] = renderNode(child, false, childPath)
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
	if isRoot {
		return ul
	} else if len(children) == 0 {
		return html.Li(
			hypp.HProps{
				"class": map[string]bool{
					"selected": n.IsSelected(),
				},
				"onclick": selectTaleByPath(path),
			},
			hypp.Text(n.Name()),
		)
	}
	return html.Li(
		hypp.HProps{
			"class": map[string]bool{
				"selected": n.IsSelected(),
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

func renderTreeView(tree Node) *hypp.VNode {
	return html.Div(
		hypp.HProps{"class": "tree-view"},
		renderNode(tree, true, nil),
	)
}

func renderRightSide(state *AdminState) *hypp.VNode {
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

func selectIFrameSize(state *AdminState, payload hypp.Payload) hypp.Dispatchable {
	event := payload.(hypp.Event)
	value := event.Target().Value()
	size := IFrameSizeFromString(value)

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
			"onchange": hypp.Action[*AdminState](selectIFrameSize),
		},
		options...,
	)
}

func toggleLandscape(state *AdminState, _ hypp.Payload) hypp.Dispatchable {
	newState := state.clone()
	newState.Settings.landscape = !newState.Settings.landscape
	return newState
}

func renderLandscapeToggle(landscape bool) *hypp.VNode {
	return html.Select(
		hypp.HProps{
			"onchange": hypp.Action[*AdminState](toggleLandscape),
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

func renderControls(state *AdminState) *hypp.VNode {
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

func RunAdmin(state *AdminState) {
	el := js.Global().Get("document").Call("getElementById", "app")
	if el.IsNull() {
		panic("Could not find element with id 'app'.")
	}
	hypp.App(hypp.AppProps[*AdminState]{
		Driver: jsd.Driver{},
		Init:   state,
		View: func(state *AdminState) *hypp.VNode {
			return html.Main(
				nil,
				renderTreeView(state.Tree),
				renderRightSide(state),
			)
		},
		Node: jsd.Node(el),
	})

	select {} // run Go forever
}
