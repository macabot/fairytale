package model

type SocketMessageType int

const (
	SocketMessageReload = iota + 1
)

type SocketMessage struct {
	Type SocketMessageType
}
