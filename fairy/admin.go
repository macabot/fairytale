package fairy

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
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

func (s *AdminState) updateFromQuery(query url.Values) {
	fmt.Println("has path", query.Has("path"))
	if query.Has("path") {
		var path []int
		if err := json.Unmarshal([]byte(query.Get("path")), &path); err != nil {
			consoleWarn("Could not parse query param 'path'.")
		} else {
			fmt.Println("set path", path)
			s.Current = path
		}
	}
	if query.Has("iFrameSize") {
		if size, err := iFrameSizeFromSlug(query.Get("iFrameSize")); err != nil {
			consoleWarn("Could not parse query param 'iFrameSize'.")
		} else {
			s.Settings.iFrameSize = size
		}
	}
	if query.Has("landscape") {
		if landscape, err := strconv.ParseBool(query.Get("landscape")); err != nil {
			consoleWarn("Could not parse query param 'landscape'.")
		} else {
			s.Settings.landscape = landscape
		}
	}
}

func (s AdminState) toQuery() url.Values {
	query := url.Values{}
	if s.Current != nil {
		b, err := json.Marshal(s.Current)
		if err != nil {
			panic("Could not JSON marshal AdminState.Current")
		}
		query.Set("path", string(b))
	}
	query.Set("iFrameSize", s.Settings.iFrameSize.Slug())
	query.Set("landscape", strconv.FormatBool(s.Settings.landscape))
	return query
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

func (i IFrameSize) Slug() string {
	return strings.ReplaceAll(i.String(), " ", "-")
}

func iFrameSizeFromSlug(s string) (IFrameSize, error) {
	return iFrameSizeFromString(strings.ReplaceAll(s, "-", " "))
}

func mustIFrameSizeFromString(s string) IFrameSize {
	size, err := iFrameSizeFromString(s)
	if err != nil {
		panic(err)
	}
	return size
}

func iFrameSizeFromString(s string) (IFrameSize, error) {
	switch s {
	case "Desktop":
		return SizeDesktop, nil
	case "iPhone 11 Pro":
		return Size_iPhone_11_Pro, nil
	default:
		return [2]int{}, fmt.Errorf("cannot convert '%s' to IFrameSize", s)
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

func historyPushState(state *AdminState) {
	href := getHref()
	href.RawQuery = state.toQuery().Encode()
	fmt.Println(">>", href.String())
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

func selectTaleByPath(path []int) hypp.Action[*AdminState] {
	return func(state *AdminState, _ hypp.Payload) hypp.Dispatchable {
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

func getHref() *url.URL {
	href, err := url.Parse(js.Global().Get("location").Get("href").String())
	if err != nil {
		panic("Could not parse window.location.href as URL.")
	}
	return href
}

func RunAdmin(state *AdminState) {
	el := js.Global().Get("document").Call("getElementById", "app")
	if el.IsNull() {
		panic("Could not find element with id 'app'.")
	}
	href := getHref()
	state.updateFromQuery(href.Query())
	hypp.App(hypp.AppProps[*AdminState]{
		Driver: jsd.Driver{},
		Init:   state,
		View: func(state *AdminState) *hypp.VNode {
			return html.Main(
				nil,
				renderTreeView(state.Tree, state.Current),
				renderRightSide(state),
			)
		},
		DispatchInitializer: func(dispatch hypp.Dispatch) hypp.Dispatch {
			return func(dispatchable hypp.Dispatchable, payload hypp.Payload) {
				switch v := dispatchable.(type) {
				case hypp.StateAndEffects[*AdminState]:
					historyPushState(v.State)
				case *AdminState:
					historyPushState(v)
				}
				dispatch(dispatchable, payload)
			}
		},
		Node: jsd.Node(el),
	})

	select {} // run Go forever
}
