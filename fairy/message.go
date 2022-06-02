package fairy

type Message[T any] struct {
	Type int
	Data T
}

const (
	MessageSelectTale = iota + 1
	MessageOperateControl
)

type OperateControlData[T any] struct {
	TalePath     []int
	ControlIndex int
	EventData    T
}
