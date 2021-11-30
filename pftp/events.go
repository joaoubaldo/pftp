package pftp

import (
	"errors"
)

type EventType string

const (
	ClientConnectEventType    EventType = "client-connect"
	ClientDisconnectEventType EventType = "client-disconnect"
	ClientCommandEventType    EventType = "client-command"
	DataTransferEventType     EventType = "data-transfer"
	ErrorEventType            EventType = "error"
)

type ClientConnectEvent struct {
	RemoteAddr  string
	ClientCount int32
}

type ClientDisconnectEvent struct {
	RemoteAddr  string
	ClientCount int32
}

type ClientCommandEvent struct {
	RemoteAddr string
	Command    string
}

type DataTransferEvent struct {
	SrcAddr string
	DstAddr string
	Bytes   int
}

type ErrorEvent struct {
	RemoteAddr   string
	ErrorMessage string
}

type Event struct {
	name    EventType
	payload interface{}
}

func (e Event) Name() EventType {
	return e.name
}

func (e Event) Payload() interface{} {
	return e.payload
}

type EventChan chan Event

func NewEventChan(bufSize int) EventChan {
	return make(EventChan, bufSize)
}

func (eventC EventChan) Send(event Event) error {
	select {
	case eventC <- event:
		return nil
	default:
		return errors.New("no receivers")
	}
}
