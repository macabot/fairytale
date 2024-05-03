package dispatch

import (
	"github.com/macabot/fairytale"
	"github.com/macabot/hypp"
)

func SelectPanelTab[S hypp.State](s *fairytale.State[S], payload hypp.Payload) hypp.Dispatchable {
	i := payload.(int)
	newState := s.Clone()
	newState.SetSelectedPanelTab(i)
	return newState
}
