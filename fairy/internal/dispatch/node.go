package dispatch

import (
	"github.com/macabot/fairytale/fairy/internal/state"
	"github.com/macabot/hypp"
)

// TODO is this function used?
func ToggleNode(path []int) hypp.Action[*state.State] {
	return func(s *state.State, _ hypp.Payload) hypp.Dispatchable {
		newState := s.Clone()
		node := s.Tree
		for _, i := range path {
			node = node.Children()[i]
		}
		node.SetIsOpen(!node.IsOpen())
		return newState
	}
}
