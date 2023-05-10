package dispatch

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func SelectPanelTabAction[S hypp.State](i int) hypp.Action[*fairytale.State[S]] {
	return func(s *fairytale.State[S], _ hypp.Payload) hypp.Dispatchable {
		newState := s.Clone()
		newState.SetSelectedPanelTab(i)
		return newState
	}
}
