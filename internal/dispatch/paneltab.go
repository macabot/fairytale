package dispatch

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func SelectPanelTab(i int) hypp.Action[*fairytale.State] {
	return func(s *fairytale.State, _ hypp.Payload) hypp.Dispatchable {
		newState := s.Clone()
		newState.SelectedPanelTab = i
		return newState
	}
}
