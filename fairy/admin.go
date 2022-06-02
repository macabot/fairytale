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
	Tree    Node
	Current []int
}

func (s AdminState) getTale(path []int) *Tale {
	node := s.Tree
	for _, i := range path {
		node = node.Children()[i]
	}
	return node.Tale()
}

func (s AdminState) CurrentTale() *Tale {
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

func renderIFrame() *hypp.VNode {
	return html.Div(
		hypp.HProps{"class": "current-tale"},
		html.Iframe(
			hypp.HProps{
				"src": "/",
			},
		),
	)
}

func renderControls(state *AdminState) *hypp.VNode {
	tale := state.CurrentTale()
	if tale == nil {
		return hypp.Text("No controls: no tale has been selected")
	}
	controls := make([]*hypp.VNode, len(tale.controls))
	for i, control := range tale.controls {
		controls[i] = control.Render(tale.state, state.Current, i)
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
				renderIFrame(),
				renderControls(state),
			)
		},
		Node: jsd.Node(el),
	})

	select {} // run Go forever
}
