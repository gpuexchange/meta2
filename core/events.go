package core

type GlobalEventType int

const (
	Shutdown GlobalEventType = iota
)

type GlobalEvent struct {
	EventType GlobalEventType
}
