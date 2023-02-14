package dispatch

import (
	"github.com/macabot/fairytale/fairy/internal/state"
	"github.com/macabot/hypp"
)

func SelectPanelTab(i int) hypp.Action[*state.State] {
	return func(s *state.State, _ hypp.Payload) hypp.Dispatchable {
		newState := s.Clone()
		newState.SelectedPanelTab = i
		return newState
	}
}
