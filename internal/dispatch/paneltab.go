package dispatch

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func SelectPanelTab(i int) hypp.Action[*fairytale.State] {
	return func(s *fairytale.State, _ hypp.Payload) hypp.Dispatchable {
		newState := s.Clone()
		newState.SetSelectedPanelTab(i)
		return newState
	}
}
