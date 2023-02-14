package dispatch

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/macabot/fairytale/fairy/internal/state"
)

func PostMessageToIFrame[T any](m state.Message[T]) {
	origin := js.Global().Get("window").Get("location").Get("origin")
	iframeEl := js.Global().Get("document").Call("querySelector", "iframe")
	b, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("fairy: cannot JSON marshal message with type '%d': %w", m.Type, err))
	}
	iframeEl.Get("contentWindow").Call("postMessage", string(b), origin)
}

func PostMessageToTopFrame[T any](m state.Message[T]) {
	origin := js.Global().Get("window").Get("location").Get("origin")
	top := js.Global().Get("top")
	b, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("fairy: cannot JSON marshal message with type '%d': %w", m.Type, err))
	}
	top.Call("postMessage", string(b), origin)
}
