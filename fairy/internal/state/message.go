package state

import (
	"github.com/macabot/hypp"
	jsd "github.com/macabot/hypp/driver/js"
)

var window hypp.Window = jsd.Driver{}.Window()

type Message[T any] struct {
	Type int
	Data T
}

const (
	MessageSelectTale = iota + 1
	MessageOperateControl
	MessageTaleEvent
	MessageRefreshApp
)

type OperateControlData[T any] struct {
	TalePath     []int
	ControlIndex int
	EventData    T
}

type TaleEvent struct {
	Key   string
	Event any
}

type MessageProps struct {
	Type         int
	Dispatchable hypp.Dispatchable
}
