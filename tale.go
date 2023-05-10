package fairytale

import (
	"encoding/json"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/macabot/hypp"
)

type TaleTarget int

const (
	TaleInsideBody TaleTarget = iota
	TaleAsBody
)

type TaleSettings struct {
	Target TaleTarget
}

var _ Node[hypp.EmptyState] = &Tale[hypp.EmptyState]{}

// Control manages the state of a Tale. Typically, a Control manages a single
// property of the state, however a Control can change the whole state.
type Control[S hypp.State] interface {
	Render(state S, talePath []int, controlIndex int) *hypp.VNode
	UpdateFromEvent(state S, event hypp.Event) hypp.Dispatchable
	UpdateFromMessage(state S, data json.RawMessage) hypp.Dispatchable
}

type TaleEvent[S hypp.State] struct {
	Path  []int
	Label string
	State S
}

// Tale is a Node that let's you develop and document a component.
type Tale[S hypp.State] struct {
	name     string
	slug     string
	state    S                   // S
	view     func(S) *hypp.VNode // func(S) *hypp.VNode
	dispatch func(hypp.Dispatchable, hypp.Payload)

	controls        []Control[S]
	settings        TaleSettings
	events          []TaleEvent[S]
	stateSubscriber func(S)
}

// New creates a new Tale.
func New[S hypp.State](
	name string,
	state S,
	view func(S) *hypp.VNode,
) *Tale[S] {
	var tale *Tale[S]
	var dispatch func(dispatchable hypp.Dispatchable, payload hypp.Payload)
	dispatch = func(dispatchable hypp.Dispatchable, payload hypp.Payload) {
		switch v := dispatchable.(type) {
		case hypp.StateAndEffects[S]:
			tale.SetState(v.State)
			for _, effect := range v.Effects {
				effect.Effecter(dispatch, effect.Payload)
			}
		case hypp.Action[S]:
			dispatch(v(tale.state, payload), nil)
		case hypp.ActionAndPayload[S]:
			dispatch(v.Action, v.Payload)
		case S: // State
			tale.SetState(v)
		default:
			panic(fmt.Errorf("fairytale: dispatchable has unexpected type '%[1]T'. Expected type 'StateAndEffects[%[2]T]', 'Action[%[2]T]', 'ActionAndPayload[%[2]T]' or '%[2]T'", dispatchable, tale.state))
		}
	}
	tale = &Tale[S]{
		name:     name,
		slug:     slug.Make(name),
		state:    state,
		view:     view,
		dispatch: dispatch,
	}
	return tale
}

func (t Tale[S]) Name() string           { return t.name }
func (t Tale[S]) Slug() string           { return t.slug }
func (t Tale[S]) Children() []Node[S]    { return nil }
func (t *Tale[S]) Tale() *Tale[S]        { return t }
func (t Tale[S]) IsOpen() bool           { return false }
func (t *Tale[S]) SetIsOpen(isOpen bool) { /* noop */ }
func (t Tale[S]) View() *hypp.VNode      { return t.view(t.state) }
func (t Tale[S]) Controls() []Control[S] { return t.controls }
func (t Tale[S]) State() S               { return t.state }
func (t Tale[S]) Settings() TaleSettings { return t.settings }
func (t *Tale[S]) ClearEvents()          { t.events = nil }
func (t Tale[S]) Events() []TaleEvent[S] { return t.events }

func (t *Tale[S]) SetStateSubscriber(stateSubscriber func(S)) {
	t.stateSubscriber = stateSubscriber
}

func (t *Tale[S]) AppendEvent(event TaleEvent[S]) {
	t.events = append(t.events, event)
}

func (t *Tale[S]) Dispatch(dispatchable hypp.Dispatchable, payload hypp.Payload) {
	t.dispatch(dispatchable, payload)
}

func (t *Tale[S]) SetState(s S) {
	t.state = s
	if t.stateSubscriber != nil {
		t.stateSubscriber(s)
	}
}

func (t *Tale[S]) WithControls(controls ...Control[S]) *Tale[S] {
	t.controls = controls
	return t
}

func (t *Tale[S]) WithSettings(settings TaleSettings) *Tale[S] {
	t.settings = settings
	return t
}
