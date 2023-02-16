package dispatch

import (
	"github.com/macabot/fairytale/internal/state"
	"github.com/macabot/hypp"
)

func OnRefreshApp(dispatchable hypp.Dispatchable) hypp.Subscription {
	return hypp.Subscription{
		Subscriber: onMessage,
		Payload: messageProps{
			Type:         messageRefreshApp,
			Dispatchable: dispatchable,
		},
	}
}

func RefreshApp(s *state.State, payload hypp.Payload) hypp.Dispatchable {
	return s.Clone()
}

/*
FIXME imports from senet to fairy types no longer work

➜  fairytale git:(feat/poc) GOOS=js GOARCH=wasm go build -o main.wasm main.go
# github.com/macabot/senet/internal/app/view/tale/page
../../internal/app/view/tale/page/gamepage.go:9:28: undefined: fairy.Tale
../../internal/app/view/tale/page/gamepage.go:10:15: undefined: fairy.NewTale
../../internal/app/view/tale/page/gamepage.go:17:9: undefined: fairy.NewCheckboxControl
# github.com/macabot/senet/internal/app/view/tale/component
../../internal/app/view/tale/component/board.go:9:25: undefined: fairy.Tale
../../internal/app/view/tale/component/board.go:11:15: undefined: fairy.NewTale
../../internal/app/view/tale/component/board.go:18:9: undefined: fairy.NewSelectControl
../../internal/app/view/tale/component/board.go:107:12: undefined: fairy.SelectOption
../../internal/app/view/tale/component/board.go:115:9: undefined: fairy.NewCheckboxControl
../../internal/app/view/tale/component/board.go:125:9: undefined: fairy.NewSelectControl
../../internal/app/view/tale/component/board.go:134:12: undefined: fairy.SelectOption
../../internal/app/view/tale/component/piece.go:9:25: undefined: fairy.Tale
../../internal/app/view/tale/component/stick.go:8:25: undefined: fairy.Tale
../../internal/app/view/tale/component/sticks.go:26:26: undefined: fairy.Tale
../../internal/app/view/tale/component/board.go:134:12: too many errors

*/
